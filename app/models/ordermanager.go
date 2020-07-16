package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
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
	om.db = cfg.PostgresDB
	return om
}

func (mgr OrderManager) generateWhereClause(q *OrderInfoQuery) (string, []interface{}, error) {
	conditions := []string{}
	values := []interface{}{}
	argNo := 1
	if q.PartOfName != "" {
		cond := fmt.Sprintf("(lower(order_name) like lower($%d) or lower(product) like lower($%d))", argNo, argNo)
		argNo++
		conditions = append(conditions, cond)
		values = append(values, "%" + q.PartOfName + "%")
	}
	if q.DateFrom != "" {
		if _, err := time.Parse(time.RFC3339, q.DateFrom); err == nil {
			cond := fmt.Sprintf("created_at >= $%d", argNo)
			argNo++
			conditions = append(conditions, cond)
			values = append(values, q.DateFrom)
		} else {
			return "", nil, utils.NewIllegalArgument("time values should conform to RFC3339")
		}
	}
	if q.DateTill != "" {
		if _, err := time.Parse(time.RFC3339, q.DateTill); err == nil {
			cond := fmt.Sprintf("created_at < $%d", argNo)
			argNo++
			conditions = append(conditions, cond)
			values = append(values, q.DateTill)
		} else {
			return "", nil, utils.NewIllegalArgument("time values should conform to RFC3339")
		}
	}
	clause := ""
	if len(conditions) > 0 {
		clause = "where " + strings.Join(conditions, " and ")
	}
	return clause, values, nil
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
func (mgr OrderManager) Search(qry *OrderInfoQuery, pg *Pagination) ([]OrderInfo, error) {
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

	whereClause, values, err := mgr.generateWhereClause(qry)
	if (err != nil) {
		return nil, err
	}
	limitOffset, values := mgr.generateLimitOffset(pg, values)

	sql := fmt.Sprintf(sqlFmt, whereClause, limitOffset)

	stmt, err := mgr.db.Prepare(sql)
	if err != nil {
		return nil, fmt.Errorf("Preparing sql statement failed. %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(values...)
	if err != nil {
		return nil, fmt.Errorf("Querying database failed. %w", err)
	}

	orderInfos := []OrderInfo{}

	for rows.Next() {
		var oi OrderInfo
		err := rows.Scan(&oi.OrderID, &oi.OrderName, &oi.CreatedAt, &oi.CustomerID, &oi.TotalAmount, &oi.DeliveredAmount)
		if (err != nil) {
			log.Printf("Structuring returned data failed.", err.Error())
		}
		orderInfos = append(orderInfos, oi)
	}
	return orderInfos, nil
}
