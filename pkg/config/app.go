package config

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var (
	db *gorm.DB
)

func Connect(){
	connectionString := "uri"

	d, err := gorm.Open("postgres", connectionString)
	
	if err  != nil{
		log.Fatal(err)
	}
	
	db = d
}

func GetDB() *gorm.DB{
	return db
}