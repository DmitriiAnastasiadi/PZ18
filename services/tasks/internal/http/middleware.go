package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"tech-ip-sem2/services/tasks/internal/client/authclient"
)

func AuthMiddleware(client *authclient.Client, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{"valid": false, "error": "unauthorized"})
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		subject, valid, err := client.Verify(ctx, token)
		if err != nil {
			log.Printf("calling grpc verify error: %v", err)
			if err == authclient.ErrUnauthenticated {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]interface{}{"valid": false, "error": "unauthorized"})
				return
			}
			if err == authclient.ErrUnavailable {
				w.WriteHeader(http.StatusServiceUnavailable)
				json.NewEncoder(w).Encode(map[string]interface{}{"error": "auth service unavailable"})
				return
			}

			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": "auth service error"})
			return
		}

		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{"valid": false, "error": "unauthorized"})
			return
		}

		ctx = context.WithValue(r.Context(), "subject", subject)
		next(w, r.WithContext(ctx))
	}
}
