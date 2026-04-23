package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {

	// tenta carregar .env local (não quebra na Vercel)
	_ = godotenv.Load()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// default importante pra cloud
	if sslmode == "" {
		sslmode = "require"
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic("Erro ao abrir conexão com banco: " + err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic("Erro ao conectar no banco: " + err.Error())
	}

	fmt.Println("✅ Conectado ao PostgreSQL!")

	DB = db
}
