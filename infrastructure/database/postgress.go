package database

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/guiPython/codepix/domain/model"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "gorm.io/driver/sqlite"
)

func init() {
	_, path, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(path)

	if err := godotenv.Load(basepath + "/../../.env"); err != nil {
		log.Fatalf("error loading .env file")
	}
}

func ConnectDB(env string) *gorm.DB {
	var connection_string string
	var database *gorm.DB
	var err error

	if env != "test" {
		connection_string = os.Getenv("connection_string")
		database, err = gorm.Open(os.Getenv("db"), connection_string)
	} else {
		connection_string = os.Getenv("connection_string_test")
		database, err = gorm.Open(os.Getenv("db_test"), connection_string)
	}
	if err != nil {
		log.Fatalf("error connection to data base: %v", err)
	}

	if os.Getenv("debug") == "true" {
		database.LogMode(true)
	}

	if os.Getenv("AutoMigrateDB") == "true" {
		database.AutoMigrate(&model.Bank{}, &model.Account{}, &model.PixKey{}, &model.Transaction{})
	}
	return database
}
