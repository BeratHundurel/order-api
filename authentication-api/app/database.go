package application

import (
	"fmt"
	"log"
	"os"

	"github.com/BeratHundurel/order-api/authentication-api/auth"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbInstance *gorm.DB

func InitializeDB() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	authToken := os.Getenv("TURSO_AUTH_TOKEN")
	databaseURL := os.Getenv("TURSO_DATABASE_URL")

	url := fmt.Sprintf("%s?authToken=%s", databaseURL, authToken)
	println(url)

	db, err := gorm.Open(sqlite.New(sqlite.Config{
		DriverName: "libsql",
		DSN:        url,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return err
	}

	dbInstance = db
	return nil
}

func MigrateDB() error {
	return dbInstance.AutoMigrate(&auth.User{})
}

func GetDB() *gorm.DB {
	return dbInstance
}
