package router

import (
	"fmt"
	envupdate "main/src/EnvUpdate"
	"main/src/envchecker"
	"main/src/logger"
	"net/http"

	"github.com/go-chi/chi"
)

func Router() *chi.Mux {
	server := chi.NewRouter()

	server.Get("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Debug("connect to api TRUE")
		http.Redirect(w, r, "/static/index.html", http.StatusMovedPermanently)
	})

	server.Get("/env/current", func(w http.ResponseWriter, r *http.Request) {
		env := envchecker.Envchecker()
		if env == "" {
			http.Error(w, "Failed to get current environment", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(env))
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