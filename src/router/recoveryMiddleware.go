package router

import (
	"main/src/logger"
	"net/http"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
logger.Logger.Error("PANIC:",err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Internal Server Error"))
			}
		}()

		next.ServeHTTP(w, r)
	})
}