package router

import (
	"fmt"
	envupdate "main/src/EnvUpdate"
	howmanychangeenv "main/src/HowManyChangeEnv"
	"main/src/envchecker"
	"main/src/logger"
	websockethelper "main/src/websocketHelper"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Router() *chi.Mux {
	server := chi.NewRouter()
	
	basePath, _ := os.Getwd()
	logPath := filepath.Join(basePath, "src/cmd/server-starter/app.log")

	server.Get("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Debug("connect to api TRUE")
		http.Redirect(w, r, "/static/index.html", http.StatusMovedPermanently)
	})

	server.Get("/logger/log", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Logger.Error("WebSocket upgrade error:", err)
			return
		}
		defer conn.Close()

		// Канал для управления горутиной
		done := make(chan struct{})
		defer close(done)

		// Запускаем отправку логов
		go func() {
			websockethelper.WebsocketHelper(logPath, func(logLine string) {
				select {
				case <-done:
					return
				default:
					err := conn.WriteMessage(websocket.TextMessage, []byte(logLine))
					if err != nil {
						logger.Logger.Warn("WebSocket write error:", err)
						return
					}
				}
			}, done)
		}()

		// Обрабатываем входящие сообщения (для поддержания соединения)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					logger.Logger.Debug("WebSocket closed unexpectedly:", err)
				}
				break
			}
		}
	})

	// Остальные обработчики остаются без изменений
	server.Get("/env/current", func(w http.ResponseWriter, r *http.Request) {
		env := envchecker.Envchecker()
		if env == "" {
			http.Error(w, "Failed to get current environment", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(env))
	})

	server.Get("/env/changes", func(w http.ResponseWriter, r *http.Request) {
		changes := howmanychangeenv.Change()
		fmt.Fprintf(w, "%d", changes)
	})

	server.Post("/env/prod", func(w http.ResponseWriter, r *http.Request) {
		err := envupdate.EnvUpdateProd()
		if err != nil {
			http.Error(w, "Failed to update environment", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Env PROD")
	})

	server.Post("/env/debug", func(w http.ResponseWriter, r *http.Request) {
		err := envupdate.EnvUpdateDebug()
		if err != nil {
			http.Error(w, "Failed to update environment", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Env DEBUG")
	})

	server.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return server
}