package auth

import (
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbInstance *gorm.DB

func GetDB() *gorm.DB {
	return dbInstance
}

func init() {
	auth_token := os.Getenv("TURSO_AUTH_TOKEN")
	database_url := os.Getenv("TURSO_DATABASE_URL")

	url := fmt.Sprintf("%s?authToken=%s", database_url, auth_token)

	db, err := gorm.Open(sqlite.New(sqlite.Config{
		DriverName: "libsql",
		DSN:        url,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dbInstance = db
}
