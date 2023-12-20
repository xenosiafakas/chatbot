package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/xenosiafakas/review-chatbot/pkg/config"
	"github.com/xenosiafakas/review-chatbot/pkg/models"
	"github.com/xenosiafakas/review-chatbot/pkg/utils"
)

var db = config.GetDB()

func SlackEvent(w http.ResponseWriter, r *http.Request) {
	body,_:= io.ReadAll(r.Body)

	var eventData map[string]interface{}

	json.Unmarshal(body, &eventData)
	fmt.Println("got it: ", string(body))


	if eventData["type"] == "receivedProduct"{
		
	var customer models.Customer
	var product models.Product

	if err:= db.Where("name = ?", eventData["user"].(map[string]interface{})["name"]).First(&customer).Error; err != nil {
		
		http.Error(w, "Error finding user in the database", http.StatusInternalServerError)
		return
	} 

	if err:= db.Where("name = ?", eventData["product"].(map[string]interface{})["name"]).First(&product).Error; err != nil {
		
		http.Error(w, "Error finding product in the database", http.StatusInternalServerError)
		return
	}
	fmt.Println("Product devlivered! : ", string(body))
	
	go utils.GetData(customer,product)
		
	} else {
		http.Error(w, "Product has not yet been delivered", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}

func GetAllReviews(w http.ResponseWriter, r *http.Request) {
	var reviews []models.Review
	db.Find(&reviews)

	res,_ := json.Marshal(reviews)
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateCustomer(w http.ResponseWriter, r *http.Request){
	CreateCustomer:= &models.Customer{}
	utils.ParseBody(r, CreateCustomer)
	b:= CreateCustomer.CreateCustomer()
	res, _ := json.Marshal(b)
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateProduct(w http.ResponseWriter, r *http.Request){
	CreateProduct:= &models.Product{}
	utils.ParseBody(r, CreateProduct)
	b:= CreateProduct.CreateProduct()
	res, _ := json.Marshal(b)
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteReview(w http.ResponseWriter, r *http.Request){
	vars:= mux.Vars(r)
	reviewId := vars["reviewId"]
	ID,_ := strconv.ParseInt(reviewId, 0,0)
	var review models.Review
	db.Where("ID=?",ID).Delete(review)
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
}