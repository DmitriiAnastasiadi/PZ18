package middleware

import (
	"context"

	"math/rand"

	"net/http"

	"strconv"

	"time"
)

func init() {

	rand.Seed(time.Now().UnixNano())

}

func RequestIDMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestID := r.Header.Get("X-Request-ID")

		if requestID == "" {

			requestID = "req-" + strconv.Itoa(rand.Intn(1000000))

		}

		ctx := context.WithValue(r.Context(), "requestID", requestID)

		r = r.WithContext(ctx)

		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r)

	})

}
