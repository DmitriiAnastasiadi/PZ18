package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"tech-ip-sem2/services/tasks/internal/client/authclient"
	httpPkg "tech-ip-sem2/services/tasks/internal/http"
	"tech-ip-sem2/shared/middleware"
)

func main() {
	port := os.Getenv("TASKS_PORT")
	if port == "" {
		port = "8082"
	}

	authAddr := os.Getenv("AUTH_GRPC_ADDR")
	if authAddr == "" {
		authAddr = "localhost:50051"
	}

	authClient, err := authclient.New(authAddr)
	if err != nil {
		log.Fatalf("failed to connect to auth gRPC (%s): %v", authAddr, err)
	}
	defer authClient.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/tasks", httpPkg.AuthMiddleware(authClient, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			httpPkg.CreateTaskHandler(w, r)
		} else if r.Method == "GET" {
			httpPkg.ListTasksHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", 405)
		}
	}))

	mux.HandleFunc("/v1/tasks/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
		if id == "" {
			http.NotFound(w, r)
			return
		}

		httpPkg.AuthMiddleware(authClient, func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "GET":
				httpPkg.GetTaskHandler(w, r, id)
			case "PATCH":
				httpPkg.UpdateTaskHandler(w, r, id)
			case "DELETE":
				httpPkg.DeleteTaskHandler(w, r, id)
			default:
				http.Error(w, "Method not allowed", 405)
			}
		})(w, r)
	})

	handler := middleware.LoggingMiddleware(middleware.RequestIDMiddleware(mux))

	server := &http.Server{Addr: ":" + port, Handler: handler}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("Tasks service starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http serve error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	log.Println("Tasks service stopped")
}
