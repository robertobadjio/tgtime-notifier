package metric

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "tgtime"
	appName   = "notifier"
)

// Metrics ???
type Metrics struct {
	messageCounter prometheus.Counter
}

var metrics *Metrics

// Init ???
func Init(_ context.Context) error {
	metrics = &Metrics{
		messageCounter: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "telegram",
				Name:      appName + "_messages_total",
				Help:      "Количество отправленных сообщений сотрудникам",
			},
		),
	}

	return nil
}

// IncMessageCounter ...
func IncMessageCounter() {
	metrics.messageCounter.Inc()
}
