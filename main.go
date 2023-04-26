package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
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
	r.HandleFunc("/modify/order/update/{id}", UpdateStatusById).Methods("PUT")

	//opening port number 8080 on localhost
	http.ListenAndServe(":8080", r)

}
