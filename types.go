package main

type ITEM struct {
	Id          string  `json:"id"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Qty         int     `json:"quantity"`
}

//Order Structure

type Order struct {
	Id       string `json:"id"`
	Status   string `json:"status"`
	Items    []ITEM `json:"items"`
	Total    int    `json:"total"`
	Currency string `json:"currencyUnit"`
}

//ORDERS structure containing an array of Order

type Orders struct {
	Orderlist []Order `json:"order"`
}
