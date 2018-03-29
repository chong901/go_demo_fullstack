package db

import (
	"github.com/aaa59891/mosi_demo_go/configs"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
)

var DB *gorm.DB
var err error

func init() {
	databaseConfig := configs.GetConfig().Database

	DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", databaseConfig.Account, databaseConfig.Password, databaseConfig.Host, databaseConfig.Port, databaseConfig.DatabaseName))

	if err != nil {
		log.Fatal(err)
	}

	mosiGoEnv := os.Getenv("MOSI_GO")

	if mosiGoEnv == "dev" {
		DB.LogMode(true)
	}

	fmt.Println("Connected to database.")
}
