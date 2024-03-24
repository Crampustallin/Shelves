package models


type Order struct {
	Quantity int 
	ProductId int 
}

type Shelve struct {
	ShelveName string
	ProductIds []int
}

type Product struct {
	ProductName string
	SecondaryShelveIds []int
}	
