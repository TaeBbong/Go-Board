package utils

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB() {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	fmt.Println("init from utils.go")
	fmt.Println(db)
	fmt.Println(err)
}
