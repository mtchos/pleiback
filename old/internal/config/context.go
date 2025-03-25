package config

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

var requestCounter uint64

func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := fmt.Sprintf("&d-&d", time.Now().UnixNano(), atomic.AddUint64(&requestCounter, 1))
		ctx := context.WithValue(r.Context(), "request_id", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
