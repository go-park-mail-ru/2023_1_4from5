package middleware

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"time"
)

type ContextKey string

var ReqID = "ReqId"

type LoggerMiddleware struct {
	logger *zap.SugaredLogger
}

func NewLoggerMiddleware(logger *zap.SugaredLogger) LoggerMiddleware {
	return LoggerMiddleware{
		logger: logger,
	}
}

func (m *LoggerMiddleware) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		reqId := fmt.Sprintf("%016x", rand.Int())[:10]
		ctx = context.WithValue(ctx, ContextKey(ReqID), reqId)

		start := time.Now()
		next.ServeHTTP(w, r.WithContext(ctx))

		m.logger.Info(r.URL.Path,
			zap.String("reqId:", reqId),
			zap.String("method", r.Method),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("url", r.URL.Path),
			zap.Time("start", start),
			zap.Duration("work_time", time.Since(start)),
		)
	})
}
