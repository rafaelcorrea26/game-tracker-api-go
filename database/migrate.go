package database

import (
	"fmt"
	"os"
)

func RunMigrations() {

	files := []string{
		"migrations/001_create_games.sql",
		"migrations/002_seed_games.sql",
		"migrations/003_create_users.sql",
	}

	for _, file := range files {

		sqlBytes, err := os.ReadFile(file)

		if err != nil {
			panic(err)
		}

		_, err = DB.Exec(string(sqlBytes))

		if err != nil {
			panic(err)
		}

		fmt.Println("Migration executada:", file)
	}

}
