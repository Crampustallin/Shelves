package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/Crampustallin/Shelves/internal/database"
	"github.com/Crampustallin/Shelves/internal/models"
)

func HandleOrdersQuery(orderNumbers []string) {
	db, err := database.NewDB()
	if err != nil {
		log.Fatal("error while connecting to database", err)
		return
	}
	defer db.Close()

	rows, err := database.QueryOrders(db, orderNumbers)
	if err != nil {
		log.Fatal("erro while executing query", err)
		return
	}
	defer rows.Close()

	shelvesMap := make(map[string][]models.Shelve)

	for rows.Next() {
		var shelve models.Shelve
		if err := rows.Scan(&shelve.OrderNumber, &shelve.ProductName, &shelve.ProductId, &shelve.Quantity,
		&shelve.MainShelf, &shelve.SecondaryShelf); err != nil {
			log.Fatal(err)
		}
		shelvesMap[shelve.MainShelf] = append(shelvesMap[shelve.MainShelf], shelve)
	}

	fmt.Printf("=+=+=+=\nСтраница сборки заказов %v\n", strings.Join(orderNumbers, ","))

	for key := range shelvesMap {
		fmt.Printf("\n===Стеллаж %s", key)
		for _, product := range shelvesMap[key] {
			message := fmt.Sprintf("%v (id=%v)\nзаказ %v, %v шт\n", product.ProductName, product.ProductId, product.OrderNumber, product.Quantity)
			if product.SecondaryShelf != "" {
				message += "доп стеллаж: " + product.SecondaryShelf 
			}
			fmt.Printf(message)
		}
	}
}
