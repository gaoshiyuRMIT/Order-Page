package models

import (
	"time"
)

// CustomerCompany customer company
type CustomerCompany struct {
	CompanyID   int
	CompanyName string
}

// Customer customer
type Customer struct {
	UserID      string
	Login       string
	Password    string
	Name        string
	CompanyID   int
	CreditCards []string
}

// OrderInfo displayable order info
type OrderInfo struct {
	CreatedAt       time.Time
	OrderName       string
	CustomerID      string
	TotalAmount     float64
	DeliveredAmount float64
}
