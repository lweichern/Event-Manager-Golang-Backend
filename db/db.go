package db

import (
	"database/sql"
	"fmt"
	"go-chi/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

var Db *sql.DB

func ConnectDb() {
	err := godotenv.Load() // access .env file
	if err != nil {
		panic("Error occured on .env file...")
	}

	// set up postgresql to open
	psqlSetup := os.Getenv("DB_RAILWAY_URL")

	// connect to postgres db
	db, errSql := sql.Open("postgres", psqlSetup)

	if errSql != nil {
		fmt.Println("There is an error while connecting to the database ", errSql)
		panic(errSql)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	CreateDbTables(db)

	Db = db
}

func ConfigDb(db *gorm.DB) {
	// Migrate the schema
	db.AutoMigrate(&models.User{}, &models.Event{}, &models.Task{})
}

func CreateDbTables(db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS events(
			id SERIAL PRIMARY KEY,
			title VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`

	_, err := db.Exec(query)

	if err != nil {
		log.Fatal(err)
	}

	query = `
		CREATE TABLE IF NOT EXISTS tasks(
			id SERIAL PRIMARY KEY,
			title VARCHAR(100) NOT NULL,
			is_done BOOLEAN DEFAULT false,
			event_id INTEGER NOT NULL,
			task_type VARCHAR(100) DEFAULT 'backlog',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err = db.Exec(query)

	if err != nil {
		log.Fatal(err)
	}
}
