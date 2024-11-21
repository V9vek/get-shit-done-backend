package config

import (
	"database/sql"
	"fmt"
	"get-shit-done/utils"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func DatabaseConnection() *sql.DB {
	err := godotenv.Load()
	utils.PanicIfError(err)

	dsn := os.Getenv("DATABASE_URL")
	driverName := "postgres"

	db, err := sql.Open(driverName, dsn)
	utils.PanicIfError(err)

	err = db.Ping()
	utils.PanicIfError(err)

	fmt.Println("db connected")
	return db
}
