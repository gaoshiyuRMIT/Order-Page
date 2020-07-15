package models

import (
	"database/sql"
	"log"

	"app/utils"
)

// OrderManager order table manager
type OrderManager struct {
	db *sql.DB
}

// NewOrderManager constructor
func NewOrderManager(cfg *utils.ConfigReader) *OrderManager {
	om := &OrderManager{}
	om.db = cfg.GetPostgresDB()
	return om
}

// TableName table name in db
func (mgr OrderManager) TableName() string {
	return "order_items"
}

// GetAll get all orders
func (mgr OrderManager) GetAll() []OrderInfo {
	sql := `
		select O.order_name, O.created_at, O.customer_id, 
			sum(OI.quantity * price_per_unit) as total_amount, 
			sum(delivered_quantity * price_per_unit) as delivered_amount
		from orders as O
			join order_items as OI on O.id = OI.order_id
			join deliveries as D on D.order_item_id = OI.id
		group by O.id`
	stmt, err := mgr.db.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}

	orderInfos := []OrderInfo{}

	for rows.Next() {
		var oi OrderInfo
		rows.Scan(&oi.OrderName, &oi.CreatedAt, &oi.CustomerID, &oi.TotalAmount, &oi.DeliveredAmount)
		orderInfos = append(orderInfos, oi)
	}
	return orderInfos
}
