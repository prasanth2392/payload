package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func gormHandlers() {
	r := mux.NewRouter()

	//to fetch all the existing order data
	r.HandleFunc("/gorm/get/orders", GetPayloadOrders).Methods("GET")

	//to fetch a single order based on id
	r.HandleFunc("/gorm/get/order/{id}", GetPayloadOrderById).Methods("GET")

	//to fetch order details based on status
	r.HandleFunc("/gorm/get/order/{status}", GetPayloadOrderByStatus).Methods("GET")

	//to add a new order
	r.HandleFunc("/gorm/add/order", GormAddOrder).Methods("POST")

	//to update order status based on id
	r.HandleFunc("/gorm/update/order/{id}", ModifyPayloadOrderStatus).Methods("PUT")

	//opening service on port 8080
	http.ListenAndServe(":8080", r)
}

func mongoHandlers() {
	r := mux.NewRouter()

	//to fetch all the existing order data
	r.HandleFunc("/mongo/get/orders", GetOrderDetailsMongo).Methods("GET")

	//to fetch a single order based on id
	r.HandleFunc("/mongo/get/order/{id}", GetOrderByIdMongo).Methods("GET")

	//to fetch order details based on status
	r.HandleFunc("/mongo/get/orders/{status}", GetOrdersByStatusMongo).Methods("GET")

	//to add a new order
	r.HandleFunc("/mongo/add/order", AddOrderMongo).Methods("POST")

	//to update order status based on id
	r.HandleFunc("/mongo/modify/order/{id}", UpdateOrderByStatusMongo).Methods("PUT")

	//opening service on port 8080
	http.ListenAndServe(":8080", r)
}

func sqlHandlers() {
	r := mux.NewRouter()

	//to create the tables Orders and Items and to Add Values in it from pay123.json file
	r.HandleFunc("/builddb", DbSetUp).Methods("GET")

	//to fetch all the existing order data
	r.HandleFunc("/get/orders", GetOrders).Methods("GET")

	//to fetch the order details by a specific order id
	r.HandleFunc("/get/order/id/{id}", GetOrderById).Methods("GET")

	//to fetch the order details by a specific order status
	r.HandleFunc("/get/orders/status/{status}", GetOrdersByStatus).Methods("GET")

	//to add new order details into the database from request body
	r.HandleFunc("/post/addorder", AddOrder).Methods("POST")

	//to update the status of a specific order from request body
	r.HandleFunc("/modify/order/id/{id}", UpdateStatusById).Methods("PUT")

	http.ListenAndServe(":8080", r)
}
