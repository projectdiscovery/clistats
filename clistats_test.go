package clistats

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func TestStartMultipleTimes(t *testing.T) {
	client, err := New()
	require.Nil(t, err)

	for i := 1; i <= 2; i++ {
		err = client.Start()
		require.Nil(t, err)

		err = client.Stop()
		require.Nil(t, err)
	}
}

func TestStartMultipleTimesWithoutStopping(t *testing.T) {
	client, err := New()
	require.Nil(t, err)

	for i := 1; i <= 2; i++ {
		err = client.Start()

		if i == 1 {
			require.Nil(t, err)
		} else {
			require.NotNil(t, err)
		}
	}

	err = client.Stop()
	require.Nil(t, err)
}

func TestMetricsHandlerPercentWithZeroTotal(t *testing.T) {
	client, err := New()
	require.Nil(t, err)

	client.AddCounter("requests", 100)
	client.AddCounter("total", 0)
	client.AddStatic("startedAt", time.Now())

	recorder := httptest.NewRecorder()
	client.metricsHandler(recorder, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	require.Equal(t, http.StatusOK, recorder.Code)

	items := make(map[string]interface{})
	require.Nil(t, json.Unmarshal(recorder.Body.Bytes(), &items))
	require.Equal(t, "0", items["percent"])
}
