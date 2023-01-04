package metrics

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricServer struct {
	Port   string
	server *http.Server
}

func NewMetricServer(port string, metricPath string) (*MetricServer, error) {
	if port == "" {
		return nil, fmt.Errorf("server port must not be empty")
	}
	if metricPath == "" {
		return nil, fmt.Errorf("server metricPath must not be empty")
	}
	// expose metric path when new server init
	http.Handle(metricPath, promhttp.Handler())
	return &MetricServer{
		Port: port,
	}, nil
}

// Blocking function! Should run like a goroutine.
// Starting the server for a certain period of time.
func (s *MetricServer) Start() {
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.Port),
		Handler: nil,
	}
	// blocking goroutine
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// Shutdown the server after "timer"..
func (serv *MetricServer) Stop(timer time.Duration) {
	time.Sleep(timer)
	serv.server.Close()
}
