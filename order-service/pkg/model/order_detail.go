package model

type OrderDetail struct {
	Id        int32 `gorm:"primary_key;column:id"`
	ProductId int32 `gorm:"column:product_id"`
	OrderId   int32 `gorm:"column:order_id"`
	Quantity  int64 `gorm:"column:quantity"`
	Price     float64
}
