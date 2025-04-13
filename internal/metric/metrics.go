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
type Metrics interface {
	IncMessageCounter()
}

// Metrics ...
type metrics struct {
	messageCounter prometheus.Counter
}

// NewMetrics ...
func NewMetrics() Metrics {
	return &metrics{
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
func (m *metrics) IncMessageCounter() {
	m.messageCounter.Inc()
}
