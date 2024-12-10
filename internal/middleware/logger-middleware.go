package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

var Logger = zap.NewExample()

func LoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload interface{}
		Logger.Info(r.Method, zap.String("url", r.URL.Path), zap.String("time", time.Now().Format("Mon Jan 2 15:04:05 MST 2006")), zap.Any("Payload", payload))
		h.ServeHTTP(w, r)
	})
}
