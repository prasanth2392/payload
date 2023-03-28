
Database : testwork

Tables : orders , items

As each Order may contain multiple Items in it , we can have two tables as per following structure

Table Description

mysql> desc orders;
+---------------+-------------+------+-----+---------+-------+
| Field         | Type        | Null | Key | Default | Extra |
+---------------+-------------+------+-----+---------+-------+
| order_id      | varchar(20) | NO   | PRI | NULL    |       |
| status        | text        | YES  |     | NULL    |       |
| total         | float       | YES  |     | NULL    |       |
| currency_unit | text        | YES  |     | NULL    |       |
+---------------+-------------+------+-----+---------+-------+
4 rows in set (0.10 sec)

mysql> desc items;
+-------------+-------------+------+-----+---------+-------+
| Field       | Type        | Null | Key | Default | Extra |
+-------------+-------------+------+-----+---------+-------+
| item_id     | int(11)     | NO   | PRI | 0       |       |
| description | text        | YES  |     | NULL    |       |
| price       | float       | YES  |     | NULL    |       |
| quantity    | int(11)     | YES  |     | NULL    |       |
| order_id    | varchar(20) | NO   | PRI |         |       |
+-------------+-------------+------+-----+---------+-------+
5 rows in set (0.01 sec)


JSON file used : pay123.json

This file contains some initial data which is to be stored into the database

We are using mux router on port number 8080 to start server and can avail services on localhost on the following api function handlers:

localhost:8080/builddb --> to create the tables Orders and Items and to Add Values in it from pay123.json file


localhost:8080/orders -->	to fetch all the existing order data (using GET)
	
localhost:8080/orders/id/{id} --> to fetch the order details by a specific order id (using GET)

localhost:8080/orders/status/{status}" --> 	to fetch the order details by a specific order status (using GET)

localhost:8080/orders --> to add new order details into the database from request body (using POST)

localhost:8080/orders/update/{id} --> to update the status of a specific order from request body (using PUT)
