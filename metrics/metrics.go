package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/APX103/openalex-go"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	defaultNamespace = "openalex"
	defaultSubsystem = ""
)

type RecorderOption func(*PrometheusRecorder)

// WithRegisterer sets a custom prometheus.Registerer.
func WithRegisterer(r prometheus.Registerer) RecorderOption {
	return func(pr *PrometheusRecorder) { pr.registerer = r }
}

// WithNamespace sets the metrics namespace (default: "openalex").
func WithNamespace(ns string) RecorderOption {
	return func(pr *PrometheusRecorder) { pr.namespace = ns }
}

// WithSubsystem sets the metrics subsystem.
func WithSubsystem(sub string) RecorderOption {
	return func(pr *PrometheusRecorder) { pr.subsystem = sub }
}

// WithBuckets sets custom histogram buckets for request duration.
func WithBuckets(buckets []float64) RecorderOption {
	return func(pr *PrometheusRecorder) { pr.buckets = buckets }
}

// PrometheusRecorder implements openalex.RequestRecorder using Prometheus metrics.
type PrometheusRecorder struct {
	registerer prometheus.Registerer
	namespace  string
	subsystem  string
	buckets    []float64

	requestDuration         *prometheus.HistogramVec
	requestDurationSummary  *prometheus.SummaryVec
	requestTotal            *prometheus.CounterVec
}

// NewPrometheusRecorder creates and registers Prometheus metrics.
func NewPrometheusRecorder(opts ...RecorderOption) *PrometheusRecorder {
	pr := &PrometheusRecorder{
		registerer: prometheus.DefaultRegisterer,
		namespace:  defaultNamespace,
		subsystem:  defaultSubsystem,
		buckets:    []float64{0.1, 0.25, 0.5, 1, 2.5, 5, 10},
	}
	for _, opt := range opts {
		opt(pr)
	}

	pr.requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: pr.namespace,
		Subsystem: pr.subsystem,
		Name:      "request_duration_seconds",
		Help:      "Duration of OpenAlex API requests in seconds.",
		Buckets:   pr.buckets,
	}, []string{"endpoint"})

	pr.requestDurationSummary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  pr.namespace,
		Subsystem:  pr.subsystem,
		Name:       "request_duration_seconds_summary",
		Help:       "Summary of OpenAlex API request duration in seconds.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"endpoint"})

	pr.requestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: pr.namespace,
		Subsystem: pr.subsystem,
		Name:      "requests_total",
		Help:      "Total number of OpenAlex API requests.",
	}, []string{"endpoint", "status"})

	pr.registerer.MustRegister(pr.requestDuration, pr.requestDurationSummary, pr.requestTotal)
	return pr
}

// RecordRequest records metrics for a single API request.
func (pr *PrometheusRecorder) RecordRequest(_ context.Context, endpoint string, duration time.Duration, statusCode int) {
	pr.requestDuration.WithLabelValues(endpoint).Observe(duration.Seconds())
	pr.requestDurationSummary.WithLabelValues(endpoint).Observe(duration.Seconds())
	pr.requestTotal.WithLabelValues(endpoint, fmt.Sprintf("%d", statusCode)).Inc()
}

// Ensure PrometheusRecorder implements openalex.RequestRecorder at compile time.
var _ openalex.RequestRecorder = (*PrometheusRecorder)(nil)
