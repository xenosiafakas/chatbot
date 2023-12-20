package models

import (
	"github.com/jinzhu/gorm"
	"github.com/xenosiafakas/review-chatbot/pkg/config"
)


var (
	db *gorm.DB
)

	type Product struct {
		gorm.Model
		Name string `json:"name"`
	}

	type Customer struct {
		gorm.Model
		Name    string   `json:"name"`
		Email   string   `json:"email"`
	}

	type Review struct {	
		gorm.Model
		Rating     int      `json:"rating"`
		Comment    string   `json:"comment"`
		CustomerName string `json:"CustomerName"`
		ProductName string `json:"ProductName"`
	}


func init(){
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Product{},&Customer{},&Review{})
}

func (r *Review) CreateReview() *Review{
	db.NewRecord(r)
	db.Create(&r)
	return r
}

func (c *Customer) CreateCustomer() *Customer{
	db.NewRecord(c)
	db.Create(&c)
	return c
}

func (p *Product) CreateProduct() *Product{
	db.NewRecord(p)
	db.Create(&p)
	return p
}