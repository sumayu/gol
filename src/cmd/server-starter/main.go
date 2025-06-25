package main

import (
	"main/src/backend/backendCFG"
	"main/src/envchecker"
	"main/src/logger"
	"main/src/mydb"
	"main/src/router"
	"net/http"
)

func main() {
	
	logger.InitLogger("prod")
	
	//pathToYaml := "../../../configYML/config.yaml"
	pathToYaml := "/app/configYML/config.yaml"
	cfg, err := configs.LoadConfigs(pathToYaml)
	if err != nil {
		logger.Logger.Error("Failed to load config", "error", err)
		return
	}
	logger.Logger.Debug("Config loaded successfully", "config", cfg)
	db, err := mydb.Database()
	if err != nil {
		logger.Logger.Error("Failed to initialize database", "error", err)
		return
	}
	defer db.Close()
	logger.Logger.Info("Server starting on", "env", envchecker.Envchecker())
	
	r := router.Router()
	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))
	logger.Logger.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Logger.Error("Server failed", "error", err)
	}

}