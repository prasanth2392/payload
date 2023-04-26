package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

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
