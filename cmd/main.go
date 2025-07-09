package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/cfioretti/ingredients-balancer/internal/infrastructure/logging"
	"github.com/cfioretti/ingredients-balancer/pkg/application"
	grpcServer "github.com/cfioretti/ingredients-balancer/pkg/infrastructure/grpc"
	pb "github.com/cfioretti/ingredients-balancer/pkg/infrastructure/grpc/proto/generated"
)

const (
	defaultPort = ":50052"
	serviceName = "ingredients-balancer"
	version     = "1.0.0"
)

var logger *logging.Logger

func main() {
	logger = logging.NewLogger(serviceName, version)

	ctx := context.Background()
	logger.WithContext(ctx).Info("Starting ingredients-balancer service")

	port := getPort()
	logger.WithField("port", port).Info("Server configuration loaded")

	ingredientsBalancerService := application.NewIngredientsBalancerService()
	server := grpcServer.NewServer(ingredientsBalancerService)

	grpcInstance := grpc.NewServer(
		grpc.UnaryInterceptor(logger.GRPCUnaryInterceptor()),
		grpc.StreamInterceptor(logger.GRPCStreamInterceptor()),
	)

	pb.RegisterIngredientsBalancerServer(grpcInstance, server)
	logger.Info("gRPC service registered successfully")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.WithError(err).Fatal("Failed to listen on port")
	}

	// graceful shutdown
	go handleShutdown(grpcInstance)

	logger.WithField("port", port).Info("gRPC server starting")
	if err := grpcInstance.Serve(lis); err != nil {
		logger.WithError(err).Fatal("Failed to serve gRPC server")
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		logger.WithField("default_port", defaultPort).Info("Using default port")
		return defaultPort
	}

	fullPort := ":" + port
	logger.WithField("configured_port", fullPort).Info("Using configured port")
	return fullPort
}

func handleShutdown(server *grpc.Server) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	sig := <-sigCh
	logger.WithField("signal", sig.String()).Info("Received shutdown signal")

	logger.Info("Shutting down server gracefully...")
	server.GracefulStop()
	logger.Info("Server stopped successfully")
}
