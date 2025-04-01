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
    WITH updated AS (
        UPDATE envchecker 
        SET env = 'prod',
            env_how_many_change = COALESCE(env_how_many_change, 0) + 1
        RETURNING *
    )
    INSERT INTO envchecker (env, env_how_many_change)
    SELECT 'prod', 1
    WHERE NOT EXISTS (SELECT 1 FROM updated)
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
    WITH updated AS (
        UPDATE envchecker 
        SET env = 'debug',
            env_how_many_change = COALESCE(env_how_many_change, 0) + 1
        RETURNING *
    )
    INSERT INTO envchecker (env, env_how_many_change)
    SELECT 'debug', 1
    WHERE NOT EXISTS (SELECT 1 FROM updated)
`)
	if err != nil {
		logger.Logger.Info("error updating env to debug:", err)
		return err
	}

	logger.Logger.Info("env = debug")
	return nil
}