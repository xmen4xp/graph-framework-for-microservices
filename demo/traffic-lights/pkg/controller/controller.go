package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"
	nexus_client "trafficlight/build/nexus-client"
	"trafficlightdemo/pkg/perf"
	"trafficlightdemo/pkg/utils"

	commonlighttrafficlightcomv1 "trafficlight/build/apis/common.trafficlight.com/v1"
	basedesiredconfiglighttrafficlightcomv1 "trafficlight/build/apis/desiredconfiglight.trafficlight.com/v1"

	"math/rand"
	"net/http"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type controller struct {
	nexusClient         *nexus_client.Clientset
	groupId             uint64
	ctxt                context.Context
	doneFn              context.CancelFunc
	signalStop          chan bool
	signalStopped       chan bool
	configMutex         sync.Mutex
	lastConfigVersion   uint64
	perfTest            bool
	maxBroadcastLatency time.Duration
	minBroadcastLatency time.Duration
}

const (
	RED    = "RED"
	YELLOW = "YELLOW"
	GREEN  = "GREEN"
)

var lightController *controller
var log *logrus.Logger

func (c *controller) processLight(light *nexus_client.LightLight) {

	desiredConfig, desiredConfigErr := c.nexusClient.RootRoot().GetDesiredConfig(c.ctxt)
	if desiredConfigErr != nil {
		log.Errorf("unable to get desired config object for light %s", light.DisplayName())
		return
	}

	desiredConfigExists := false
	var lightConfig *basedesiredconfiglighttrafficlightcomv1.LightConfig
	desiredLightConfig, err := desiredConfig.GetLights(c.ctxt, light.DisplayName())
	if nexus_client.IsChildNotFound(err) {
		lightConfig = &basedesiredconfiglighttrafficlightcomv1.LightConfig{
			ObjectMeta: v1.ObjectMeta{
				Name: light.DisplayName(),
			},
			Spec: basedesiredconfiglighttrafficlightcomv1.LightConfigSpec{
				GroupID:          light.Spec.GroupNumber,
				LightID:          light.Spec.LightNumber,
				LightDuration:    commonlighttrafficlightcomv1.LightDuration{},
				AutoReportStatus: true,
				OnboardTime:      time.Now().Format("2006-01-02 15:04:05.999999999 -0700 MST"),
			},
		}
	} else {
		lightConfig = desiredLightConfig.LightConfig
		desiredConfigExists = true
	}

	config, err := c.nexusClient.RootRoot().GetConfig(c.ctxt)
	if err == nil && config != nil {
		lightConfig.Spec.LightDuration = config.Spec.LightDuration
		lightConfig.Spec.AutoReportStatus = config.Spec.AutoReportStatus
		lightConfig.Spec.DesiredLightAll = config.Spec.DesiredLightAll
		lightConfig.Spec.ConfigInfo.Timestamp = time.Now().Format("2006-01-02 15:04:05.999999999 -0700 MST")
		lightConfig.Spec.ConfigInfo.Version = config.ResourceVersion
	}

	if desiredConfigExists {
		desiredLightConfig.Update(c.ctxt)
	} else {
		desiredConfig.AddLights(c.ctxt, lightConfig)
	}

	log.Infof("light Add event received for %s", light.DisplayName())
}

func (c *controller) processConfig(config *nexus_client.ConfigConfig) {
	currVersion, err := strconv.ParseUint(config.ResourceVersion, 10, 64)
	if err != nil {
		log.Panicf("error deducing resource version from version string %s in config object %s", config.ResourceVersion, config.DisplayName())
	}

	if currVersion <= c.lastConfigVersion {
		log.Debugf("ignoring stale object update with version: %s. Last processed version %v", config.ResourceVersion, c.lastConfigVersion)
	}

	c.configMutex.Lock()
	defer c.configMutex.Unlock()

	desiredConfig, desiredConfigErr := c.nexusClient.RootRoot().GetDesiredConfig(c.ctxt)
	if desiredConfigErr != nil {
		log.Errorf("unable to get desired config object handle")
		return
	}

	lightsIter := desiredConfig.GetAllLightsIter(c.ctxt)
	for light, _ := lightsIter.Next(c.ctxt); light != nil; light, _ = lightsIter.Next(c.ctxt) {
		light.Spec.LightDuration = config.Spec.LightDuration
		light.Spec.AutoReportStatus = config.Spec.AutoReportStatus
		light.Spec.DesiredLightAll = config.Spec.DesiredLightAll
		light.Spec.ConfigInfo.Timestamp = time.Now().Format("2006-01-02 15:04:05.999999999 -0700 MST")
		light.Spec.ConfigInfo.Version = config.ResourceVersion

		// Update the light
		light.Update(c.ctxt)
	}

	log.Infof("config %s with version %s processed successfully", config.DisplayName(), config.ResourceVersion)
}

func (c *controller) processConfigUpdate(old, new *nexus_client.ConfigConfig) {

	if old.Spec.LightDuration != new.Spec.LightDuration ||
		old.Spec.AutoReportStatus != new.Spec.AutoReportStatus ||
		old.Spec.DesiredLightAll != new.Spec.DesiredLightAll ||
		c.perfTest == true {
		log.Debugf("processing update to config %s resource version %s", new.DisplayName(), new.ResourceVersion)
		c.processConfig(new)
	} else {
		log.Debugf("Ignoring unconsequential update to config %s resource version %s", new.DisplayName(), new.ResourceVersion)
	}

}

func ControllerInit(nexus_client *nexus_client.Clientset, groupId uint64, logger *logrus.Logger, ctxt context.Context, doneFn context.CancelFunc) {

	if lightController == nil {
		lightController = &controller{
			nexusClient: nexus_client,
			groupId:     groupId,
			ctxt:        ctxt,
			doneFn:      doneFn,
			configMutex: sync.Mutex{},
		}
		log = logger
	}

	lightController.nexusClient.RootRoot().Subscribe()
	lightController.nexusClient.RootRoot().Config().Subscribe()
	lightController.nexusClient.RootRoot().Config().Group("*").Subscribe()

	lightController.nexusClient.RootRoot().DesiredConfig().Subscribe()
	lightController.nexusClient.RootRoot().DesiredConfig().Lights("*").Subscribe()
	lightController.nexusClient.RootRoot().DesiredConfig().Lights("*").Status("*").Subscribe()

	lightController.nexusClient.RootRoot().Inventory().Subscribe()
	lightController.nexusClient.RootRoot().Inventory().Lights("*").Subscribe()

	lightController.nexusClient.RootRoot().Inventory().Lights("*").RegisterAddCallback(lightController.processLight)
	lightController.nexusClient.RootRoot().Config().RegisterAddCallback(lightController.processConfig)
	lightController.nexusClient.RootRoot().Config().RegisterUpdateCallback(lightController.processConfigUpdate)

	http.HandleFunc("/broadcast", lightController.broadcastLatencyCalcHandler)
	http.HandleFunc("/lightcoldboot", lightController.lightColdBootLatencyCalcHandler)
	http.HandleFunc("/systemstatus", lightController.systemStatus)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func randRange(min, max int) uint32 {
	return uint32(rand.Intn(max-min) + min)
}

func printBroadcastHeader(outputString string, numGroups, numLights int, min, max int64) string {

	outputString = fmt.Sprintf("%s\n", outputString)
	outputString = fmt.Sprintf("%s Total Groups:   %3d\n", outputString, numGroups)
	outputString = fmt.Sprintf("%s Total Lights:   %3d\n", outputString, numLights)
	outputString = fmt.Sprintf("%s Max Latency :   %3d ms\n", outputString, max)
	outputString = fmt.Sprintf("%s Min Latency :   %3d ms\n", outputString, min)

	outputString = fmt.Sprintf("%s\n", outputString)

	outputString = fmt.Sprintf("%s\n------------------------------------------------", outputString)
	outputString = fmt.Sprintf("%s\n     Group    |  ObjectCount   | Latency (ms)   ", outputString)
	outputString = fmt.Sprintf("%s\n------------------------------------------------", outputString)
	return outputString
}

func printBroadcastGroup(outputString string, groupId int, objCount, latency int64) string {
	outputString = fmt.Sprintf("%s\n     %5d    |     %5d     |    %4d   ",
		outputString, groupId, objCount, latency)
	return outputString
}

func (c *controller) broadcastLatencyCalcHandler(w http.ResponseWriter, req *http.Request) {

	type broadcastLightPerfData struct {
		configTime time.Time
		latency    int64
		bcPerf     perf.PerfInfo
	}

	c.perfTest = true

	io.WriteString(w, "BroadcastLatencyCalcHandler!\n")

	config, configErr := c.nexusClient.RootRoot().GetConfig(c.ctxt)
	if configErr != nil {
		log.Panicf("failure to read config object")
	}

	desiredConfig, dcerr := c.nexusClient.RootRoot().GetDesiredConfig(c.ctxt)
	if dcerr != nil {
		log.Panicf("failure to read desired config object")
	}

	// Clear counters from prev runs
	lightsIter := desiredConfig.GetAllLightsIter(c.ctxt)
	for light, _ := lightsIter.Next(c.ctxt); light != nil; light, _ = lightsIter.Next(c.ctxt) {
		requestURL := fmt.Sprintf("http://%s%d.%s:%s/clear",
			utils.GROUP_DNS_PREFIX, light.Spec.GroupID,
			utils.GROUP_SERVICE_NAME, utils.GROUP_PORT_NUM)

		res, err := http.Get(requestURL)
		if err != nil {
			log.Panicf("failed to get response from %s with error %v", requestURL, err)
		}
		defer res.Body.Close()
	}

	// Trigger broadcast event
	config.Spec.LightDuration.GreenLightSeconds = randRange(1, 30)
	config.Spec.LightDuration.YellowLightSeconds = randRange(1, 10)
	config.Spec.LightDuration.RedLightSeconds = randRange(1, 25)
	log.Infof("Current config version: %s", config.ResourceVersion)
	err := config.Update(c.ctxt)
	if err != nil {
		log.Panicf("failure to update config object")
	}

	log.Infof("New config version: %s", config.ResourceVersion)
	log.Infof("Waiting for the configuration to propagate...")

	// Wait for event propgation
	time.Sleep(60 * time.Second)

	// Query perf data
	lightsIter = desiredConfig.GetAllLightsIter(c.ctxt)
	groupPerfInfo := make(map[int]broadcastLightPerfData)

	numLights := 0
	for light, _ := lightsIter.Next(c.ctxt); light != nil; light, _ = lightsIter.Next(c.ctxt) {
		if light.Spec.ConfigInfo.Version != config.ResourceVersion {
			log.Errorf("config propagation has not completed for light config %s / %s. LightConfig version %s",
				light.DisplayName(), light.Name, light.Spec.ConfigInfo.Version)
			continue
		}

		if _, ok := groupPerfInfo[light.Spec.GroupID]; !ok {
			requestURL := fmt.Sprintf("http://%s%d.%s:%s/broadcastinfo",
				utils.GROUP_DNS_PREFIX, light.Spec.GroupID,
				utils.GROUP_SERVICE_NAME, utils.GROUP_PORT_NUM)

			res, err := http.Get(requestURL)
			if err != nil {
				log.Panicf("failed to get response from %s with error %v", requestURL, err)
			}
			defer res.Body.Close()

			decoder := json.NewDecoder(res.Body)
			var data perf.PerfInfo
			err = decoder.Decode(&data)
			if err != nil {
				log.Panicf("failed to decode response from %s with error %v", requestURL, err)
			}

			configTime := utils.GetTime(light.Spec.ConfigInfo.Timestamp)
			lastObjTime := utils.GetTime(data.LastObjectTime)
			latency := lastObjTime.Sub(utils.GetTime(light.Spec.ConfigInfo.Timestamp)).Milliseconds()
			groupPerfInfo[light.Spec.GroupID] = broadcastLightPerfData{
				configTime: configTime,
				latency:    latency,
				bcPerf:     data,
			}
		}
		numLights++
	}

	/* 		lightConfigTime, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", light.Spec.ConfigInfo.Timestamp)
	   		if err != nil {
	   			log.Errorf("failed to convert ligt config time string %s to duration due to error %v", light.Spec.ConfigInfo.Timestamp, err)
	   		}

	   		status, statusErr := light.GetStatus(c.ctxt, light.DisplayName())
	   		if statusErr != nil {
	   			log.Panicf("status not found for light config %s", light.DisplayName())
	   		}

	   		statusTime, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", status.Spec.Config.Timestamp)
	   		if err != nil {
	   			log.Errorf("failed to convert status time string %s to duration due to error %v", status.Spec.Config.Timestamp, err)
	   		}
	   		log.Infof("Light Config Time: %v  Status Time: %v latency: %v ms", lightConfigTime, statusTime, statusTime.Sub(lightConfigTime).Milliseconds())

	   		currLatency := statusTime.Sub(lightConfigTime)
	   		if c.minBroadcastLatency == 0 {
	   			c.minBroadcastLatency = currLatency
	   		} else if currLatency < c.minBroadcastLatency {
	   			c.minBroadcastLatency = currLatency
	   		}

	   		if currLatency > c.maxBroadcastLatency {
	   			c.maxBroadcastLatency = currLatency
	   		} */

	var minLatency, maxLatency int64
	outputString := ""
	for group := 0; group < len(groupPerfInfo); group++ {
		if perf, ok := groupPerfInfo[group]; ok {
			if minLatency == 0 {
				minLatency = perf.latency
				maxLatency = perf.latency
			} else if perf.latency < minLatency {
				minLatency = perf.latency
			} else if perf.latency > maxLatency {
				maxLatency = perf.latency
			}
			outputString = printBroadcastGroup(outputString, group, perf.bcPerf.ObjectCount, perf.latency)
		}
	}

	headerString := ""
	headerString = printBroadcastHeader(headerString, len(groupPerfInfo), numLights, minLatency, maxLatency)

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, headerString+outputString+"\n")
	c.perfTest = false
}

func printLightColdBootHeader(outputString string, numGroups, numLights int, min, max int64) string {

	outputString = fmt.Sprintf("%s\n", outputString)
	outputString = fmt.Sprintf("%s Total Groups:   %3d\n", outputString, numGroups)
	outputString = fmt.Sprintf("%s Total Lights:   %3d\n", outputString, numLights)
	outputString = fmt.Sprintf("%s Max Latency :   %3d ms\n", outputString, max)
	outputString = fmt.Sprintf("%s Min Latency :   %3d ms\n", outputString, min)

	outputString = fmt.Sprintf("%s\n", outputString)

	outputString = fmt.Sprintf("%s\n------------------------------------------------", outputString)
	outputString = fmt.Sprintf("%s\n     Group    |  ObjectCount   | Latency (ms)   ", outputString)
	outputString = fmt.Sprintf("%s\n------------------------------------------------", outputString)
	return outputString
}

func printLightColdBootGroup(outputString string, groupId int, objCount, latency int64) string {
	outputString = fmt.Sprintf("%s\n     %5d    |     %5d     |    %4d   ",
		outputString, groupId, objCount, latency)
	return outputString
}

func (c *controller) lightColdBootLatencyCalcHandler(w http.ResponseWriter, req *http.Request) {
	type coldBootLightPerfData struct {
		bootPerf perf.PerfInfo
	}

	c.perfTest = true
	desiredConfig, dcerr := c.nexusClient.RootRoot().GetDesiredConfig(c.ctxt)
	if dcerr != nil {
		log.Panicf("failure to read desired config object")
	}

	// Query cold boot data
	lightsIter := desiredConfig.GetAllLightsIter(c.ctxt)
	groupPerfInfo := make(map[int]coldBootLightPerfData)
	numLights := 0
	for light, _ := lightsIter.Next(c.ctxt); light != nil; light, _ = lightsIter.Next(c.ctxt) {
		if _, ok := groupPerfInfo[light.Spec.GroupID]; !ok {
			requestURL := fmt.Sprintf("http://%s%d.%s:%s/bootinfo",
				utils.GROUP_DNS_PREFIX, light.Spec.GroupID,
				utils.GROUP_SERVICE_NAME, utils.GROUP_PORT_NUM)

			res, err := http.Get(requestURL)
			if err != nil {
				log.Panicf("failed to get response from %s with error %v", requestURL, err)
			}
			defer res.Body.Close()

			decoder := json.NewDecoder(res.Body)
			var data perf.PerfInfo
			err = decoder.Decode(&data)
			if err != nil {
				log.Panicf("failed to decode response from %s with error %v", requestURL, err)
			}

			groupPerfInfo[light.Spec.GroupID] = coldBootLightPerfData{
				bootPerf: data,
			}
		}
		numLights++
	}

	var minLatency, maxLatency int64
	outputString := ""
	for group := 0; group < len(groupPerfInfo); group++ {
		if perf, ok := groupPerfInfo[group]; ok {
			if minLatency == 0 {
				minLatency = perf.bootPerf.LatencyInMilliSeconds
				maxLatency = perf.bootPerf.LatencyInMilliSeconds
			} else if perf.bootPerf.LatencyInMilliSeconds < minLatency {
				minLatency = perf.bootPerf.LatencyInMilliSeconds
			} else if perf.bootPerf.LatencyInMilliSeconds > maxLatency {
				maxLatency = perf.bootPerf.LatencyInMilliSeconds
			}
			outputString = printLightColdBootGroup(outputString, group, perf.bootPerf.ObjectCount, perf.bootPerf.LatencyInMilliSeconds)
		}
	}

	headerString := ""
	headerString = printLightColdBootHeader(headerString, len(groupPerfInfo), numLights, minLatency, maxLatency)

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, headerString+outputString+"\n")
	c.perfTest = false

	io.WriteString(w, "LightColdBootLatencyCalcHandler!\n")
}

type groupSummary struct {
	numRed      int
	numYellow   int
	numGreen    int
	lightStatus map[int]string
}

func printSummaryHeader(outputString string, numGroups, numLights int) string {

	outputString = fmt.Sprintf("%s\n", outputString)
	outputString = fmt.Sprintf("%s Total Groups:   %3d\n", outputString, numGroups)
	outputString = fmt.Sprintf("%s Total Lights:   %3d\n", outputString, numLights)
	outputString = fmt.Sprintf("%s\n", outputString)

	outputString = fmt.Sprintf("%s\n---------------------------------------------------------------", outputString)
	outputString = fmt.Sprintf("%s\n     Group    |  Total Lights |   # RED  | # YELLOW| # GREEN   ", outputString)
	outputString = fmt.Sprintf("%s\n---------------------------------------------------------------", outputString)
	return outputString
}

func printSummaryGroup(outputString string, groupId int, summary groupSummary) string {
	outputString = fmt.Sprintf("%s\n     %5d    |     %5d     |    %4d   |    %5d  |    %d     ",
		outputString, groupId, len(summary.lightStatus), summary.numRed, summary.numYellow, summary.numGreen)
	return outputString
}

func (c *controller) systemStatus(w http.ResponseWriter, req *http.Request) {

	numGroups := 200
	numLightsPerGroup := 500

	numGroupsString := req.URL.Query().Get("groups")
	log.Infof("numGroupsString: %s", numGroupsString)
	if numGroupsString != "" {
		if id, err := strconv.Atoi(numGroupsString); err == nil {
			numGroups = id
			log.Infof("numGroups: %d", numGroups)

		}
	}

	numLightsPerGroupString := req.URL.Query().Get("lights-per-group")
	log.Infof("numLightsPerGroupString: %s", numLightsPerGroupString)

	if numLightsPerGroupString != "" {
		if id, err := strconv.Atoi(numLightsPerGroupString); err == nil {
			numLightsPerGroup = id
			log.Infof("numGroups: %d", numLightsPerGroup)

		}
	}

	desiredConfig, dcerr := c.nexusClient.RootRoot().GetDesiredConfig(c.ctxt)
	if dcerr != nil {
		log.Panicf("failure to read desired config object")
	}

	systemStatus := make([]groupSummary, numGroups)
	for group := 0; group < numGroups; group++ {
		systemStatus[group] = groupSummary{
			lightStatus: make(map[int]string),
		}
	}
	groupsFound := make(map[int]struct{})

	numLights := 0
	lightsIter := desiredConfig.GetAllLightsIter(c.ctxt)
	for light, _ := lightsIter.Next(c.ctxt); light != nil; light, _ = lightsIter.Next(c.ctxt) {
		status, statusErr := light.GetStatus(c.ctxt, light.DisplayName())
		if statusErr != nil {
			log.Errorf("status not found for light config %s", light.DisplayName())
			continue
		}
		color := status.Spec.AutoStatus.CurrentLight
		if color == utils.RED {
			systemStatus[light.Spec.GroupID].numRed++
		} else if color == utils.YELLOW {
			systemStatus[light.Spec.GroupID].numYellow++
		} else if color == utils.GREEN {
			systemStatus[light.Spec.GroupID].numGreen++
		}
		systemStatus[light.Spec.GroupID].lightStatus[light.Spec.LightID] = utils.GetLightChar(color)
		numLights++
		groupsFound[light.Spec.GroupID] = struct{}{}
	}

	outputString := ""
	outputString = printSummaryHeader(outputString, len(groupsFound), numLights)
	for group := 0; group < numGroups; group++ {
		if len(systemStatus[group].lightStatus) > 0 {
			outputString = printSummaryGroup(outputString, group, systemStatus[group])
		}
	}
	outputString = fmt.Sprintf("%s\n", outputString)

	// outputString = fmt.Sprintf("%s-------------------\n", outputString)
	//outputString = fmt.Sprintf("%s # of Lights: %d\n", outputString, numLights)
	//outputString = fmt.Sprintf("%s Minimum Broadcast Latency: %d ms\n", outputString, c.minBroadcastLatency.Milliseconds())
	//outputString = fmt.Sprintf("%s Maximum Broadcast Latency: %d ms\n", outputString, c.maxBroadcastLatency.Milliseconds())

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, outputString)

}
