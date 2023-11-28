package models

import (
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

type Namespaces struct {
	Name string `json:"name"`
}

type Dataset struct {
	ID      uint   `json:"id" gorm:"primary_key"`
	Name    string `json:"name"`
	Owner   string `json:"owner"`
	Content string `json:"content"`
}

type CreateDatasetInput struct {
	Name    string `json:"name" binding:"required"`
	Owner   string `json:"owner" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func SetupDB() *gorm.DB {
	postgres_user := os.Getenv("POSTGRES_USER")
	postgres_pass := os.Getenv("POSTGRES_PASSWORD")
	postgres_db := os.Getenv("POSTGRES_DB")
	postgres_host := os.Getenv("POSTGRES_HOST")
	postgres_port := os.Getenv("POSTGRES_PORT")

	dsn := "postgres" + "://" + postgres_user + ":" +
		postgres_pass + "@" + postgres_host + ":" + postgres_port +
		"/" + postgres_db

	retries := 10
	log.Println("Trying to connect to DB")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	for err != nil {
		log.Println("Attempt: " + strconv.Itoa(retries))
		log.Println(err)
		if retries > 1 {
			retries--
			time.Sleep(5 * time.Second)
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			log.Println()
			continue
		}
		panic(err)
	}
	log.Println("Establish connection to DB")

	db.AutoMigrate(&Dataset{})
	return db
}
