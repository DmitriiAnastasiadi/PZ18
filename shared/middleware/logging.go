package middleware

import (
	"log"

	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestID := r.Context().Value("requestID")

		if requestID != nil {

			log.Printf("[%s] %s %s", requestID, r.Method, r.URL.Path)

		} else {

			log.Printf("%s %s", r.Method, r.URL.Path)

		}

		next.ServeHTTP(w, r)

	})

}
