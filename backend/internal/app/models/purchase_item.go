package models

type PurchaseItem struct {
	GormModel
	PurchaseID uint    `json:"purchaseID"` // Foreign key to Purchase
	ProductID  uint    `json:"productID"`  // Foreign key to Product
	Product    Product `json:"product"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	TotalPrice float64 `json:"totalPrice"`
}
