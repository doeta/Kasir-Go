package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique" json:"username"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
}

type Product struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `json:"name"`
	Price      int       `json:"price"`
	Stock      int       `json:"stock"`
	CategoryID uint      `json:"category_id"`
}

type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

type Transaction struct {
	ID              uint                `gorm:"primaryKey" json:"id"`
	UserID          uint                `json:"user_id"`
	User            User                `gorm:"foreignKey:UserID" json:"user"` 
	TotalAmount     int                 `json:"total_amount"`
	TransactionDate time.Time           `json:"transaction_date"`
	Details         []TransactionDetail `gorm:"foreignKey:TransactionID" json:"details"`
}

type TransactionDetail struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	TransactionID uint    `json:"transaction_id"`
	ProductID     uint    `json:"product_id"`
	Quantity      int     `json:"quantity"`
	Subtotal      int     `json:"subtotal"`
}


type PaymentMethod struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}