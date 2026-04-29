package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type metricKey struct {
	route  string
	method string
	status int
}

type metricsRegistry struct {
	mu       sync.Mutex
	started  time.Time
	requests map[metricKey]uint64
}

func newMetricsRegistry() *metricsRegistry {
	return &metricsRegistry{
		started:  time.Now(),
		requests: make(map[metricKey]uint64),
	}
}

func (m *metricsRegistry) record(route string, method string, status int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := metricKey{
		route:  route,
		method: method,
		status: status,
	}

	m.requests[key]++
}

func (m *metricsRegistry) snapshot() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	keys := make([]metricKey, 0, len(m.requests))
	for key := range m.requests {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i int, j int) bool {
		if keys[i].route != keys[j].route {
			return keys[i].route < keys[j].route
		}
		if keys[i].method != keys[j].method {
			return keys[i].method < keys[j].method
		}
		return keys[i].status < keys[j].status
	})

	lines := make([]string, 0, len(keys))
	for _, key := range keys {
		lines = append(lines, fmt.Sprintf(
			`gitops_api_requests_total{route="%s",method="%s",status="%d"} %d`,
			key.route,
			key.method,
			key.status,
			m.requests[key],
		))
	}

	return lines
}

func (m *metricsRegistry) handler(cfg config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		w.WriteHeader(http.StatusOK)

		var builder strings.Builder
		builder.WriteString("# HELP gitops_api_build_info Build information for the API service.\n")
		builder.WriteString("# TYPE gitops_api_build_info gauge\n")
		builder.WriteString(fmt.Sprintf(
			`gitops_api_build_info{service="%s",version="%s",environment="%s"} 1`+"\n",
			cfg.serviceName,
			cfg.version,
			cfg.environment,
		))
		builder.WriteString("# HELP gitops_api_uptime_seconds Seconds since the API process started.\n")
		builder.WriteString("# TYPE gitops_api_uptime_seconds gauge\n")
		builder.WriteString("gitops_api_uptime_seconds " + strconv.FormatFloat(time.Since(m.started).Seconds(), 'f', 0, 64) + "\n")
		builder.WriteString("# HELP gitops_api_requests_total Total number of HTTP requests by route, method, and status code.\n")
		builder.WriteString("# TYPE gitops_api_requests_total counter\n")

		for _, line := range m.snapshot() {
			builder.WriteString(line)
			builder.WriteString("\n")
		}

		_, _ = w.Write([]byte(builder.String()))
	})
}

func instrument(route string, metrics *metricsRegistry, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(recorder, r)
		metrics.record(route, r.Method, recorder.status)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(statusCode int) {
	s.status = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}
