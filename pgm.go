package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//ITEM Structure

type ITEM struct {
	Id          string  `json:"Id"`
	Description string  `json:"Description"`
	Price       float32 `json:"Price"`
	Qty         int     `json:"Quantity"`
}

//Order Structure

type Order struct {
	Id       string `json:"Id"`
	Status   string `json:"Status"`
	Items    []ITEM `json:"Items"`
	Total    int    `json:"Total"`
	Currency string `json:"CurrencyUnit"`
}

//ORDERS structure containing an array of Order

type Orders struct {
	Orders []Order `json:"Order"`
}

//DataBase Connectivity Credentials

const (
	username = "root"
	password = "root"
	hostname = "127.0.0.1:3306"
	dbname   = "work"
)

//To return DB connection String

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

//To connect to database and return a DB object

func dbConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn(""))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return nil, err
	}
	//defer db.Close()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return nil, err
	}
	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return nil, err
	}
	log.Printf("rows affected %d\n", no)

	db.Close()
	db, err = sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return nil, err
	}
	//defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return nil, err
	}
	log.Printf("Connected to DB %s successfully\n", dbname)
	return db, nil
}

//To create ITEMS table in the database

func createItemTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS items(item_id int, description text, 
        price float, quantity int,order_id varchar(20), 
		primary key(item_id,order_id),foreign key(order_id) references orders(order_id))`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating product table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return err
	}
	log.Printf("Rows affected when creating table: %d", rows)
	return nil
}

//To create ORDERS table in the database

func createOrderTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS orders(order_id varchar(20) primary key, status text, total float, currency_unit text)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating order table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return err
	}
	log.Printf("Rows affected when creating table: %d", rows)
	return nil
}

//To insert records in the ORDERS table

func insert_orders(db *sql.DB, ord Order) error {
	query := "INSERT INTO orders VALUES (?, ?,?,?)"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, ord.Id, ord.Status, ord.Total, ord.Currency)
	if err != nil {
		log.Printf("Error %s when inserting row into ORDERS table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}

	for i := 0; i < len(ord.Items); i++ {
		insert_items(db, ord.Items[i], ord.Id)
	}
	log.Printf("%d ORDER added ", rows)
	return nil
}

//To insert records in the ITEMS table

func insert_items(db *sql.DB, item ITEM, id string) error {
	query := "INSERT INTO items VALUES (?, ?,?,?,?)"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, item.Id, item.Description, item.Price, item.Qty, id)
	if err != nil {
		log.Printf("Error %s when inserting row into ITEMS table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d ITEMS Added ", rows)
	return nil
}

//To update the Order Status based on Order ID

func update_status(db *sql.DB, id string, status string) error {
	query := "UPDATE ORDERS SET STATUS = ? WHERE ORDER_ID = ?"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, status, id)
	if err != nil {
		log.Printf("Error %s when updating order status", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d Order Updated ", rows)
	return nil
}

//To fetch Order Details based on Order ID

func search_orderid(db *sql.DB, id string) error {
	query1 := "SELECT * FROM ORDERS WHERE ORDER_ID = ?"
	res1, err := db.Query(query1, id)
	if err != nil {
		log.Printf("Error %s while Searching Orders by ID", err)
		return err
	}
	query2 := "SELECT * FROM ITEMS WHERE ORDER_ID = ?"
	res2, err := db.Query(query2, id)
	if err != nil {
		log.Printf("Error %s while Searching Orders by ID", err)
		return err
	}
	var ord1 Order
	var item ITEM
	for res1.Next() {

		err = res1.Scan(&ord1.Id, &ord1.Status, &ord1.Total, &ord1.Currency)
		if err != nil {
			log.Printf("%s", err)
			return err
		}
		log.Printf("Order Details of %s", id)
		log.Printf("Status : %s", ord1.Status)
		log.Printf("Total : %d", ord1.Total)
		log.Printf("Currency Unit : %s", ord1.Currency)
	}
	item_count := 0
	for res2.Next() {
		item_count++
		err = res2.Scan(&item.Id, &item.Description, &item.Price, &item.Qty, &id)
		if err != nil {
			log.Printf("%s", err)
			return err
		}
		log.Printf("Item %d", item_count)
		log.Printf("Item Id : %s", item.Id)
		log.Printf("Description : %s", item.Description)
		log.Printf("Item Price : %f", item.Price)
		log.Printf("Item quantity : %d", item.Qty)
	}
	return nil

}

//To fetch Order Details based on Order Status

func search_orderstatus(db *sql.DB, status string) error {
	query1 := "SELECT * FROM ORDERS WHERE STATUS = ?"
	res1, err := db.Query(query1, status)
	if err != nil {
		log.Printf("Error %s while Searching Orders by STATUS", err)
		return err
	}

	var ord1 Order
	var item ITEM
	for res1.Next() {

		err = res1.Scan(&ord1.Id, &ord1.Status, &ord1.Total, &ord1.Currency)
		if err != nil {
			log.Printf("%s", err)
			return err
		}
		log.Printf("%s Order Details", status)
		log.Printf("Order ID : %s", ord1.Id)

		query2 := "SELECT * FROM ITEMS WHERE ORDER_ID = ?"
		res2, err := db.Query(query2, ord1.Id)
		if err != nil {
			log.Printf("Error %s while Searching Orders by ID", err)
			return err
		}
		item_count := 0
		for res2.Next() {
			item_count++
			err = res2.Scan(&item.Id, &item.Description, &item.Price, &item.Qty, &ord1.Id)
			if err != nil {
				log.Printf("%s", err)
				return err
			}
			log.Printf("Item %d", item_count)
			log.Printf("Item Id : %s", item.Id)
			log.Printf("Description : %s", item.Description)
			log.Printf("Item Price : %f", item.Price)
			log.Printf("Item quantity : %d", item.Qty)
		}

		log.Printf("Total : %d", ord1.Total)
		log.Printf("Currency Unit : %s", ord1.Currency)
	}

	return nil
}

// Main Entry

func main() {
	//DataBase Connectivity

	db, err := dbConnection()
	if err != nil {
		log.Printf("Error %s when getting db connection", err)
		return
	}
	defer db.Close()
	log.Printf("Successfully connected to database")

	//Loading JSON file

	jsonfile, err := os.Open("pay123.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("File Opened")

	defer jsonfile.Close()
	byteValue, _ := ioutil.ReadAll(jsonfile)

	var ord Orders

	//Decoding json file bytecode to an Orders object
	json.Unmarshal(byteValue, &ord)

	//Inserting data into the database from the Orders object

	for i := 0; i < len(ord.Orders); i++ {
		insert_orders(db, ord.Orders[i])
	}

	//To update the status of an Order based on Order Id
	/*
		update_status(db, "abcdef-2", "payment pending")
	*/

	//To search Order Details based on Order Id
	/*
		search_orderid(db, "abcdef-2")
	*/

	//To search Order Details based on Order Status
	/*
		search_orderstatus(db, "payment pending")
	*/
}
