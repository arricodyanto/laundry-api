package entity

type BillDetail struct {
	Id           int `json:"id"`
	Bill_Id      int `json:"billId"`
	Product_Id   int `json:"productId"`
	ProductPrice int `json:"productPrice"`
	Quantity     int `json:"qty"`
}
