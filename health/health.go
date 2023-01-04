package health

import (
	"context"

	metric_model "api/models/metric"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// HealthServer for the Health Check gRPC API
type HealthServer struct {
	HealthMetric metric_model.IMetricCount
}

// Check is used for health checks
func (s *HealthServer) Check(ctx context.Context,
	in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {

	// send to the health metric
	s.HealthMetric.Add(1)

	// This is where you can implement checks of your service status
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

// Watch is not implemented
func (s *HealthServer) Watch(in *grpc_health_v1.HealthCheckRequest,
	srv grpc_health_v1.Health_WatchServer) error {

	return status.Error(codes.Unimplemented, "Watch is not implemented")
}
