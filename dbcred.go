package main

import (
	"fmt"
)

const (
	username = "root"
	password = "root"
	hostname = "127.0.0.1:3306"
	dbname   = "testwork"

	jsonfile = "pay123.json"
)

//To return DB connection String

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=True", username, password, hostname, dbName)
}
