package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/projectdiscovery/clistats"
)

func main() {
	statistics, err := clistats.New()
	if err != nil {
		log.Fatal(err)
	}
	statistics.AddCounter("requests", 0)
	statistics.AddCounter("errors", 0)
	statistics.AddStatic("startedAt", time.Now())
	statistics.AddDynamic("rps", clistats.NewRequestsPerSecondCallback(clistats.RequestPerSecondCallbackOptions{
		StartTimeFieldID:  "startedAt",
		RequestsCounterID: "requests",
	}))
	printMutex := &sync.Mutex{}

	statistics.Start(func(stats clistats.StatisticsClient) {
		requests, _ := stats.GetCounter("requests")
		errors, _ := stats.GetCounter("errors")
		startedAt, _ := stats.GetStatic("startedAt")
		rps, _ := stats.GetDynamic("rps")

		data := fmt.Sprintf("Requests: [%d/%d] StartedAt: %s RPS: %s", requests, errors, clistats.String(startedAt), clistats.String(rps(stats)))
		printMutex.Lock()
		log.Printf("%s\r\n", data)
		printMutex.Unlock()
	}, 1*time.Second)

	statistics.IncrementCounter("requests", 1)
	time.Sleep(3 * time.Second)
	statistics.IncrementCounter("requests", 1)
	statistics.IncrementCounter("requests", 1)
	statistics.IncrementCounter("requests", 1)
	statistics.IncrementCounter("requests", 1)
	statistics.IncrementCounter("errors", 1)
	statistics.IncrementCounter("requests", 1)
	statistics.IncrementCounter("requests", 1)
	time.Sleep(3 * time.Second)
	statistics.IncrementCounter("requests", 1)
	statistics.IncrementCounter("requests", 1)
	statistics.IncrementCounter("requests", 1)
}
