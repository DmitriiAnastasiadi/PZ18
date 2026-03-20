package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tech-ip-sem2/proto/authpb"
	"tech-ip-sem2/services/auth/internal/grpcserver"
	httpHandlers "tech-ip-sem2/services/auth/internal/http"
	"tech-ip-sem2/shared/middleware"

	"google.golang.org/grpc"
)

func main() {
	httpPort := os.Getenv("AUTH_PORT")
	if httpPort == "" {
		httpPort = "8081"
	}

	grpcPort := os.Getenv("AUTH_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start gRPC server
	grpcLis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen on grpc port %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, &grpcserver.AuthGRPCServer{})

	go func() {
		log.Printf("Auth gRPC service starting on port %s", grpcPort)
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatalf("gRPC serve error: %v", err)
		}
	}()

	// HTTP endpoint (optional compatibility)
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/auth/login", httpHandlers.LoginHandler)
	mux.HandleFunc("/v1/auth/verify", grpcserver.VerifyHandler)

	httpServer := &http.Server{Addr: ":" + httpPort, Handler: middleware.LoggingMiddleware(middleware.RequestIDMiddleware(mux))}

	go func() {
		log.Printf("Auth HTTP service starting on port %s", httpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http serve error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("http server shutdown error: %v", err)
	}
	grpcServer.GracefulStop()
	log.Println("Auth service stopped")
}
