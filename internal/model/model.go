package model

import "time"

type Order struct {
	OrderUID    string    `json:"order_uid"`
	TrackNumber string    `json:"track_number"`
	Entry       string    `json:"entry"`
	Locale      string    `json:"locale"`
	CustomerID  string    `json:"customer_id"`
	DateCreated time.Time `json:"date_created"`
	OofShard    string    `json:"oof_shard"`
	Delivery    Delivery
	Payment     Payment
	OrderItems  []OrderItem
}

// Delivery структура для доставки
type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

// Payment структура для платежей
type Payment struct {
	Transaction  string    `json:"transaction"`
	RequestID    string    `json:"request_id"`
	Currency     string    `json:"currency"`
	Provider     string    `json:"provider"`
	Amount       int       `json:"amount"`
	PaymentDate  time.Time `json:"payment_dt"`
	Bank         string    `json:"bank"`
	DeliveryCost int       `json:"delivery_cost"`
	GoodsTotal   int       `json:"goods_total"`
	CustomFee    int       `json:"custom_fee"`
}

// OrderItem структура для элементов заказа
type OrderItem struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NMID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}