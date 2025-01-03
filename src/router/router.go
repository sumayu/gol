package router

import (
	"fmt"
	envupdate "main/src/EnvUpdate"
	"net/http"

	"github.com/go-chi/chi"
)

func Router() *chi.Mux {

	server := chi.NewRouter()
server.Get("/",func(w http.ResponseWriter, r *http.Request){
			fmt.Fprintf(w, "connect to api TRUE")
		})
		server.Get("/env/prod",func(w http.ResponseWriter, r *http.Request){
			fmt.Fprintf(w, "Env PROD")
			envupdate.EnvUpdateProd()
			
		})
		server.Get("/env/debug",func(w http.ResponseWriter, r *http.Request){
			fmt.Fprintf(w, "Env DEBUG")
			envupdate.EnvUpdateDebug()	
		})
return server
	}	


	