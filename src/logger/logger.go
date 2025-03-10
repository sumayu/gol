package logger

import (
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

var filname string = "app.log"
var Logger *slog.Logger

func InitLogger(env string) {
	switch env {
	case "debug":
		Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "local":
		Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		logFile, err := os.OpenFile(filname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			slog.Error("Failed to open log file:", err)
			return
		}
		Logger = slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	go monitorLogFileSize()
}

func monitorLogFileSize() {
	for {
		time.Sleep(60 * time.Second)

		fileSize, errSize := os.Stat(filname)
		if errSize != nil {
			slog.Error("Failed to open size file:", errSize)
			continue
		}

		// ВАЖНО INFOOOOOO ЕСЛИ НУЖНО ПОМЕНЯТЬ МАКСИМАЛЬНЫЙ РАЗМЕР ЛОГОВ ТО ТУТ МЕНЯЕШЬ 1000 БАЙТ НА 100000 ИЛИ СКОЛЬКО НУЖНО ИЛИ ПРОСТО МОЖНО 
		// удалить эту часть кода и тогда у логгов не будет лимита 
		if fileSize.Size() >= 1000 {
			err := os.Truncate(filname, 0)
			if err != nil {
				slog.Error("Failed to truncate log file:", err)
			}
		}
		// Тут логи из app.log удаляются если достигают размера 1000 байт (при кодировке utf 8 в среднем 1 символ ==  1 байт (8бит))
	}
}

func PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				Logger.Error("Panic recovered",
					"error", err,
					"stack", string(debug.Stack()),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}