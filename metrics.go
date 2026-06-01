package chi

import (
	"context"
	"net/http"
	"time"
)

// MetricsCollector is the interface for collecting route hit metrics.
type MetricsCollector interface {
	RecordHit(ctx context.Context, r *http.Request, metric RouteMetric)
}

// RouteMetric contains data about a single route hit.
type RouteMetric struct {
	Pattern   string
	Method    string
	Path      string
	Duration  time.Duration
	URLParams RouteParams
}

// SimpleMetricsCollector is a basic MetricsCollector that invokes a callback function.
type SimpleMetricsCollector struct {
	Callback func(metric RouteMetric)
}

// RecordHit calls the callback with the given RouteMetric.
func (s *SimpleMetricsCollector) RecordHit(ctx context.Context, r *http.Request, metric RouteMetric) {
	if s.Callback != nil {
		s.Callback(metric)
	}
}
