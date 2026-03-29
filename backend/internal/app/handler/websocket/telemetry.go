package websocket

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	meter = otel.Meter("websocket-server")
)

var (
	activeConns, _ = meter.Int64UpDownCounter("ws_active_connections",
		metric.WithDescription("Number of currently open WebSocket sessions"))
	msgSentCounter, _ = meter.Int64Counter("ws_messages_sent_total",
		metric.WithDescription("Total number of sent messages"))
	msgRecvCounter, _ = meter.Int64Counter("ws_messages_received_total",
		metric.WithDescription("Total number of received messages"))
	connDuration, _ = meter.Float64Histogram("ws_connection_duration_seconds",
		metric.WithDescription("Duration of WebSocket sessions"))
)
