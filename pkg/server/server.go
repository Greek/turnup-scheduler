package server

import (
	"context"
	"fmt"
	realLog "log"
	"log/slog"
	"net"
	"os"
	"turnup-scheduler/internal/logging"
	"turnup-scheduler/pkg/scheduler"

	pb "turnup-scheduler/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedSchedulerServiceServer
	Scheduler *scheduler.Scheduler
}

func authUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log := logging.BuildLogger("authUnaryInterceptor")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Warn("No metadata found in context")
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}
	tokens := md["authorization"]
	if len(tokens) == 0 {
		log.Warn("No authorization token provided")
		return nil, status.Error(codes.Unauthenticated, "missing authorization token")
	}
	expectedToken := os.Getenv("AUTH_TOKEN")
	if expectedToken == "" {
		log.Error("AUTH_TOKEN not set in environment")
		return nil, status.Error(codes.Internal, "server misconfiguration")
	}
	if tokens[0] != expectedToken {
		log.Warn("Invalid authorization token", slog.String("provided", tokens[0]))
		return nil, status.Error(codes.PermissionDenied, "invalid authorization token")
	}
	return handler(ctx, req)
}

func InitializeGrpcServer(port int, scheduler *scheduler.Scheduler) {
	log := logging.BuildLogger("InitializeGrpcServer")

	log.Info("Creating gRPC server")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		realLog.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authUnaryInterceptor),
	)
	pb.RegisterSchedulerServiceServer(grpcServer, &Server{
		Scheduler: scheduler,
	})
	grpcServer.Serve(lis)
}
