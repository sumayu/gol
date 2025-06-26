package logger

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"time"
)
var pathToApp string
var filename string
func init() {

	if isDocker() {
		filename, _= getProjectRoot()
	} else {
	 filename  = "main/src/cmd/server-starter/app.log"
	}
}

var Logger *slog.Logger

func InitLogger(env string) {
	var handler slog.Handler

	switch env {
	case "debug", "local":
		logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			slog.Error("Failed to open log file:", err)
			return
		}

		multiWriter := io.MultiWriter(os.Stdout, logFile)
		handler = slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
			Level: slog.LevelDebug, // Уровень Debug
		})

	case "prod":
		logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			slog.Error("Failed to open log file:", err)
			return
		}

		multiWriter := io.MultiWriter(os.Stdout, logFile)
		handler = slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
			Level: slog.LevelInfo, // Уровень Info (меньше деталей, чем Debug)
		})

	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	Logger = slog.New(handler)
	go monitorLogFileSize()
}

// isDocker проверяет, работает ли приложение в Docker.
func isDocker() bool {
	_, err := os.Stat("/.dockerenv")
	return err == nil
}

// monitorLogFileSize следит за размером лог-файла и очищает его при превышении.
func monitorLogFileSize() {
	for {
		time.Sleep(10 * time.Minute)

		fileInfo, err := os.Stat(filename)
		if err != nil {
			slog.Error("Failed to check log file size:", err)
			continue
		}
		if fileInfo.Size() >= 1_000_000 {
			if err := os.Truncate(filename, 0); err != nil {
				slog.Error("Failed to truncate log file:", err)
			}
		}
	}
}
func PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if Logger != nil {
					Logger.Error("Panic recovered",
						"error", err,
						"stack", string(debug.Stack()),
					)
				}
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}


func getProjectRoot() (string, error) {
	_, filename, _, _ := runtime.Caller(0) 
	return filepath.Abs(filepath.Join(filepath.Dir(filename), "../../..")) 
}

func init() {
	root, err := getProjectRoot()
	if err != nil {
		panic("Failed to get project root: " + err.Error())
	}

	filename = filepath.Join(root, "gol", "src", "cmd", "server-starter", "app.log")
}