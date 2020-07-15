package models

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

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
	defer stmt.Close()

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

func (mgr OrderManager) generateWhereClause(q *OrderInfoQuery) (string, []interface{}) {
	conditions := []string{}
	values := []interface{}{}
	argNo := 1
	if q.PartOfOrderName != "" {
		cond := fmt.Sprintf("lower(order_name) like lower($%d)", argNo)
		argNo++
		conditions = append(conditions, cond)
		values = append(values, "%" + q.PartOfOrderName + "%")
	}
	if q.PartOfProductName != "" {
		cond := fmt.Sprintf("lower(product) like lower($%d)", argNo)
		argNo++
		conditions = append(conditions, cond)
		values = append(values, "%" + q.PartOfProductName, "%")
	}
	if _, err := time.Parse(time.RFC3339, q.DateFrom); err == nil {
		cond := fmt.Sprintf("created_at >= $%d", argNo)
		argNo++
		conditions = append(conditions, cond)
		values = append(values, q.DateFrom)
	}
	if _, err := time.Parse(time.RFC3339, q.DateTill); err == nil {
		cond := fmt.Sprintf("created_at < $%d", argNo)
		argNo++
		conditions = append(conditions, cond)
		values = append(values, q.DateTill)
	}
	clause := ""
	if len(conditions) > 0 {
		clause = "where " + strings.Join(conditions, " and ")
	}
	return clause, values
}

// Search search for orders
func (mgr OrderManager) Search(qry *OrderInfoQuery) []OrderInfo {
	sqlFmt := `
		select O.order_name, O.created_at, O.customer_id, 
			sum(OI.quantity * price_per_unit) as total_amount, 
			sum(delivered_quantity * price_per_unit) as delivered_amount
		from orders as O
			join order_items as OI on O.id = OI.order_id
			join deliveries as D on D.order_item_id = OI.id
		%s
		group by O.id`

	whereClause, values := mgr.generateWhereClause(qry)

	sql := fmt.Sprintf(sqlFmt, whereClause)

	stmt, err := mgr.db.Prepare(sql)
	if err != nil {
		log.Fatal("preparing statement failed,", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(values...)
	if err != nil {
		log.Fatal("querying failed,", err)
	}

	orderInfos := []OrderInfo{}

	for rows.Next() {
		var oi OrderInfo
		rows.Scan(&oi.OrderName, &oi.CreatedAt, &oi.CustomerID, &oi.TotalAmount, &oi.DeliveredAmount)
		orderInfos = append(orderInfos, oi)
	}
	return orderInfos
}
