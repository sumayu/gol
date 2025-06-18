package router

import (
	"context"
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
	server.Use(RecoveryMiddleware)
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
	
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
	
		go func() {
			err := websockethelper.WebsocketHelper(ctx, logPath, func(logLine string) {
				select {
				case <-ctx.Done():  // Проверяем отмену контекста
					return
				default:
					if err := conn.WriteMessage(websocket.TextMessage, []byte(logLine)); err != nil {
						logger.Logger.Warn("WebSocket write error:", err)
						cancel()  // Отменяем контекст при ошибке
						return
					}
				}
			})
			
			if err != nil {
				logger.Logger.Error("WebsocketHelper failed:", err)
			}
		}()
	
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					logger.Logger.Debug("WebSocket closed:", err)
				}
				return
			}
		}
	})
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