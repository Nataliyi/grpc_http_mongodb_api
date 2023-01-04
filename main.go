package main

import (
	"api/filter"
	"api/health"
	"api/services"
	"api/util"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	pb "api/proto/gen/go"

	"api/metrics"
	logger_models "api/models/logger"
	metric_models "api/models/metric"
	store_models "api/models/store"
	watcher_models "api/models/watcher"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

var (
	store         store_models.IStore
	watcher       watcher_models.IWatcher
	logger        logger_models.ILogger
	errorsCounter metric_models.IMetricCount
	healthCounter metric_models.IMetricCount
)

func main() {

	// load the config from the env
	cfg := &Config{}
	if err := cfg.ReadEnv(); err != nil {
		log.Fatalf("error reading env, %v", err)
	}
	// set the default config settings from consts
	cfg.SetDefaultConsts()

	// register Prometheus metrics
	healthCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "custom_api_heath_check",
		Help: "Tracking api health",
	})
	errorsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "custom_api_errors",
		Help: "The total number of api errors",
	})

	// create the metric server and start monitoring metrics
	metricsServer, err := metrics.NewMetricServer(
		cfg.MetricsSettings.Port, cfg.MetricsSettings.Path)
	if err != nil {
		log.Fatalf("failed to init metris server: %v", err)
	}
	go metricsServer.Start()
	defer metricsServer.Stop(cfg.MetricsSettings.ServerRuntime)

	// run and listen the watcher
	watcher, err = services.NewPubSubWatcher("mock", "example", "")
	if err != nil {
		log.Fatalf("failed to init watcher: %v", err)
	}
	go watcher.Listen()
	defer watcher.Close()

	// init the store client
	store, err = services.NewMongoStore(context.Background(),
		&store_models.StoreConfig{
			Login:    cfg.DBSettings.Login,
			Password: cfg.DBSettings.Password,
			Addr:     cfg.DBSettings.Addr,
			Port:     cfg.DBSettings.Port,
			DB:       cfg.DBSettings.DB,
			Table:    cfg.DBSettings.Table,
		},
	)
	if err != nil {
		log.Fatalf("failed to init store client: %v", err)
	}

	// init the logger
	logger, err = services.NewCustomLogger(
		cfg.LogsSettings.Prefix, cfg.LogsSettings.Path, cfg.LogsSettings.Frequency)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	// grpc serve
	serve(cfg)
}

func serve(cfg *Config) {

	addr := fmt.Sprintf("%s:%s", cfg.GRPCSettings.Host, cfg.GRPCSettings.Port)

	// create a listener on TCP port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCSettings.Port))
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// init the gRPC server
	grpcServer := grpc.NewServer(grpc.KeepaliveParams(
		keepalive.ServerParameters{
			MaxConnectionIdle:     cfg.GRPCSettings.ConnDeadlineDuration,
			MaxConnectionAge:      cfg.GRPCSettings.ConnDeadlineDuration,
			MaxConnectionAgeGrace: cfg.GRPCSettings.ConnDeadlineDuration,
			Time:                  cfg.GRPCSettings.ConnDeadlineDuration,
			Timeout:               cfg.GRPCSettings.ConnDeadlineDuration,
		}),
		grpc.ConnectionTimeout(cfg.GRPCSettings.ConnDeadlineDuration),
		grpc.MaxConcurrentStreams(uint32(cfg.GRPCSettings.MaxConcurrentStreams)),
	)
	// register the gRPC server
	pb.RegisterUsersStoreServer(grpcServer,
		&services.Server{
			MaxProcessingGoroutines: cfg.GRPCSettings.MaxGoriutinesPerStream,
			Store:                   store,
			Filter:                  &filter.BsonHelper{},
			ErrorsMetric:            errorsCounter,
			WatcherCh:               watcher.GetChannel(),
			Logger:                  logger,
		},
	)

	// keepalive probes
	grpc_health_v1.RegisterHealthServer(grpcServer,
		&health.HealthServer{
			HealthMetric: healthCounter,
		},
	)
	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	// serve the gRPC server
	log.Printf("Serving gRPC on %s", addr)
	go func() {
		log.Fatalln(grpcServer.Serve(lis))
	}()

	// create a client connection to the gRPC server
	conn, err := grpc.DialContext(
		context.Background(),
		addr,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	// register the gRPC server endpoint
	gwmux := runtime.NewServeMux()
	err = pb.RegisterUsersStoreHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.GRPCSettings.GatewayPort),
		Handler: cors(gwmux),
	}

	log.Printf("Serving gRPC gateway on %s:%s", cfg.GRPCSettings.Host,
		cfg.GRPCSettings.GatewayPort)

	log.Fatalln(gwServer.ListenAndServe())

}

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if util.AllowedOrigin(r.Header.Get("Origin")) {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")
		}
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
