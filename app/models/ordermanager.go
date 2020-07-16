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

func (mgr OrderManager) generateWhereClause(q *OrderInfoQuery) (string, []interface{}) {
	conditions := []string{}
	values := []interface{}{}
	argNo := 1
	if q.PartOfName != "" {
		cond := fmt.Sprintf("(lower(order_name) like lower($%d) or lower(product) like lower($%d))", argNo, argNo)
		argNo++
		conditions = append(conditions, cond)
		values = append(values, "%" + q.PartOfName + "%")
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

func (mgr OrderManager) generateLimitOffset(pg *Pagination, values []interface{}) (string, []interface{}) {
	argNo := len(values) + 1
	sql := fmt.Sprintf("limit $%d offset $%d", argNo, argNo+1)
	argNo++
	limit := pg.PageSize
	offset := (pg.PageNo - 1) * pg.PageSize
	values = append(values, limit, offset)
	return sql, values
}

// Search search for orders
func (mgr OrderManager) Search(qry *OrderInfoQuery, pg *Pagination) []OrderInfo {
	sqlFmt := `
		select O.id, O.order_name, O.created_at, O.customer_id, 
			sum(OI.quantity * price_per_unit) as total_amount, 
			sum(delivered_quantity * price_per_unit) as delivered_amount
		from orders as O
			join order_items as OI on O.id = OI.order_id
			join deliveries as D on D.order_item_id = OI.id
		%s
		group by O.id
		%s`

	whereClause, values := mgr.generateWhereClause(qry)
	limitOffset, values := mgr.generateLimitOffset(pg, values)

	sql := fmt.Sprintf(sqlFmt, whereClause, limitOffset)

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
		rows.Scan(&oi.OrderID, &oi.OrderName, &oi.CreatedAt, &oi.CustomerID, &oi.TotalAmount, &oi.DeliveredAmount)
		orderInfos = append(orderInfos, oi)
	}
	return orderInfos
}
