package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// WSConnections is a metric to show the number of active web socket connections being handled
	WSConnections = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "commentron",
		Subsystem: "websocket",
		Name:      "connections",
		Help:      "Number of active web socket connections",
	}, []string{"claim"})

	// UserLoadOverall Number of active users
	UserLoadOverall = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "commentron",
		Subsystem: "apis",
		Name:      "user_load",
		Help:      "Number of active users",
	})

	// UserLoadByAPI Number of active calls by api
	UserLoadByAPI = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "commentron",
		Subsystem: "apis",
		Name:      "api_load",
		Help:      "Number of active calls by api",
	}, []string{"version", "service", "method"})

	// Durations The durations of the individual api calls
	Durations = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "commentron",
		Subsystem: "apis",
		Name:      "duration",
		Help:      "The durations of the individual api calls",
	}, []string{"version", "service", "method"})

	// SDKDurations The durations of the individual api calls
	SDKDurations = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "commentron",
		Subsystem: "sdk",
		Name:      "duration",
		Help:      "The durations of the individual sdk api calls",
	}, []string{"method"})

	// SDKClaimCache is a metric to show the miss hit ration of the claim cache
	SDKClaimCache = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "commentron",
		Subsystem: "cache",
		Name:      "sdk_claim",
		Help:      "SDK claim cache miss/hit",
	}, []string{"type"})
	// JobsDuration The type of job and the duration it runs for
	JobsDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "commentron",
		Subsystem: "jobs",
		Name:      "duration",
		Help:      "Runs of each job measuring duration",
	}, []string{"job"})
	// CommentsClassified reports the number of successful comment classifications
	CommentsClassified = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "commentron",
		Subsystem: "moderation",
		Name:      "comments_classified",
		Help:      "Number of successful comments classified",
	})
	// PollingCallsForClassifierJob is a metric to show the number of calls to the classifier job
	PollingCallsForClassifierJob = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "commentron",
		Subsystem: "moderation",
		Name:      "polling_calls_for_classifier_job",
		Help:      "Number of calls to the classifier job",
	})
)

// SDKCall helper function for observing the duration
func SDKCall(start time.Time, callType string) {
	duration := time.Since(start).Seconds()
	SDKDurations.WithLabelValues(callType).Observe(duration)
}

// Job helper function for observing the duration
func Job(start time.Time, name string) {
	duration := time.Since(start).Seconds()
	JobsDuration.WithLabelValues(name).Observe(duration)
}
