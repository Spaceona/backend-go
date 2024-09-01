package logging

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestsReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "spacesona_requests_received_total",
		Help: "total number of http requests received",
	})
	RequestsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "spacesona_requests_processed_total",
		Help: "total number of http requests processed",
	})
	RequestsAuthenticated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "spacesona_requests_authenticated_total",
		Help: "total number of http requests authenticated",
	})
	RequestsAuthenticatedFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "spacesona_requests_authenticated_failed_total",
		Help: "total number of http requests failed authenticated",
	})
	UpdateStatusDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "spacesona_update_status_duration_seconds",
		Help: "Duration of updating status of board",
	})
	GetStatusDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "spacesona_get_status_duration_seconds",
		Help: "Duration of getting status of board",
	})
)
