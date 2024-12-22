package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dbHost := viper.GetString("database.host")
	dbPort := os.Getenv("database.port")
	dbUser := os.Getenv("database.user")
	dbPassword := os.Getenv("database.password")
	dbName := os.Getenv("database.dbname")
	dbSsl := os.Getenv("database.sslmode")

	pgi := "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai"
	dsn := fmt.Sprintf(pgi, dbHost, dbUser, dbPassword, dbName, dbPort, dbSsl)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	DB = db
	log.Println("Database connection established")
}
