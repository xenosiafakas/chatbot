package routes

import (
	"github.com/gorilla/mux"
	"github.com/xenosiafakas/review-chatbot/pkg/controllers"
)

var ChatBotRoutes = func(router *mux.Router) {

	router.HandleFunc("/slack/event", controllers.SlackEvent).Methods("POST")

	router.HandleFunc("/reviews", controllers.GetAllReviews).Methods("GET")
	router.HandleFunc("/customer", controllers.CreateCustomer).Methods("POST")
	router.HandleFunc("/product", controllers.CreateProduct).Methods("POST")
	router.HandleFunc("/reviews/{reviewId}", controllers.DeleteReview).Methods("DELETE")
}
