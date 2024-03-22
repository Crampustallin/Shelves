package database

import (
	"database/sql"
	"strconv"
	"strings"

	_ "github.com/lib/pq"

	"github.com/Crampustallin/Shelves/internal/config"
)

type DB interface {
	Query(query string, args ...any) (*sql.Rows, error)
}


func NewDB() (*sql.DB, error) {
	conf := config.NewConfig()
	connectionString := conf.ConnectionString()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func QueryOrders(db DB, orderNumbers []string) (*sql.Rows, error) {
	query := `SELECT 
	orders.order_number, 
	products.product_name,
	products.id as product_id,
	order_summaries.quantity, 
	main_shelf.shelve_name AS main_shelf_name, 
	COALESCE(
		(SELECT STRING_AGG(secondary_shelf.shelve_name, ',') 
		FROM product_shelves 
		JOIN shelves AS secondary_shelf ON product_shelves.shelve_id = secondary_shelf.ID 
		WHERE product_shelves.product_id = products.ID 
		AND product_shelves.is_main = FALSE), 
		NULL
	) AS secondary_shelves_names
	FROM orders
	JOIN order_summaries ON orders.ID = order_summaries.order_id
	JOIN products ON order_summaries.product_id = products.ID
	LEFT JOIN product_shelves main_ps ON products.ID = main_ps.product_id AND main_ps.is_main = TRUE
	LEFT JOIN shelves main_shelf ON main_ps.shelve_id = main_shelf.ID
	where orders.order_number in (`

	orderLen := len(orderNumbers)
	placeholders := make([]string, orderLen, orderLen) 
	interfaceSlice := make([]interface{}, orderLen, orderLen) 

	for i := range placeholders {
		placeholders[i] = "$" + strconv.Itoa(i+1)
		interfaceSlice[i] = orderNumbers[i]
	}
	query += strings.Join(placeholders, ",") + ") order by main_shelf_name, product_id;"
	return db.Query(query, interfaceSlice...) 
}

// TODO: select product_id, quantity from orders, order_summaries where orders.order_number == number and orders.id == order_summaries
// select product_name from products where products.product_id == product ids from order_summaries
// select shelves.id shelves.shelve_name from product_shelves, shelves where product_id == product_id and shelve_id == shelves.id 
// retrieve subshevles
// save all data
