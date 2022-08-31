package clistats

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExampleDynamicCallbackRps(t *testing.T) {
	client, err := New()
	require.Nil(t, err)

	client.AddCounter("requests", 1000)
	client.AddStatic("startTime", time.Now())

	requests, ok := client.GetCounter("requests")
	require.True(t, ok)
	start, ok := client.GetStatic("startTime")
	require.True(t, ok)
	startTime := start.(time.Time)
	rps := float64(requests) / time.Since(startTime).Seconds()
	require.True(t, rps > 0)
}

func TestDynamicCallback_Elapsedtime(t *testing.T) {
	client, err := New()
	require.Nil(t, err)

	client.AddStatic("startTime", time.Now())

	time.Sleep(time.Second)

	start, ok := client.GetStatic("startTime")
	require.True(t, ok)
	startTime := start.(time.Time)

	elapsed := time.Since(startTime).Seconds()
	require.True(t, elapsed > 0)
}
