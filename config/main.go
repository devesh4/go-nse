package config

import (
	"os"

	"github.com/devesh44/gin-poc/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EnvVariables struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBName     string
	DBPassword string

	JWTSecret string
	APIPort   string
}
type Queue []map[int]float64

var EnvVar *EnvVariables

type MainConfiguration struct {
	DB *gorm.DB

	Router *gin.Engine
}

var MainConf = &MainConfiguration{}

func LoadMainConfiguration() {

	MainConf.Router = gin.Default()
	pref := "TRADING_"
	EnvVar = &EnvVariables{
		DBHost:     os.Getenv(pref + "DB_HOST"),
		DBPort:     os.Getenv(pref + "DB_PORT"),
		DBUser:     os.Getenv(pref + "DB_USER"),
		DBName:     os.Getenv(pref + "DB_NAME"),
		DBPassword: os.Getenv(pref + "DB_PASSWORD"),
		JWTSecret:  os.Getenv(pref + "JWT_SECRET"),
		APIPort:    os.Getenv(pref + "API_PORT"),
	}

	MainConf.DB = db.NewPGConnection(&db.DBConf{
		DBHost:     EnvVar.DBHost,
		DBPort:     EnvVar.DBPort,
		DBUser:     EnvVar.DBUser,
		DBName:     EnvVar.DBName,
		DBPassword: EnvVar.DBPassword,
		SSL:        "disable",
	})
}
