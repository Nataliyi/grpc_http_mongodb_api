package consts

import "time"

const (
	GRPC_PORT                      string        = "8080"
	GRPC_MAX_CONCURRENT_STREAMS    int           = 10
	GRPC_MAX_GOROUTINES_PER_STREAM int           = 30
	GRPC_CONN_DEADLINE_DURATION    time.Duration = 5 * time.Minute
)
