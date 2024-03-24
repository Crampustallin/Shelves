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

	ids := ParseToInterface(orderNumbers)
	rows, err := database.QueryOrders(db, ids)
	if err != nil {
		log.Fatal("error while executing query orders", err)
		return
	}

	ordersMap := make(map[int]string)
	ids = make([]interface{}, len(orderNumbers), len(orderNumbers))

	for rows.Next() {
		var orderId int
		var orderNumber string
		if err := rows.Scan(&orderId, &orderNumber); err != nil {
			log.Fatal(err)
			return
		}
		ordersMap[orderId] = orderNumber
		ids = append(ids, orderId)
	}
	rows.Close()

	rows, err = database.QuerySummaries(db, ids)
	if err != nil {
		log.Fatal("error while executing query products", err)
		return
	}

	orderInfo := make(map[int]map[int]models.Order)

	for rows.Next() {
		var orderId int
		var productId int
		var quantity int
		if err := rows.Scan(&orderId, &productId, &quantity); err != nil {
			log.Fatal(err)
			return
		}
		if orderInfo[productId] == nil {
			orderInfo[productId] = make(map[int]models.Order)
		}
		orderInfo[productId][orderId] = models.Order{Quantity: quantity}
	}
	rows.Close()

	ids = make([]interface{}, len(orderInfo), len(orderInfo))
	count := 0
	for key := range orderInfo {
		ids[count] = key
		count++
	}

	rows, err = database.QueryProducts(db,ids)
	if err != nil {
		log.Fatal("error while executing query products")
	}

	productsMap := make(map[int]models.Product)

	for rows.Next() {
		var productId int
		var productName string
		if err := rows.Scan(&productId, &productName); err != nil {
			log.Fatal(err)
		}
		productsMap[productId] = models.Product{ProductName: productName} 
	}
	rows.Close()

	ids = make([]interface{}, len(productsMap), len(productsMap))
	count = 0
	for key := range productsMap {
		ids[count] = key
		count++
	}

	rows, err = database.QueryShelves(db, ids)
	if err != nil {
		log.Fatal("error while executing shelves query")
	}

	mainShelvesMap := make(map[int]models.Shelve)
	shelvesMap := make(map[int]string)
	for rows.Next() {
		var shelveId int
		var productId int
		var isMain bool
		if err := rows.Scan(&shelveId, &productId, &isMain); err != nil {
			log.Fatal(err)
		}
		if isMain {
			if val, ok := mainShelvesMap[shelveId]; ok {
				val.ProductIds = append(val.ProductIds, productId)
				mainShelvesMap[shelveId] = val
			} else {
				mainShelvesMap[shelveId] = models.Shelve{ProductIds: make([]int, 0, 1)}
				val.ProductIds = append(val.ProductIds, productId)
				mainShelvesMap[shelveId] = val
			}
		} else if val, ok := productsMap[productId]; ok {
			val.SecondaryShelveIds = append(val.SecondaryShelveIds, shelveId)
			productsMap[productId] = val
		}
		shelvesMap[shelveId] = ""
	}
	rows.Close()

	ids = make([]interface{}, len(shelvesMap), len(shelvesMap))
	count = 0
	for key := range shelvesMap {
		ids[count] = key
		count++
	}

	rows, err = database.QueryShelveNames(db, ids)
	if err != nil {
		log.Fatal("error while querying shevle names")
	}
	
	for rows.Next() {
		var shelveId int
		var shelveName string
		if err := rows.Scan(&shelveId, &shelveName); err != nil {
			log.Fatal(err)
		}
		shelvesMap[shelveId] = shelveName
	}

	PrintShelves(mainShelvesMap, shelvesMap, orderNumbers, productsMap, orderInfo, ordersMap)
}

func PrintShelves(mainShelvesMap map[int]models.Shelve, shelvesMap map[int]string, orderNumbers []string, productsMap map[int]models.Product,
orderInfo map[int]map[int]models.Order, ordersMap map[int]string) {
	fmt.Printf("=+=+=+=\nСтраница сборки заказов %v\n\n", strings.Join(orderNumbers, ","))
	for key, val := range mainShelvesMap {
		fmt.Printf("===Стеллаж %s\n", shelvesMap[key])
		for _, productId := range val.ProductIds {
			for orderId, orderVal := range orderInfo[productId] { 
				fmt.Printf("%s (id=%v)\n", productsMap[productId].ProductName, productId)
				fmt.Printf("заказ %s, %v шт\n", ordersMap[orderId], orderVal.Quantity)
				s := ""
				for _, secondaryId := range productsMap[productId].SecondaryShelveIds {
					s += shelvesMap[secondaryId] + ","
				}
				if s != "" {
					fmt.Printf("доп стеллаж: %s\n", s[:len(s)-1])
				}
				fmt.Print("\n")
			}
		}
	}
}

func ParseToInterface[T any](src []T) []interface{} {
	result := make([]interface{},len(src), len(src))
	for i := range src {
		result[i] = src[i]
	}
	return result
}
