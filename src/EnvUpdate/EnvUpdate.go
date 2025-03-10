package envupdate

import (
	"database/sql"
	"main/src/envchecker"
	"main/src/logger"
	mydb "main/src/mydb"

) 
func EnvUpdateProd() {
		
	db, _ := mydb.Database()
	_, err := db.Exec (`
TRUNCATE TABLE envchecker;
INSERT INTO envchecker (env) VALUES ('prod');

	`)
	
	if err != nil {
		logger.Logger.Info("error message to db",err)
	}
	logger.Logger.Info("env = " , envchecker.Envchecker())

	}
	func EnvUpdateDebug() {
		db ,_:= mydb.Database()
		_, err := db.Exec (`
TRUNCATE TABLE envchecker;
INSERT INTO envchecker (env) VALUES ('debug');

		`)
		if err != nil {
			logger.Logger.Info("error message to db",err)

		}
		logger.Logger.Info("env = " , envchecker.Envchecker())
	}
		func Rt () (*sql.DB,error) {
			db , error := mydb.Database()
			return db, error
		}