package models

import (
	"time"
)

// OrderInfo displayable order info
type OrderInfo struct {
	CreatedAt       time.Time
	OrderName       string
	CustomerID      string
	CustomerName	string
	CustomerCompany	string
	TotalAmount     float64
	DeliveredAmount float64
}

func (oi *OrderInfo) SetCustomerInfo(ci CustomerInfo) {
	oi.CustomerName = ci.Name
	oi.CustomerCompany = ci.CompanyName
}

// OrderInfoQuery query for orders
type OrderInfoQuery struct {
	PartOfOrderName   string
	PartOfProductName string
	DateFrom          string // conforms to RFC3339
	DateTill          string // conforms to RFC3339
}

// CustomerInfo customer info
type CustomerInfo struct {
	Name string `bson:"name,omitempty"`
	CompanyID int `bson:"company_id,omitempty"`
	CompanyName string `bson:"company_name,omitempty"`
}