package config

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
)

var (
	//DB ...
	DB *gorm.DB
)

// InitDB connection
func InitDB() {
	var err error

	//connect to postgres database
	DB, err = gorm.Open("postgres", os.Getenv("DB_URI"))
	if err != nil {
		log.Fatal("Failed to connect to database:\n" + err.Error())
	}
	if err = DB.DB().Ping(); err != nil {
		log.Fatal("Failed to ping database:\n" + err.Error())
	}
	DB.SingularTable(true)

	log.Println("Successfully connected to database!")
}
