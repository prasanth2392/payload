package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

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

//to display all the existing orders

func GetOrders(w http.ResponseWriter, r *http.Request) {
	db, err := DbConnection()
	defer db.Close()
	if err != nil {
		log.Printf(err.Error())
	} else {
		query1 := "SELECT * FROM ORDERS"
		res1, err := db.Query(query1)
		if err != nil {
			log.Printf("Error %s while fetching all orders", err)
			return
		}

		var ord1 Order
		var item ITEM
		var fullorder Orders
		for res1.Next() {

			err = res1.Scan(&ord1.Id, &ord1.Status, &ord1.Total, &ord1.Currency)
			if err != nil {
				log.Printf("%s", err)
				return
			}
			query2 := "SELECT * FROM ITEMS WHERE ORDER_ID = ?"
			res2, err := db.Query(query2, ord1.Id)
			if err != nil {
				log.Printf("Error %s while Searching Orders by ID", err)
				return
			}

			ord1.Items = nil
			for res2.Next() {
				err = res2.Scan(&item.Id, &item.Description, &item.Price, &item.Qty, &ord1.Id)
				if err != nil {
					log.Printf("%s", err)
					return
				}
				ord1.Items = append(ord1.Items, item)
			}
			fullorder.Orderlist = append(fullorder.Orderlist, ord1)

		}
		w.Header().Set("Content-Type", "json")
		json.NewEncoder(w).Encode(fullorder.Orderlist)

	}
}

//to display order details of a specific id

func GetOrderById(w http.ResponseWriter, r *http.Request) {
	db, err := DbConnection()
	defer db.Close()
	if err != nil {
		log.Printf(err.Error())
	} else {
		vars := mux.Vars(r)
		searchid := vars["id"]
		fmt.Print(searchid)
		query1 := "SELECT * FROM ORDERS where order_id = ?"
		res1, err := db.Query(query1, searchid)
		if err != nil {
			log.Printf("Error %s while Searching Orders by ID", err)
			return
		}

		var ord1 Order
		var item ITEM
		var fullorder Orders
		for res1.Next() {
			err = res1.Scan(&ord1.Id, &ord1.Status, &ord1.Total, &ord1.Currency)
			if err != nil {
				log.Printf("%s", err)
				return
			}
			query2 := "SELECT * FROM ITEMS WHERE ORDER_ID = ?"
			res2, err := db.Query(query2, ord1.Id)
			if err != nil {
				log.Printf("Error %s while Searching Orders by ID", err)
				return
			}
			for res2.Next() {
				err = res2.Scan(&item.Id, &item.Description, &item.Price, &item.Qty, &ord1.Id)
				if err != nil {
					log.Printf("%s", err)
					return
				}
				ord1.Items = append(ord1.Items, item)
			}
			fullorder.Orderlist = append(fullorder.Orderlist, ord1)

		}
		w.Header().Set("Content-Type", "json")
		json.NewEncoder(w).Encode(fullorder.Orderlist)

	}
}

func GetOrdersByStatus(w http.ResponseWriter, r *http.Request) {
	db, err := DbConnection()
	defer db.Close()
	if err != nil {
		log.Printf(err.Error())
	} else {
		vars := mux.Vars(r)
		searchstatus := vars["status"]
		query1 := "SELECT * FROM ORDERS where status = ?"
		res1, err := db.Query(query1, searchstatus)
		if err != nil {
			log.Printf("Error %s while Searching Orders by status", err)
			return
		}

		var ord1 Order
		var item ITEM
		var fullorder Orders
		for res1.Next() {

			err = res1.Scan(&ord1.Id, &ord1.Status, &ord1.Total, &ord1.Currency)
			if err != nil {
				log.Printf("%s", err)
				return
			}
			query2 := "SELECT * FROM ITEMS WHERE ORDER_ID = ?"
			res2, err := db.Query(query2, ord1.Id)
			if err != nil {
				log.Printf("Error %s while Searching Orders by ID", err)
				return
			}
			for res2.Next() {
				err = res2.Scan(&item.Id, &item.Description, &item.Price, &item.Qty, &ord1.Id)
				if err != nil {
					log.Printf("%s", err)
					return
				}
				ord1.Items = append(ord1.Items, item)
			}
			fullorder.Orderlist = append(fullorder.Orderlist, ord1)

		}
		w.Header().Set("Content-Type", "json")
		json.NewEncoder(w).Encode(fullorder.Orderlist)

	}
}

// to add the details of a new order

func AddOrder(w http.ResponseWriter, r *http.Request) {
	db, err := DbConnection()
	defer db.Close()
	if err != nil {
		log.Printf(err.Error())
	} else {
		var neword Order
		json.NewDecoder(r.Body).Decode(&neword)
		query1 := "INSERT INTO ORDERS VALUES (?,?,?,?)"
		res, err := db.Exec(query1, neword.Id, neword.Status, neword.Total, neword.Currency)
		if err != nil {
			log.Printf(err.Error())
		} else {
			_, err := res.LastInsertId()
			if err != nil {
				json.NewEncoder(w).Encode("[error : order not inserted]")
			} else {
				itemcount := len(neword.Items)
				for i := 0; i < itemcount; i++ {
					curitem := neword.Items[i]
					query2 := "INSERT INTO ITEMS VALUES(?,?,?,?,?)"
					_, err := db.Exec(query2, curitem.Id, curitem.Description, curitem.Price, curitem.Qty, neword.Id)
					if err != nil {
						log.Printf(err.Error())
					} else {
						if err != nil {
							json.NewEncoder(w).Encode("[error : order not inserted]")
						} else {
							json.NewEncoder(w).Encode("[error : order inserted successfully]")
						}
					}
				}
			}
		}

	}
}

// to update the status of an order

func UpdateStatusById(w http.ResponseWriter, r *http.Request) {
	db, err := DbConnection()
	defer db.Close()
	curord := Order{}
	json.NewDecoder(r.Body).Decode(&curord)
	fmt.Printf(curord.Id, curord.Status)
	if err != nil {
		log.Printf(err.Error())
	} else {
		vars := mux.Vars(r)
		id := vars["id"]
		sql := "UPDATE ORDERS SET STATUS = ? WHERE ORDER_ID = ?"
		res, err := db.Exec(sql, curord.Status, id)
		if err != nil {
			log.Printf(err.Error())
		} else {
			_, err := res.RowsAffected()
			if err != nil {
				w.Write([]byte("Record not updated"))
			} else {
				w.Write([]byte("Record updated successfully"))
			}

		}

	}
}
