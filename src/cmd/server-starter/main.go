// main.go
package main

import (
	"log"
	"main/src/backend/backendCFG"
	"main/src/envchecker"
	"main/src/logger" // Импортируйте пакет logger
	"main/src/mydb"
	"main/src/router"
	"net/http"
)

func main() {
//	pathToYaml := "../../../configYML/config.yaml"
	pathToYaml := "/app/configYML/config.yaml"
	a, err := configs.LoadConfigs(pathToYaml)
	

	// Инициализируйте логгер
	logger.InitLogger(envchecker.Envchecker())

	if err == nil {
		logger.Logger.Debug("configs -> starter(main.go) | TRUE |")
	} else {
		msg := "configs -> starter(main.go) | ERROR |"
		log.Fatalf(msg, a, err)
	}
	log.Println("server start on", envchecker.Envchecker(), "env")
	
mydb.Database()	



















	errorSERVER := http.ListenAndServe(":8080", router.Router())
if err != nil {
	logger.Logger.Info("error with create server API (router)", errorSERVER)
}
}

