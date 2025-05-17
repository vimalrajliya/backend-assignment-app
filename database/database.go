package database

import (
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/vimalrajliya/backend-assignment-app/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func ConnectDB() {

	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	println("Connection Opened to Database")
	database.AutoMigrate(&models.User{})
	DB = Dbinstance{Db: database}
}

var Client *redis.Client

func ConnectRedis() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	println("Connection Opened to Redis")
}
