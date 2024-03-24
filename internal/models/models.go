package models


type Order struct {
	Quantity int 
	ProductId int 
}

type Shelve struct {
	ShelveName string
	isMain bool
	ProductIds []int
}

type Product struct {
	ProductName string
	SecondaryShelveIds []int
}	
