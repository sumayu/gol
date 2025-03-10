package configs

import (
	"log"


	"github.com/ilyakaznacheev/cleanenv"
)

type Сonfig struct {
	Env string `yaml:"env" env-default:"prod"`
	DB
	Server
}

type DB struct {
	BdPath string `yaml:"bdpath"`
}
type Server struct {
	Adr string `yaml:"adr"`
}

	func LoadConfigs(filePath string) (*Сonfig, error ) { 
		
	var cfg Сonfig	

err := cleanenv.ReadConfig(filePath, &cfg)
if err != nil {
	msg := ("   ||  configs -> starter(main.go) | ERROR |")

	log.Fatal(err,msg, "|backend/config.go| ?CHECK PATH? ")
	return nil,err 
} else {
	return &cfg , nil
}
}
