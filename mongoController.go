package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

/*
func getMongoConnection() {
	_, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	uri := "mongodb+srv://prasanthquest:pras%40mongo123@cluster0.odrho4i.mongodb.net/test"
	client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

}*/

func AddOrderMongo(response http.ResponseWriter, request *http.Request) {
	var neworder Order

	json.NewDecoder(request.Body).Decode(&neworder)
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	uri := "mongodb+srv://prasanthquest:pras%40mongo123@cluster0.odrho4i.mongodb.net/test"
	client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	collection := client.Database("work").Collection("payload")
	collection.InsertOne(ctx, neworder)

	json.NewEncoder(response).Encode("order inserted successfully]")
	response.Header().Set("Content-Type", "text")
}

func GetOrderDetailsMongo(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "json")
	var orderCollection []Order
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()

	uri := "mongodb+srv://prasanthquest:pras%40mongo123@cluster0.odrho4i.mongodb.net/test"
	client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	collection := client.Database("work").Collection("payload")

	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)

	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var order Order
		cursor.Decode(&order)
		orderCollection = append(orderCollection, order)
	}
	json.NewEncoder(response).Encode(orderCollection)

}

func GetOrderByIdMongo(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "json")
	var orderData Order
	vars := mux.Vars(request)
	searchid := vars["id"]
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()

	uri := "mongodb+srv://prasanthquest:pras%40mongo123@cluster0.odrho4i.mongodb.net/test"
	client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	collection := client.Database("work").Collection("payload")
	filter := bson.D{{"id", searchid}}
	err := collection.FindOne(ctx, filter).Decode(&orderData)
	if err != nil {
		//log.Printf(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(response).Encode(orderData)

}

func GetOrdersByStatusMongo(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "json")
	var orderData Order
	vars := mux.Vars(request)
	statuskey := vars["status"]
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()

	uri := "mongodb+srv://prasanthquest:pras%40mongo123@cluster0.odrho4i.mongodb.net/test"
	client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	collection := client.Database("work").Collection("payload")
	filter := bson.D{{"status", statuskey}}
	err := collection.FindOne(ctx, filter).Decode(&orderData)
	if err != nil {
		//log.Printf(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(response).Encode(orderData)

}

func UpdateOrderByStatusMongo(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text")
	var orderData Order
	vars := mux.Vars(request)
	keyid, _ := vars["id"]

	json.NewDecoder(request.Body).Decode(&orderData)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()

	uri := "mongodb+srv://prasanthquest:pras%40mongo123@cluster0.odrho4i.mongodb.net/test"
	client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	collection := client.Database("work").Collection("payload")
	filter := bson.D{{"id", keyid}}
	update := bson.M{"$set": bson.M{"status": orderData.Status}}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		//log.Printf(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
	}
	msg := fmt.Sprintf("updated successfully", result.ModifiedCount, "rows")
	json.NewEncoder(response).Encode(msg)

}
