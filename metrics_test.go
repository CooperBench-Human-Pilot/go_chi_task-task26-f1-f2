package chi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSimpleMetricsCollector(t *testing.T) {
	var recorded []RouteMetric
	collector := &SimpleMetricsCollector{
		Callback: func(metric RouteMetric) {
			recorded = append(recorded, metric)
		},
	}

	r := NewRouter()
	r.SetMetricsCollector(collector)
	r.Get("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/hello/world", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if len(recorded) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(recorded))
	}

	m := recorded[0]
	if m.Method != http.MethodGet {
		t.Errorf("expected method GET, got %s", m.Method)
	}
	if m.Path != "/hello/world" {
		t.Errorf("expected path /hello/world, got %s", m.Path)
	}
	if m.Pattern != "/hello/{name}" {
		t.Errorf("expected pattern /hello/{name}, got %s", m.Pattern)
	}
	if m.Duration <= 0 {
		t.Errorf("expected positive duration, got %v", m.Duration)
	}
}

func TestMetricsCollectorNotCalledOnNotFound(t *testing.T) {
	var called bool
	collector := &SimpleMetricsCollector{
		Callback: func(metric RouteMetric) {
			called = true
		},
	}

	r := NewRouter()
	r.SetMetricsCollector(collector)
	r.Get("/exists", func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest(http.MethodGet, "/notfound", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if called {
		t.Error("metrics collector should not be called for unmatched routes")
	}
}

func TestMetricsCollectorDuration(t *testing.T) {
	delay := 10 * time.Millisecond
	var recorded RouteMetric
	collector := &SimpleMetricsCollector{
		Callback: func(metric RouteMetric) {
			recorded = metric
		},
	}

	r := NewRouter()
	r.SetMetricsCollector(collector)
	r.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
	})

	req := httptest.NewRequest(http.MethodGet, "/slow", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if recorded.Duration < delay {
		t.Errorf("expected duration >= %v, got %v", delay, recorded.Duration)
	}
}

func TestSimpleMetricsCollectorNilCallback(t *testing.T) {
	collector := &SimpleMetricsCollector{}

	// Should not panic with nil callback.
	collector.RecordHit(context.Background(), nil, RouteMetric{})
}

func TestNoMetricsCollector(t *testing.T) {
	r := NewRouter()
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	// Should not panic when no collector is set.
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
