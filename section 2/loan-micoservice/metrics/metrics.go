package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	ApprovalRate = prometheus.NewCounter(
		prometheus.CounterOpts{Name: "underwriting_approved_total", Help: "Number of approved applications"})
	RejectionRate = prometheus.NewCounter(
		prometheus.CounterOpts{Name: "underwriting_rejected_total", Help: "Number of rejected applications"})
	CreditLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{Name: "credit_score_request_duration_seconds", Help: "Credit check duration"})
)

func init() {
	prometheus.MustRegister(ApprovalRate, RejectionRate, CreditLatency)
}
