package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestQueryOrders(t *testing.T) {
	targetQuery := `SELECT * FROM test WHERE id IN (`
	formatParamQuery(&targetQuery, []int{1,2,3})
	expected := "SELECT * FROM test WHERE id IN ($1,$2,$3)"
	assert.Equal(t, expected, targetQuery)
}
