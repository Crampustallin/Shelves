package database

import (
	"testing"
	"database/sql"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
	T assert.TestingT
}

func (m *MockDB) Query(query string, args ...any) (*sql.Rows, error) {
	expectedQuery := `SELECT 
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
	where orders.order_number in ($1,$2,$3) order by main_shelf_name, product_id;`
    assert.Equal(m.T, expectedQuery, query)

    rows := new(sql.Rows)   
    return rows, nil
}

func TestQueryOrders(t *testing.T) {
	mockDB := new(MockDB)
	mockDB.T = t
	mockDB.On("Query", mock.Anything, mock.AnythingOfType("[]interface {}")).Return(nil, nil)

	_, err := QueryOrders(mockDB, []string{"10", "15", "20"})
	assert.NoError(t, err)
}
