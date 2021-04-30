package metrics

import (
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStartPushMetrics(t *testing.T) {
	testMetricName := "testMetric"

	testMetric := prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "Tests",
		Name:      testMetricName,
		Help:      "Should be in r.Body",
	}, nil)
	testMetric.Add(1)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		res := string(resBytes)
		t.Log(res)

		if !strings.Contains(res, testMetricName) {
			t.Fatal("r.Body doesn't contains out test metric!")
		}

		w.WriteHeader(202)

	}))
	defer ts.Close()

	pusher := push.New(ts.URL, "my_job").Gatherer(stdprometheus.DefaultGatherer)
	err := pusher.Push()
	if err != nil {
		t.Fatal("can't push to server", err)
	}

}
