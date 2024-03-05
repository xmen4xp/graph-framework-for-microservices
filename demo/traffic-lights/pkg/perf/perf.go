package perf

import (
	"sync"
	"time"
	"trafficlightdemo/pkg/utils"
)

var (
	perfTestObjectCount     int64
	perfObjectSubscribeTime time.Time
	perfTestFirstObjectTime time.Time
	perfTestLastObjectTime  time.Time

	perfMutex sync.RWMutex
)

type PerfInfo struct {
	ObjectCount           int64
	LatencyInMilliSeconds int64
	SubscriptionTime      string
	FirstObjectTime       string
	LastObjectTime        string
}

func Reset() {
	perfMutex.Lock()
	defer perfMutex.Unlock()

	perfTestObjectCount = 0
	perfObjectSubscribeTime = time.Time{}
	perfTestFirstObjectTime = time.Time{}
	perfTestLastObjectTime = time.Time{}
}

func IncTestObjectCount() {
	perfMutex.Lock()
	defer perfMutex.Unlock()
	perfTestObjectCount++
}

func GetTestObjectCount() int64 {
	perfMutex.RLock()
	defer perfMutex.RUnlock()
	return perfTestObjectCount
}

func SetObjectTime(ts time.Time) {
	perfMutex.Lock()
	defer perfMutex.Unlock()

	if perfTestFirstObjectTime.IsZero() {
		perfTestFirstObjectTime = ts
		perfTestLastObjectTime = ts
	} else {
		perfTestLastObjectTime = ts
	}
}

func SetObjectSubscribeTime(ts time.Time) {
	perfMutex.Lock()
	defer perfMutex.Unlock()
	perfObjectSubscribeTime = ts
}

func GetFirstObjectTime() time.Time {
	perfMutex.RLock()
	defer perfMutex.RUnlock()
	return perfTestFirstObjectTime
}

func GetLastObjectTime() time.Time {
	perfMutex.RLock()
	defer perfMutex.RUnlock()
	return perfTestLastObjectTime
}

func GetBootInfo() *PerfInfo {
	perfMutex.RLock()
	defer perfMutex.RUnlock()

	var latency int64
	if !perfObjectSubscribeTime.IsZero() && !perfTestLastObjectTime.IsZero() {
		latency = perfTestLastObjectTime.Sub(perfObjectSubscribeTime).Milliseconds()
	}
	return &PerfInfo{
		ObjectCount:           perfTestObjectCount,
		LatencyInMilliSeconds: latency,
		SubscriptionTime:      utils.GetTimeString(perfObjectSubscribeTime),
		FirstObjectTime:       utils.GetTimeString(perfTestFirstObjectTime),
		LastObjectTime:        utils.GetTimeString(perfTestLastObjectTime),
	}
}

func GetBroadcastInfo() *PerfInfo {
	perfMutex.RLock()
	defer perfMutex.RUnlock()

	return &PerfInfo{
		ObjectCount:     perfTestObjectCount,
		FirstObjectTime: utils.GetTimeString(perfTestFirstObjectTime),
		LastObjectTime:  utils.GetTimeString(perfTestLastObjectTime),
	}
}
