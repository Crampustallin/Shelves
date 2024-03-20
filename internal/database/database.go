package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/lib/pq"

	"github.com/Crampustallin/Shelves/internal/config"
)

func NewDB() (*sql.DB, error) {
	conf := config.NewConfig()
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", conf.DBHost, conf.DBPort,
	conf.DBUser, conf.DBPassword, conf.DBName)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func QueryOrders(db *sql.DB, orderNumbers []string) (*sql.Rows, error) {
	query := `SELECT 
	orders.order_number, 
	products.product_name, 
	order_summaries.quantity, 
	main_shelf.shelve_name AS main_shelf_name, 
	COALESCE(
		(SELECT STRING_AGG(secondary_shelf.shelve_name, ', ') 
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
	placeholders := make([]string, len(orderNumbers), len(orderNumbers))
	for i := range placeholders {
		placeholders[i] = "$" + strconv.Itoa(i+1)
	}
	query += strings.Join(placeholders, ",") + ")"
	return db.Query(query, orderNumbers)
}

