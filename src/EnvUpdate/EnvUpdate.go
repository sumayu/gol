package envupdate

import (
	"main/src/logger"
	mydb "main/src/mydb"
)

func EnvUpdateProd() error {
	db, err := mydb.Database()
	if err != nil {
		logger.Logger.Info("error connecting to db:", err)
		return err
	}

	_, err = db.Exec(`
		TRUNCATE TABLE envchecker;
		INSERT INTO envchecker (env) VALUES ('prod');
	`)
	if err != nil {
		logger.Logger.Info("error updating env to prod:", err)
		return err
	}

	logger.Logger.Info("env = prod")
	return nil
}

func EnvUpdateDebug() error {
	db, err := mydb.Database()
	if err != nil {
		logger.Logger.Info("error connecting to db:", err)
		return err
	}

	_, err = db.Exec(`
		TRUNCATE TABLE envchecker;
		INSERT INTO envchecker (env) VALUES ('debug');
	`)
	if err != nil {
		logger.Logger.Info("error updating env to debug:", err)
		return err
	}

	logger.Logger.Info("env = debug")
	return nil
}