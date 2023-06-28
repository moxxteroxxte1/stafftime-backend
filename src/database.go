package main
/*
import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/moxxteroxxte1/stafftime-backend/src/models"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func Connect() *gorm.DB {
	dsn := fmt.Sprintf("host=stafftime-backend-database user=%s password=%s dbname=%s port=5432 sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	fmt.Println("connected!")

	fmt.Println("migrating models")
	db.AutoMigrate(&models.Contract{},&models.Payment{},&models.Role{},&models.Shift{},&models.Status{},&models.User{},&models.Tokens{})

	return db
}*/