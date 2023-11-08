package transactions

import (
	"challenge-goapi/config"
	"challenge-goapi/entity"
	"challenge-goapi/utils"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionDetails struct {
	Id           string `json:"id"`
	Bill_Id      string `json:"billId"`
	Product_Id   string `json:"productId"`
	ProductPrice int    `json:"productPrice"`
	Quantity     int    `json:"qty"`
}

type transaction struct {
	Id          int                  `json:"id"`
	BillDate    string               `json:"billDate"`
	EntryDate   string               `json:"entryDate"`
	FinishDate  string               `json:"finishDate"`
	Employee_Id string               `json:"employeeId"`
	Customer_Id string               `json:"customerId"`
	BillDetails []transactionDetails `json:"billDetails"`
}

func CreateNewTransaction(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	var newTransaction transaction
	err := c.ShouldBind(&newTransaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Not Binding": err.Error()})
		return
	}

	// Parse format yang didapat dari request body ke dalam db
	formattedBill := entity.Bill{
		Id:          newTransaction.Id,
		BillDate:    utils.FormatStringToTime(newTransaction.BillDate),
		EntryDate:   utils.FormatStringToTime(newTransaction.EntryDate),
		FinishDate:  utils.FormatStringToTime(newTransaction.FinishDate),
		Employee_Id: utils.FormatStringToInt(newTransaction.Employee_Id),
		Customer_Id: utils.FormatStringToInt(newTransaction.Customer_Id),
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	billId := createBill(formattedBill, c, tx)
	var formattedBillDetails []entity.BillDetail
	for _, billDetail := range newTransaction.BillDetails {
		formattedBillDetail := entity.BillDetail{
			Id:           utils.FormatStringToInt(billDetail.Id),
			Bill_Id:      billId,
			Product_Id:   utils.FormatStringToInt(billDetail.Product_Id),
			ProductPrice: billDetail.ProductPrice,
			Quantity:     billDetail.Quantity,
		}
		formattedBillDetails = append(formattedBillDetails, formattedBillDetail)
	}

	billDetailsId := createBillDetail(formattedBillDetails, billId, c, tx)

	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Commited": err.Error()})
		return
	}
	newTransaction.Id = billId
	for i, _ := range newTransaction.BillDetails {
		newTransaction.BillDetails[i].Id = utils.FormatIntToString(billDetailsId[i])
		newTransaction.BillDetails[i].Bill_Id = utils.FormatIntToString(billId)
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Successfully Created New Transaction", "data": newTransaction})
}

func createBill(bill entity.Bill, c *gin.Context, tx *sql.Tx) int {
	var maxBillId int
	queryMaxId := "SELECT MAX(id) FROM trx_bill;"

	err := tx.QueryRow(queryMaxId).Scan(&maxBillId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Created Bill": err.Error()})
		return 0
	}
	newBillId := maxBillId + 1

	queryInsert := "INSERT INTO trx_bill (id, bill_date, entry_date, finish_date, employee_id, customer_id) VALUES ($1, $2, $3, $4, $5, $6);"

	_, err = tx.Exec(queryInsert, newBillId, bill.BillDate, bill.EntryDate, bill.FinishDate, bill.Employee_Id, bill.Customer_Id)
	utils.Validate(err, "Failed to Inserted New Bill", c, tx)
	return newBillId
}

func createBillDetail(billDetails []entity.BillDetail, billId int, c *gin.Context, tx *sql.Tx) []int {
	var newBillDetailsid []int
	for _, billDetail := range billDetails {
		var maxBillDetailId int
		queryMaxId := "SELECT MAX(id) FROM trx_bill_detail;"

		err := tx.QueryRow(queryMaxId).Scan(&maxBillDetailId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error Created Bill Detail": err.Error()})
			return []int{}
		}
		newBillDetailid := maxBillDetailId + 1

		queryInsert := "INSERT INTO trx_bill_detail (id, bill_id, product_id, product_price, qty) VALUES ($1, $2, $3, $4, $5);"

		_, err = tx.Exec(queryInsert, newBillDetailid, billId, billDetail.Product_Id, billDetail.ProductPrice, billDetail.Quantity)
		utils.Validate(err, "Failed to Inserted New Bill Details", c, tx)

		newBillDetailsid = append(newBillDetailsid, newBillDetailid)
	}
	return newBillDetailsid
}
