package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "tgtime"
	appName   = "notifier"
)

// Metrics ...
type Metrics struct {
	messageCounter prometheus.Counter
}

// NewMetrics ...
func NewMetrics() *Metrics {
	return &Metrics{
		messageCounter: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "telegram",
				Name:      appName + "_messages_total",
				Help:      "Количество отправленных сообщений сотрудникам",
			},
		),
	}
}

// IncMessageCounter ...
func (m *Metrics) IncMessageCounter() {
	m.messageCounter.Inc()
}
