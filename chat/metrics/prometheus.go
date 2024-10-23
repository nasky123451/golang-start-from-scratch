package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	MessageSendCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "chat_message_sent_total",
		Help: "Total number of chat messages sent",
	})

	MessageReceiveCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "chat_message_received_total",
		Help: "Total number of chat messages received",
	})
)

func InitMetrics() {
	prometheus.MustRegister(MessageSendCounter)
	prometheus.MustRegister(MessageReceiveCounter)
}
