package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestRecordRequest(t *testing.T) {
	reg := prometheus.NewRegistry()
	r := NewPrometheusRecorder(WithRegisterer(reg))

	r.RecordRequest(context.Background(), "works", 150*time.Millisecond, 200)
	r.RecordRequest(context.Background(), "works", 50*time.Millisecond, 404)
	r.RecordRequest(context.Background(), "authors", 300*time.Millisecond, 200)

	mfs, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather metrics: %v", err)
	}

	duration := findMetric(mfs, "openalex_request_duration_seconds")
	if duration == nil {
		t.Fatal("duration metric not found")
	}
	if len(duration.GetMetric()) != 2 {
		t.Fatalf("expected 2 duration buckets (works, authors), got %d", len(duration.GetMetric()))
	}

	total := findMetric(mfs, "openalex_requests_total")
	if total == nil {
		t.Fatal("total metric not found")
	}
	if len(total.GetMetric()) != 3 {
		t.Fatalf("expected 3 total counters (works/200, works/404, authors/200), got %d", len(total.GetMetric()))
	}
}

func TestRecordRequestNetworkError(t *testing.T) {
	reg := prometheus.NewRegistry()
	r := NewPrometheusRecorder(WithRegisterer(reg))

	r.RecordRequest(context.Background(), "works", 5*time.Second, 0)

	mfs, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather metrics: %v", err)
	}

	total := findMetric(mfs, "openalex_requests_total")
	if total == nil {
		t.Fatal("total metric not found")
	}

	found := false
	for _, m := range total.GetMetric() {
		for _, l := range m.GetLabel() {
			if l.GetName() == "status" && l.GetValue() == "0" {
				found = true
			}
		}
	}
	if !found {
		t.Fatal("expected status=0 label for network error")
	}
}

func TestCustomBuckets(t *testing.T) {
	reg := prometheus.NewRegistry()
	custom := []float64{0.01, 0.05, 0.1}
	r := NewPrometheusRecorder(WithRegisterer(reg), WithBuckets(custom))

	r.RecordRequest(context.Background(), "works", 75*time.Millisecond, 200)

	mfs, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather metrics: %v", err)
	}

	duration := findMetric(mfs, "openalex_request_duration_seconds")
	if duration == nil {
		t.Fatal("duration metric not found")
	}
	if len(duration.GetMetric()[0].GetHistogram().GetBucket()) != 3 {
		t.Fatalf("expected 3 custom buckets, got %d", len(duration.GetMetric()[0].GetHistogram().GetBucket()))
	}
}

func findMetric(mfs []*dto.MetricFamily, name string) *dto.MetricFamily {
	for _, mf := range mfs {
		if mf.GetName() == name {
			return mf
		}
	}
	return nil
}

func TestTextExporter(t *testing.T) {
	reg := prometheus.NewRegistry()
	r := NewPrometheusRecorder(WithRegisterer(reg))

	r.RecordRequest(context.Background(), "works", 150*time.Millisecond, 200)
}
