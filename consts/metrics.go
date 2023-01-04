package consts

import "time"

const (
	METRICS_PORT           string        = "9090"
	METRICS_SERVER_TIMEOUT time.Duration = 1 * time.Minute
	METRICS_PATH           string        = "/metrics"
)
