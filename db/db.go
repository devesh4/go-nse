package db

import (
	"fmt"
	"log"

	"github.com/devesh44/gin-poc/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConf struct {
	// host url
	DBHost string

	// port
	DBPort string

	// user
	DBUser string

	// name
	DBName string

	// password
	DBPassword string

	// ssl
	SSL string
}

func NewPGConnection(conf *DBConf) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		conf.DBHost, conf.DBUser, conf.DBPassword, conf.DBName, conf.DBPort, conf.SSL, "Asia/Kolkata")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	db.AutoMigrate(
		model.User{},
		model.Order{},
		model.OrderLines{},
		model.OrderHistory{},
	)
	return db

}
