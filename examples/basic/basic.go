package main

import (
	"fmt"
	"log"
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

	go func() {
		tick := time.NewTicker(1 * time.Second)
		defer tick.Stop()
		for range tick.C {
			requests, _ := statistics.GetCounter("requests")
			errors, _ := statistics.GetCounter("errors")
			startedAt, _ := statistics.GetStatic("startedAt")
			rps, _ := statistics.GetDynamic("rps")

			data := fmt.Sprintf("Requests: [%d/%d] StartedAt: %s RPS: %s", requests, errors, clistats.String(startedAt), clistats.String(rps(statistics)))
			log.Printf("%s\r\n", data)
		}
	}()

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
