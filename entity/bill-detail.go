package entity

type BillDetail struct {
	Id           string `json:"id"`
	Bill_Id      string `json:"billId"`
	Product_Id   string `json:"productId"`
	ProductPrice int    `json:"productPrice"`
	Quantity     int    `json:"qty"`
}
