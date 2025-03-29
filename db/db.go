package db

import (
	"fmt"
	"os"

	"github.com/cmerin0/SimpleCarsApp/models"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstace struct {
	Db *gorm.DB // Pointer to a GORM database connection
}

var DB Dbinstace // variable of instace of database

func ConnectDB() {
	// Data Source Name
	dsn := fmt.Sprintf("host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	// Connecting to database through PostgreSQL driver of GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	// Checking if we receive an error
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	// If not erros, we are connected to database
	log.Info("Connected to database")
	db.Logger = logger.Default.LogMode(logger.Info)

	// Running migrations
	log.Info("Running migrations")
	db.AutoMigrate(&models.Make{}, &models.Car{}, &models.User{})

	// Set the value of global db variable
	DB = Dbinstace{Db: db}

}
