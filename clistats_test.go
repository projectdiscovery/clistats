package clistats

import "time"

func ExampleDynamicCallback_rps(client StatisticsClient) interface{} {
	requests, _ := client.GetCounter("requests")
	start, _ := client.GetStatic("startTime")
	startTime := start.(time.Time)

	return float64(requests) / time.Since(startTime).Seconds()
}

func ExampleDynamicCallback_Elapsedtime(client StatisticsClient) interface{} {
	start, _ := client.GetStatic("startTime")
	startTime := start.(time.Time)

	return time.Since(startTime).Seconds()
}
