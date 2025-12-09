package config

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
	// Godotenv has been called at this point in time
	/*
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
			return
		}
	*/

	username, exists := os.LookupEnv("DB_USERNAME")

	if !exists {
		log.Fatal("DB_USERNAME environment variable is not set")
		return
	}

	password, exists := os.LookupEnv("DB_PASSWORD")

	if !exists {
		log.Fatal("DB_PASSWORD environment variable is not set")
		return
	}

	database, exists := os.LookupEnv("DB_DATABASE")

	if !exists {
		log.Fatal("DB_DATABASE environment variable is not set")
		return
	}

	host, exists := os.LookupEnv("DB_HOST")

	if !exists {
		log.Fatal("DB_HOST environment variable is not set")
		return
	}

	port, exists := os.LookupEnv("DB_PORT")

	if !exists {
		log.Fatal("DB_PORT environment variable is not set")
		return
	}

	var err error
	DB, err = sql.Open(
		"mysql",
		username+":"+password+
			"@tcp("+host+":"+port+")/"+database+
			"?parseTime=true&loc=Asia%2FJakarta",
	)

	if err != nil {
		log.Fatal(err)
	}

	DB.SetMaxOpenConns(20)
	DB.SetMaxIdleConns(20)
	DB.SetConnMaxLifetime(1 * time.Hour)
	DB.SetConnMaxIdleTime(5 * time.Minute)

	if err := DB.Ping(); err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	log.Println("Connected to database successfully")
}
