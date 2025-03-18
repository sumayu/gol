package router

import (
	"fmt"
	envupdate "main/src/EnvUpdate"
	"main/src/logger"
	"net/http"

	"github.com/go-chi/chi"
)

func Router() *chi.Mux {
	server := chi.NewRouter()

	// Маршрут для главной страницы
	server.Get("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Debug("connect to api TRUE ")
		http.Redirect(w,r, "/static/index.html", http.StatusMovedPermanently)


	})

	// Маршрут для переключения на PROD
	server.Get("/env/prod", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Env PROD")
		envupdate.EnvUpdateProd()
	})

	// Маршрут для переключения на DEBUG
	server.Get("/env/debug", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Env DEBUG")

		envupdate.EnvUpdateDebug()
	})

	// Маршрут для статических файлов (HTML, CSS, JS)
	server.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return server
}