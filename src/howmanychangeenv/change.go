	package howmanychangeenv

	import (
		"main/src/mydb"
	)
	var saveHowManyChange int
	func Change() int {
	db, _ :=  mydb.Database()
		var count int
		err := db.QueryRow(`
			SELECT COALESCE(env_how_many_change, 0) 
			FROM envchecker
		`).Scan(&count)
		
		if err != nil || count == 0 {
			db.Exec("UPDATE envchecker SET env_how_many_change = 0")
			return 0
		}
		return count
	}