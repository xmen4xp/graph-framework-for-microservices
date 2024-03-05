package light

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
	baselighttrafficlightcomv1 "trafficlight/build/apis/light.trafficlight.com/v1"
	nexus_client "trafficlight/build/nexus-client"
	"trafficlightdemo/pkg/perf"
	utils "trafficlightdemo/pkg/utils"

	basedesiredconfiglighttrafficlightcomv1 "trafficlight/build/apis/desiredconfiglight.trafficlight.com/v1"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type LightGroup struct {
	nexusClient *nexus_client.Clientset
	ctxt        context.Context
	groupId     int
	numLights   int
	doneFn      context.CancelFunc
}

type Light struct {
	nexusClient                 *nexus_client.Clientset
	groupId                     int
	lightId                     int
	ctxt                        context.Context
	doneFn                      context.CancelFunc
	currentLight                string
	lightMutex                  sync.RWMutex
	signalStop                  chan bool
	signalStopped               chan bool
	autoStatusReportingEnabled  bool
	autoStatusReportingPeriodic time.Duration
	autoStatusReportingOnUpdate bool
}

var (
	Lights            sync.Map
	TrafficLightGroup *LightGroup

	perfTestType       string
	perfTestInProgress bool
	perfTestMutex      sync.RWMutex

	perfTestObjectCount     int64
	perfTestFirstObjectTime time.Time
	perfTestLastObjectTime  time.Time
)

var log *logrus.Logger

func (l *LightGroup) getName(id int) string {
	return fmt.Sprintf("light-%d-%d", l.groupId, id)
}

func (l *Light) getName() string {
	return fmt.Sprintf("light-%d-%d", l.groupId, l.lightId)
}

func (g *LightGroup) getCachedLight(lightId int) *Light {
	l, ok := Lights.Load(g.getName(lightId))
	if ok {
		return l.(*Light)
	}

	newl := &Light{
		nexusClient:                 g.nexusClient,
		groupId:                     g.groupId,
		lightId:                     lightId,
		ctxt:                        g.ctxt,
		doneFn:                      g.doneFn,
		autoStatusReportingPeriodic: time.Duration(time.Minute),
		autoStatusReportingOnUpdate: true,
	}

	Lights.Store(newl.getName(), newl)

	return newl
}

func (g *LightGroup) processInventory(inv *nexus_client.InventoryInventory) {

	for lightId := 0; lightId < g.numLights; lightId++ {

		g.getCachedLight(lightId)

		invLight, lightErr := inv.AddLights(g.ctxt, &baselighttrafficlightcomv1.Light{
			ObjectMeta: v1.ObjectMeta{
				Name: g.getName(lightId),
			},
			Spec: baselighttrafficlightcomv1.LightSpec{
				GroupNumber: g.groupId,
				LightNumber: lightId,
			},
		})

		if lightErr != nil && !errors.IsAlreadyExists(lightErr) {
			log.Panicf("failed to creation light %d in group %d with error %+v", lightId, g.groupId, lightErr)
		}

		log.Infof("Light %d in group %d created with object id %s", lightId, g.groupId, invLight.Name)
	}
}

func (l *Light) setCurrentLight(lightColor string) {
	l.lightMutex.Lock()
	defer l.lightMutex.Unlock()

	l.currentLight = lightColor
}

func (l *Light) getCurrentLight() string {
	l.lightMutex.RLock()
	defer l.lightMutex.RUnlock()

	return l.currentLight
}

func (l *Light) autoStatusReporting(stopStatusReporting chan bool) {

	timer := time.NewTimer(l.autoStatusReportingPeriodic)
	for {
		select {
		case <-timer.C:
			light, err := l.nexusClient.RootRoot().DesiredConfig().GetLights(l.ctxt, l.getName())

			status := basedesiredconfiglighttrafficlightcomv1.LightInfo{
				CurrentLight:         l.getCurrentLight(),
				LastStatusReportTime: time.Now().Format("2006-01-02 15:04:05.999999999 -0700 MST"),
			}

			if err == nil {
				currStatus, statusErr := light.GetStatus(l.ctxt, l.getName())
				if statusErr == nil {
					currStatus.Spec.AutoStatus = status
					statusErr = currStatus.Update(l.ctxt)
					if statusErr != nil {
						log.Panicf("unable to status auto status with error %v", statusErr)
					}
				} else if nexus_client.IsChildNotFound(statusErr) {
					_, statusCreateErr := light.AddStatus(l.ctxt,
						&basedesiredconfiglighttrafficlightcomv1.LightStatus{
							ObjectMeta: v1.ObjectMeta{
								Name: l.getName(),
							},
							Spec: basedesiredconfiglighttrafficlightcomv1.LightStatusSpec{
								AutoStatus: status,
							},
						})
					if statusCreateErr != nil {
						log.Panicf("error %v when creating status %s", statusCreateErr, l.getName())
					}
				}
			}
			// Restart the timer
			timer = time.NewTimer(l.autoStatusReportingPeriodic)
		case <-stopStatusReporting:
			return
		}
	}
}

func (l *Light) signal(config *basedesiredconfiglighttrafficlightcomv1.LightConfigSpec) {
	defer func() {
		l.signalStopped <- true
	}()

	// Figure out light durations.
	redDuration, yellowDuration, greenDuration := uint32(25), uint32(5), uint32(30)
	if config.LightDuration.RedLightSeconds > 0 {
		redDuration = config.LightDuration.RedLightSeconds
	}
	if config.LightDuration.YellowLightSeconds > 0 {
		yellowDuration = config.LightDuration.YellowLightSeconds
	}
	if config.LightDuration.GreenLightSeconds > 0 {
		greenDuration = config.LightDuration.GreenLightSeconds
	}

	if config.AutoReportStatus {
		stopStatusReporting := make(chan bool)
		defer func() {
			stopStatusReporting <- true
		}()
		go l.autoStatusReporting(stopStatusReporting)
	}

	if config.DesiredLightAll != "" {
		if config.DesiredLightAll != utils.RED && config.DesiredLightAll != utils.YELLOW && config.DesiredLightAll != utils.GREEN {
			log.Panicf("invalid config received for all light configuration")
		}
		log.Debugf("Light group-%d/id-%d configured to static light color: %s", config.GroupID, config.LightID, config.DesiredLightAll)
		for {
			select {
			case <-l.signalStop:
				log.Infof("recieved signal stop request")
				return
			}
		}
	}

	timer := time.NewTimer(time.Duration(time.Duration(redDuration) * time.Second))
	l.setCurrentLight(utils.RED)
	log.Debugf("Light group-%v/id-%v set to: %s", config.GroupID, config.LightID, l.getCurrentLight())
	for {
		log.Debugf("Light group-%v/id-%v set to: %s", config.GroupID, config.LightID, l.getCurrentLight())
		select {
		case <-timer.C:
			if l.getCurrentLight() == utils.RED {
				timer = time.NewTimer(time.Duration(utils.RandRange(1, greenDuration)) * time.Second)
				l.setCurrentLight(utils.GREEN)
			} else if l.getCurrentLight() == utils.GREEN {
				timer = time.NewTimer(time.Duration(utils.RandRange(1, yellowDuration)) * time.Second)
				l.setCurrentLight(utils.YELLOW)
			} else {
				timer = time.NewTimer(time.Duration(utils.RandRange(1, redDuration)) * time.Second)
				l.setCurrentLight(utils.RED)
			}

			if l.autoStatusReportingOnUpdate {
				// set status
				l.updateLightChangeStatus()
			}
		case <-l.signalStop:
			log.Infof("recieved signal stop request")
			return
		}
	}
}

func (l *Light) updateLightChangeStatus() {
	light, err := l.nexusClient.RootRoot().DesiredConfig().GetLights(l.ctxt, l.getName())
	if err != nil {
		log.Panicf("unable to light %s with error %v", l.getName(), err)
	}

	currStatus, statusErr := light.GetStatus(l.ctxt, l.getName())
	if statusErr == nil {
		currStatus.Spec.AutoStatus.CurrentLight = l.getCurrentLight()
	}
	currStatus.Update(l.ctxt)
}

func (g *LightGroup) updateLightBroadcastLatency(light *nexus_client.DesiredconfiglightLightConfig, timestamp time.Time) {

	configTime := utils.GetTime(light.Spec.ConfigInfo.Timestamp)
	latency := timestamp.Sub(configTime).Milliseconds()

	configStatus := basedesiredconfiglighttrafficlightcomv1.LightConfigStatus{
		Timestamp: utils.GetTimeString(timestamp),
		Version:   light.ResourceVersion,
	}

	currStatus, statusErr := light.GetStatus(g.ctxt, g.getName(light.Spec.LightID))
	if statusErr == nil {
		if currStatus.Spec.BroadcastLatency.Min == 0 ||
			currStatus.Spec.BroadcastLatency.Min > latency {
			currStatus.Spec.BroadcastLatency.Min = latency
		}

		if currStatus.Spec.BroadcastLatency.Max == 0 ||
			currStatus.Spec.BroadcastLatency.Max < latency {
			currStatus.Spec.BroadcastLatency.Max = latency
		}

		if g.groupId == light.Spec.GroupID {
			currStatus.Spec.Config = configStatus
		}

	} else if nexus_client.IsChildNotFound(statusErr) {
		lightSpec := basedesiredconfiglighttrafficlightcomv1.LightStatusSpec{
			BroadcastLatency: basedesiredconfiglighttrafficlightcomv1.Latency{
				Min: latency,
				Max: latency,
			},
		}

		if g.groupId == light.Spec.GroupID {
			lightSpec.Config = configStatus
		}

		_, statusCreateErr := light.AddStatus(g.ctxt,
			&basedesiredconfiglighttrafficlightcomv1.LightStatus{
				ObjectMeta: v1.ObjectMeta{
					Name: g.getName(light.Spec.LightID),
				},
				Spec: lightSpec,
			})

		if statusCreateErr != nil {
			log.Panicf("error %v when creating status %s", statusCreateErr, g.getName(light.Spec.LightID))
		}
	}
}

func (g *LightGroup) updateLightUnicastLatency(light *nexus_client.DesiredconfiglightLightConfig, timestamp time.Time) {

	if g.groupId != light.Spec.GroupID {
		return
	}

	configTime := utils.GetTime(light.Spec.ConfigInfo.Timestamp)
	latency := timestamp.Sub(configTime).Milliseconds()

	configStatus := basedesiredconfiglighttrafficlightcomv1.LightConfigStatus{
		Timestamp: utils.GetTimeString(timestamp),
		Version:   light.ResourceVersion,
	}

	currStatus, statusErr := light.GetStatus(g.ctxt, g.getName(light.Spec.LightID))
	if statusErr == nil {
		if currStatus.Spec.UnicastLatency.Min == 0 ||
			currStatus.Spec.UnicastLatency.Min > latency {
			currStatus.Spec.UnicastLatency.Min = latency
		}

		if currStatus.Spec.UnicastLatency.Max == 0 ||
			currStatus.Spec.UnicastLatency.Max < latency {
			currStatus.Spec.UnicastLatency.Max = latency
		}
		currStatus.Spec.Config = configStatus

	} else if nexus_client.IsChildNotFound(statusErr) {
		lightSpec := basedesiredconfiglighttrafficlightcomv1.LightStatusSpec{
			BroadcastLatency: basedesiredconfiglighttrafficlightcomv1.Latency{
				Min: latency,
				Max: latency,
			},
		}
		lightSpec.Config = configStatus

		_, statusCreateErr := light.AddStatus(g.ctxt,
			&basedesiredconfiglighttrafficlightcomv1.LightStatus{
				ObjectMeta: v1.ObjectMeta{
					Name: g.getName(light.Spec.LightID),
				},
				Spec: lightSpec,
			})

		if statusCreateErr != nil {
			log.Panicf("error %v when creating status %s", statusCreateErr, g.getName(light.Spec.LightID))
		}
	}
}

func (g *LightGroup) updateLightStatus(light *nexus_client.DesiredconfiglightLightConfig, timestamp time.Time) {
	configStatus := basedesiredconfiglighttrafficlightcomv1.LightConfigStatus{
		Timestamp: timestamp.Format("2006-01-02 15:04:05.999999999 -0700 MST"),
		Version:   light.ResourceVersion,
	}

	currStatus, statusErr := light.GetStatus(g.ctxt, g.getName(light.Spec.LightID))
	if statusErr == nil {
		currStatus.Spec.Config = configStatus

		statusErr = currStatus.Update(g.ctxt)
		if statusErr != nil {
			log.Panicf("unable to status auto status with error %v", statusErr)
		}
	} else if nexus_client.IsChildNotFound(statusErr) {
		_, statusCreateErr := light.AddStatus(g.ctxt,
			&basedesiredconfiglighttrafficlightcomv1.LightStatus{
				ObjectMeta: v1.ObjectMeta{
					Name: g.getName(light.Spec.LightID),
				},
				Spec: basedesiredconfiglighttrafficlightcomv1.LightStatusSpec{
					Config: configStatus,
				},
			})
		if statusCreateErr != nil {
			log.Panicf("error %v when creating status %s", statusCreateErr, g.getName(light.Spec.LightID))
		}
	}
}

func (g *LightGroup) processAddLight(dcLight *nexus_client.DesiredconfiglightLightConfig) {
	timestamp := time.Now()

	log.Debugf("processAddLight: object %s", dcLight.DisplayName())
	if isPerfTestInProgress() {
		mode := getPerfTestType()
		log.Debugf("processAddLight: perfTestInProgress with mode %s", mode)
		if mode == utils.BROADCAST || mode == utils.BOOT {
			// g.updateLightBroadcastLatency(dcLight, timestamp)
			perf.IncTestObjectCount()
			perf.SetObjectTime(timestamp)
		} else if mode == utils.UNICAST {
			g.updateLightBroadcastLatency(dcLight, timestamp)
		} else if mode == utils.BOOT {
			perf.IncTestObjectCount()
			perf.SetObjectTime(timestamp)
		}

		return
	}

	if g.groupId != dcLight.Spec.GroupID {
		// This event is not meant for this light group.
		return
	}

	g.updateLightStatus(dcLight, timestamp)

	l := g.getCachedLight(dcLight.Spec.LightID)

	// Check if the light needs to be restarted.
	if l.signalStop != nil && l.signalStopped != nil {
		l.signalStop <- true
		<-l.signalStopped
	}

	go func() {
		l.signalStop = make(chan bool)
		l.signalStopped = make(chan bool)
		l.signal(&dcLight.Spec)
	}()
}

func (g *LightGroup) processUpdateLight(oldDCLight, newDCLight *nexus_client.DesiredconfiglightLightConfig) {
	timestamp := time.Now()

	log.Debugf("processUpdateLight: object %s", newDCLight.DisplayName())
	if isPerfTestInProgress() {
		mode := getPerfTestType()
		log.Debugf("processAddLight: perfTestInProgress with mode %s", mode)
		if mode == utils.BROADCAST || mode == utils.BOOT {
			// g.updateLightBroadcastLatency(newDCLight, timestamp)
			perf.IncTestObjectCount()
			perf.SetObjectTime(timestamp)
		} else if mode == utils.UNICAST {
			g.updateLightBroadcastLatency(newDCLight, timestamp)
		}

		return
	}

	if g.groupId != newDCLight.Spec.GroupID {
		// This event is not meant for this light group.
		return
	}
	log.Debugf("Light group-%v/id-%v update received", newDCLight.Spec.GroupID, newDCLight.Spec.LightID)

	restartSignal := false

	if oldDCLight.Spec.DesiredLightAll != newDCLight.Spec.DesiredLightAll ||
		oldDCLight.Spec.LightDuration != newDCLight.Spec.LightDuration ||
		oldDCLight.Spec.AutoReportStatus != newDCLight.Spec.AutoReportStatus {
		restartSignal = true
	}

	l := g.getCachedLight(newDCLight.Spec.LightID)

	g.updateLightStatus(newDCLight, timestamp)

	if restartSignal {

		// Check if the light needs to be restarted.
		if l.signalStop != nil && l.signalStopped != nil {
			l.signalStop <- true
			<-l.signalStopped
		}

		go func() {
			l.signalStop = make(chan bool)
			l.signalStopped = make(chan bool)
			l.signal(&newDCLight.Spec)
		}()
	}

	/* 	if oldDCLight.Spec.LastManualStatusRequestedTime != newDCLight.Spec.LastManualStatusRequestedTime {

		status := basedesiredconfiglighttrafficlightcomv1.LightInfo{
			CurrentLight:         l.getCurrentLight(),
			LastStatusReportTime: time.Now().String(),
		}

		currStatus, statusErr := newDCLight.GetStatus(l.ctxt, l.getName())
		if statusErr == nil {
			currStatus.Spec.ManualStatus = status
			statusErr = currStatus.Update(l.ctxt)
			if statusErr != nil {
				log.Panicf("unable to status auto status with error %v", statusErr)
			}
		} else if nexus_client.IsChildNotFound(statusErr) {
			_, statusCreateErr := newDCLight.AddStatus(l.ctxt,
				&basedesiredconfiglighttrafficlightcomv1.LightStatus{
					ObjectMeta: v1.ObjectMeta{
						Name: l.getName(),
					},
					Spec: basedesiredconfiglighttrafficlightcomv1.LightStatusSpec{
						ManualStatus: status,
					},
				})
			if statusCreateErr != nil {
				log.Panicf("error %v when creating status %s", statusCreateErr, l.getName())
			}
		}

	} */

}

func (g *LightGroup) processPerfTestInit(mode string, inprogress bool) {
	perfTestMutex.Lock()
	defer perfTestMutex.Unlock()

	perfTestType = mode
	perfTestInProgress = inprogress
}

func (g *LightGroup) processPerfTestAdd(obj *nexus_client.PerftestPerfTest) {
	perfTestMutex.Lock()
	defer perfTestMutex.Unlock()

	if perfTestInProgress != obj.Spec.TestInProgress {
		// Test in progress status has flipped from the last time we read
		os.Exit(0)
	}
	perfTestType = obj.Spec.Mode
}

func isPerfTestInProgress() bool {
	perfTestMutex.RLock()
	defer perfTestMutex.RUnlock()
	return perfTestInProgress
}

func getPerfTestType() string {
	perfTestMutex.RLock()
	defer perfTestMutex.RUnlock()

	return perfTestType
}

func (g *LightGroup) processPerfTestUpdate(old, new *nexus_client.PerftestPerfTest) {
	if old.Spec.TestInProgress != new.Spec.TestInProgress {
		// Die, so we can come back in test mode
		os.Exit(0)
	}

	if new.Spec.TestInProgress && old.Spec.Mode != new.Spec.Mode {
		// Die, so we can come back in the new test mode
		os.Exit(0)
	}

	g.processPerfTestAdd(new)
}

func (g *LightGroup) getBroadcastInfo(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(perf.GetBroadcastInfo())
}

func (g *LightGroup) getBootInfo(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(perf.GetBootInfo())
}

func (g *LightGroup) clear(w http.ResponseWriter, req *http.Request) {
	perf.Reset()
}

func LightGroupInit(client *nexus_client.Clientset, groupId, numLights int, logger *logrus.Logger, ctxt context.Context, doneFn context.CancelFunc) {
	if TrafficLightGroup == nil {
		TrafficLightGroup = &LightGroup{
			nexusClient: client,
			groupId:     groupId,
			ctxt:        ctxt,
			doneFn:      doneFn,
			numLights:   numLights,
		}
	}
	log = logger

	// Get before subscribe so we can read from the DB directly at boot time.
	testObj, err := TrafficLightGroup.nexusClient.RootRoot().Config().GetPerfTest(ctxt)
	if err == nil {
		TrafficLightGroup.processPerfTestInit(testObj.Spec.Mode, testObj.Spec.TestInProgress)
	} else if nexus_client.IsNotFound(err) {
		TrafficLightGroup.processPerfTestInit("", false)
	} else {
		log.Panic("error while trying to get the test object: %v", err)
	}

	TrafficLightGroup.nexusClient.RootRoot().Subscribe()
	TrafficLightGroup.nexusClient.RootRoot().Config().Subscribe()
	TrafficLightGroup.nexusClient.RootRoot().Config().Group("*").Subscribe()
	TrafficLightGroup.nexusClient.RootRoot().Config().PerfTest().Subscribe()
	TrafficLightGroup.nexusClient.RootRoot().Config().PerfTest().RegisterAddCallback(TrafficLightGroup.processPerfTestAdd)
	TrafficLightGroup.nexusClient.RootRoot().Config().PerfTest().RegisterUpdateCallback(TrafficLightGroup.processPerfTestUpdate)

	TrafficLightGroup.nexusClient.RootRoot().DesiredConfig().Subscribe()
	TrafficLightGroup.nexusClient.RootRoot().DesiredConfig().Lights("*").Subscribe()
	// TrafficLightGroup.nexusClient.RootRoot().DesiredConfig().Lights("*").Status("*").Subscribe()

	TrafficLightGroup.nexusClient.RootRoot().Inventory().Subscribe()
	TrafficLightGroup.nexusClient.RootRoot().Inventory().Lights("*").Subscribe()
	perf.SetObjectSubscribeTime(time.Now())
	TrafficLightGroup.nexusClient.RootRoot().DesiredConfig().Lights("*").RegisterAddCallback(TrafficLightGroup.processAddLight)
	TrafficLightGroup.nexusClient.RootRoot().DesiredConfig().Lights("*").RegisterUpdateCallback(TrafficLightGroup.processUpdateLight)

	TrafficLightGroup.nexusClient.RootRoot().Inventory().RegisterAddCallback(TrafficLightGroup.processInventory)

	http.HandleFunc("/bootinfo", TrafficLightGroup.getBootInfo)
	http.HandleFunc("/broadcastinfo", TrafficLightGroup.getBroadcastInfo)
	http.HandleFunc("/clear", TrafficLightGroup.clear)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
