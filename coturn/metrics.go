package coturn

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Active Sessions Metrics
	activeSessions = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "coturn_active_sessions_total",
		Help: "The total number of active sessions",
	})

	activeSessionsAge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "coturn_active_sessions_age_seconds",
		Help: "Average age of active sessions",
	})

	activeSessionsExpiration = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "coturn_active_sessions_expiration_seconds",
		Help: "Average expiration of active sessions",
	})

	activeSessionsSentBytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "coturn_active_sessions_transmit_bytes",
		Help: "Average sent bytes of active sessions",
	})

	activeSessionsRecvBytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "coturn_active_sessions_receive_bytes",
		Help: "Average received bytes of active sessions",
	})

	activeSessionsTotalBytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "coturn_active_sessions_total_bytes",
		Help: "Average send and received bytes of active sessions",
	})

	// All Sessions Metrics
	sessionsAge = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "coturn_sessions_age_seconds",
		Help: "Summary age of all finished sessions",
	})

	sessionsSentBytes = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "coturn_sessions_transmit_bytes",
		Help: "Sent bytes of all finished sessions",
	})

	sessionsRecvBytes = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "coturn_sessions_receive_bytes",
		Help: "Received bytes of all finished sessions",
	})

	sessionsTotalBytes = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "coturn_sessions_total_bytes",
		Help: "Total send and received bytes of all finished sessions",
	})
)

func init() {
	// Active Sessions
	prometheus.MustRegister(activeSessions)
	prometheus.MustRegister(activeSessionsAge)
	prometheus.MustRegister(activeSessionsExpiration)
	prometheus.MustRegister(activeSessionsSentBytes)
	prometheus.MustRegister(activeSessionsRecvBytes)
	prometheus.MustRegister(activeSessionsTotalBytes)

	// All Sesssions
	prometheus.MustRegister(sessionsAge)
	prometheus.MustRegister(sessionsSentBytes)
	prometheus.MustRegister(sessionsRecvBytes)
	prometheus.MustRegister(sessionsTotalBytes)

}
