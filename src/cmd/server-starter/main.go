package main

import (
	"main/src/backend/backendCFG"
	"main/src/envchecker"
	"main/src/logger" // Импортируйте пакет logger
	"main/src/mydb"
	"main/src/router"
	"net/http"
)

func main() {
	// Инициализация логгера
	logger.InitLogger("prod")

	// Загрузка конфигурации
	pathToYaml := "../../../configYML/config.yaml"
	// pathToYaml := "/app/configYML/config.yaml"
	cfg, err := configs.LoadConfigs(pathToYaml)
	if err != nil {
		logger.Logger.Error("Failed to load config", "error", err)
		return
	}
	logger.Logger.Debug("Config loaded successfully", "config", cfg)

	// Инициализация базы данных
	db, err := mydb.Database()
	if err != nil {
		logger.Logger.Error("Failed to initialize database", "error", err)
		return
	}
	defer db.Close()

	// Логирование окружения
	logger.Logger.Info("Server starting on", "env", envchecker.Envchecker())

	// Создаем роутер
	r := router.Router()

	// Отдаем статические файлы (HTML, CSS, JS)
	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Запуск сервера
	logger.Logger.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Logger.Error("Server failed", "error", err)
	}
}