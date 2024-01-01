package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func ConnectToDB() (db *sql.DB) {
	// Load environment variables from the .env file
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}
	// database config
	cfg := mysql.Config{
		Net:                     os.Getenv("DB_PROTOCOL"),
		User:                    os.Getenv("MYSQL_USER"),
		Passwd:                  os.Getenv("MYSQL_PASSWORD"),
		Addr:                    os.Getenv("DB_ADDR"),
		DBName:                  os.Getenv("MYSQL_DATABASE"),
		ParseTime:               true, // Set this if you want to parse time values
		AllowNativePasswords:    true,
		CheckConnLiveness:       false,
		MaxAllowedPacket:        0,
		AllowOldPasswords:       false,
		InterpolateParams:       false,
		AllowCleartextPasswords: false,
	}

	dsn := cfg.FormatDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	log.Print("CONFIG: ", cfg)
	log.Print("DSN: ", dsn)

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	return db
}
