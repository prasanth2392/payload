package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"io/ioutil"

	"log"
	"net/http"

	"os"
	//	_ "github.com/go-sql-driver/mysql"
	//	_ "github.com/gorilla/mux"
)

const (
	username = "root"
	password = "root"
	hostname = "127.0.0.1:3306"
	dbname   = "testwork"

	jsonfile = "pay123.json"
)

func Test123() string {
	fmt.Print("okay")
	return fmt.Sprintf("okay")
}

//To return DB connection String

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

//To connect to database and return a DB object

func DbConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn(dbname))

	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return nil, err
	}
	//log.Printf("Connected to DB %s successfully\n", dbname)
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

func insert_orders(db *sql.DB) error {
	jf, err := os.Open("pay123.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("File Opened")
	defer jf.Close()
	byteValue, _ := ioutil.ReadAll(jf)
	var ord Orders

	//Decoding json file bytecode to an Orders object
	json.Unmarshal(byteValue, &ord)

	//Inserting data into the database from the Orders object

	for i := 0; i < len(ord.Orderlist); i++ {
		neword := ord.Orderlist[i]

		query := "INSERT INTO orders VALUES (?, ?,?,?)"
		ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()
		stmt, err := db.PrepareContext(ctx, query)
		if err != nil {
			log.Printf("Error %s when preparing SQL statement", err)
			return err
		}
		defer stmt.Close()
		res, err := stmt.ExecContext(ctx, neword.Id, neword.Status, neword.Total, neword.Currency)
		if err != nil {
			log.Printf("Error %s when inserting row into ORDERS table", err)
			return err
		}
		rows, err := res.RowsAffected()
		if err != nil {
			log.Printf("Error %s when finding rows affected", err)
			return err
		}

		for j := 0; j < len(neword.Items); j++ {
			insert_items(db, neword.Items[j], neword.Id)
		}
		log.Printf("%d ORDER added ", rows)

	}
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

// to create tables and insert values from json file

func DbSetUp(w http.ResponseWriter, r *http.Request) {
	db, err := DbConnection()
	defer db.Close()
	var msg string
	if err != nil {
		log.Printf(err.Error())
	} else {
		err := createOrderTable(db)
		if err != nil {
			log.Printf(err.Error())
		} else {
			msg += "Orders Table Created\n"
		}
		err1 := createItemTable(db)
		if err1 != nil {
			log.Printf(err1.Error())
		} else {
			msg += "Items Table Created\n"
		}
		err2 := insert_orders(db)
		if err2 != nil {
			log.Printf(err.Error())
		} else {
			msg += "insertion done"
		}
	}
	w.Write([]byte(msg))
}
