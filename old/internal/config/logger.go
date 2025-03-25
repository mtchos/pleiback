package config

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		slog.Info("incoming request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"request_id", middleware.GetReqID(r.Context()))

		next.ServeHTTP(w, r)

		slog.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start).String(),
			"request_id", middleware.GetReqID(r.Context()))
	})
}

func SetupLogger() {
	opts := &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelInfo,
		ReplaceAttr: removeTime,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func removeTime(_ []string, attr slog.Attr) slog.Attr {
	if attr.Key == "time" {
		if t, ok := attr.Value.Any().(time.Time); ok {
			attr.Value = slog.StringValue(t.Format(time.RFC3339))
		}
	}

	return attr
}
