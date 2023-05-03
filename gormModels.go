package main

import (
	"gorm.io/gorm"
)

type OrderItem struct {
	gorm.Model
	ItemId         string  `json:"id" gorm:"primaryKey"`
	Description    string  `json:"description"`
	Price          float32 `json:"price"`
	Qty            int     `json:"quantity"`
	PayloadOrderId string  `gorm:"size:191" gorm:"primaryKey" gorm:"foreignKey:Id" `
}

//Order Structure

type PayloadOrder struct {
	gorm.Model
	Id         string      `json:"id" gorm:"size:191" gorm:"primaryKey"`
	Status     string      `json:"status"`
	OrderItems []OrderItem `json:"items"`
	Total      int         `json:"total"`
	Currency   string      `json:"currencyUnit"`
}

type PayloadOrders struct {
	Orderlist []PayloadOrder `json:"order"`
}
