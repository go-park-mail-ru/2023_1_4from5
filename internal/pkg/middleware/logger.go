package middleware

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
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

		next.ServeHTTP(w, r.WithContext(ctx))

		m.logger.Info(r.URL.Path,
			zap.String("method", r.Method),
			zap.String("remote_addr", r.RemoteAddr),
		)
	})
}
