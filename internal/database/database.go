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

func QueryOrders(db DB, orderNumbers []interface{}) (*sql.Rows, error) {
	query := `SELECT orders.id, orders.order_number FROM orders WHERE orders.order_number IN (`
	formatParamQuery(&query, orderNumbers)
	return db.Query(query, orderNumbers...)
}

func QuerySummaries(db DB, orderIds []interface{}) (*sql.Rows, error) {
	query := `SELECT os.order_id, 
	os.product_id,
	os.quantity
	FROM order_summaries os 
	WHERE os.order_id IN (`
	formatParamQuery(&query, orderIds)
	return db.Query(query, orderIds...)
}

func QueryProducts(db DB, productIds []interface{}) (*sql.Rows, error) {
	query := `SELECT products.id, products.product_name FROM products WHERE products.id IN (`
	formatParamQuery(&query, productIds)
	return db.Query(query, productIds...)
}

func QueryShelves(db DB, productIds []interface{}) (*sql.Rows, error) {
	query := `SELECT ps.shelve_id, ps.product_id, ps.is_main FROM product_shelves ps
	WHERE ps.product_id IN (`
	formatParamQuery(&query, productIds)
	return db.Query(query, productIds...)
}

func QueryShelveNames(db DB, shelvesId []interface{}) (*sql.Rows, error) {
	query := `SELECT shelves.id, shelves.shelve_name FROM shelves 
	WHERE shelves.id IN (`
	formatParamQuery(&query, shelvesId)
	return db.Query(query, shelvesId...)
}

func formatParamQuery[T interface{}](query *string, parameters []T) {
	orderLen := len(parameters)
	placeholders := make([]string, orderLen, orderLen) 

	for i := range placeholders {
		placeholders[i] = "$" + strconv.Itoa(i+1)
	}
	*query += strings.Join(placeholders, ",") + ")"
}
