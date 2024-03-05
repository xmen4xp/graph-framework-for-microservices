package utils

import (
	"log"
	"math/rand"
	"time"
)

const (
	GROUP_DNS_PREFIX   = "light-connector-"
	GROUP_SERVICE_NAME = "light-connector"
	GROUP_PORT_NUM     = "8080"
)

const (
	RED    = "RED"
	YELLOW = "YELLOW"
	GREEN  = "GREEN"
)

const (
	BROADCAST = "BROADCAST"
	UNICAST   = "UNICAST"
	BOOT      = "BOOT"
)

func GetLightChar(color string) string {
	if color == RED {
		return "R"
	} else if color == GREEN {
		return "G"
	} else if color == YELLOW {
		return "Y"
	}
	return "X"
}

func RandRange(min, max uint32) uint32 {
	return uint32(rand.Intn(int(max-min)) + int(min))
}

func GetTimeString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.999999999 -0700 MST")
}

func GetTime(tStr string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", tStr)
	if err != nil {
		log.Panic("determining time from time string %s failed with error %v", tStr, err)
	}
	return t
}
