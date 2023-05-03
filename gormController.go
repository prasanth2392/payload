package main

import (
	//"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	//"github.com/jinzhu/gorm"
	"gorm.io/gorm"
)

var DNS string = dsn("payload")

func GormDBConnection() *gorm.DB {
	db, err := gorm.Open(mysql.Open(DNS), &gorm.Config{})
	if err != nil {
		log.Printf("Can't connect to sql")
		return nil
	}
	return db
}

func GormDBSetup() {
	db := GormDBConnection()
	//Creating payload_orders and order_items tables
	db.AutoMigrate(&PayloadOrder{}, &OrderItem{})

	jf, err := os.Open("pay123.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jf.Close()

	//loading json data into bytecode
	byteValue, err := ioutil.ReadAll(jf)
	if err != nil {
		fmt.Println(err)
	}

	var neworder PayloadOrders

	//Decoding json file bytecode to an PayloadOrders object
	json.Unmarshal(byteValue, &neworder)

	//Inserting json data into database
	for i := 0; i < len(neworder.Orderlist); i++ {
		db.Create(&neworder.Orderlist[i])
	}
}

func GormAddOrder(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text")
	db := GormDBConnection()
	var order PayloadOrder
	json.NewDecoder(request.Body).Decode(&order)
	db.Create(&order)
	json.NewEncoder(response).Encode("Order Added Successfully")
}

func GetPayloadOrders(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "json")
	db := GormDBConnection()
	var orders []PayloadOrder
	db.Model(&PayloadOrder{}).Preload("OrderItems").Find(&orders)
	json.NewEncoder(response).Encode(orders)
}

func GetPayloadOrderById(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "json")
	db := GormDBConnection()
	params := mux.Vars(request)
	var id string = params["id"]
	var order PayloadOrder
	db.Preload("OrderItems").Where("id = ?", id).First(&order)
	json.NewEncoder(response).Encode(order)
}

func GetPayloadOrderByStatus(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "json")
	db := GormDBConnection()
	params := mux.Vars(request)
	var status string = params["status"]
	var order PayloadOrder
	//db.First(&order, "?id = ", params["id"])
	db.Preload("OrderItems").Where("status = ?", status).First(&order)
	json.NewEncoder(response).Encode(order)
}

func ModifyPayloadOrderStatus(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text")
	db := GormDBConnection()
	params := mux.Vars(request)
	var id string = params["id"]
	var order PayloadOrder
	db.Preload("OrderItems").Where("id = ?", id).First(&order)
	json.NewDecoder(request.Body).Decode(&order)
	db.Save(&order)
	json.NewEncoder(response).Encode("Status Updated Successfully")
}
