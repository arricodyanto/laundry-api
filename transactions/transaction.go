package transactions

import (
	"challenge-goapi/config"
	"challenge-goapi/entity"
	"challenge-goapi/utils"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transaction struct {
	Id          string              `json:"id"`
	BillDate    string              `json:"billDate"`
	EntryDate   string              `json:"entryDate"`
	FinishDate  string              `json:"finishDate"`
	Employee_Id string              `json:"employeeId"`
	Customer_Id string              `json:"customerId"`
	BillDetails []entity.BillDetail `json:"billDetails"`
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
	defer utils.ErrorRecover(c)

	// Parse format yang didapat dari request body ke dalam db
	formattedBill := entity.Bill{
		Id:          newTransaction.Id,
		BillDate:    utils.FormatStringToTime(newTransaction.BillDate),
		EntryDate:   utils.FormatStringToTime(newTransaction.EntryDate),
		FinishDate:  utils.FormatStringToTime(newTransaction.FinishDate),
		Employee_Id: newTransaction.Employee_Id,
		Customer_Id: newTransaction.Customer_Id,
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	billId := createBill(formattedBill, c, tx)
	billDetailsId := createBillDetail(newTransaction.BillDetails, billId, c, tx)
	productPrices := updateProductPrice(billDetailsId, c, tx)

	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Update newTransaction value with returning new datas
	newTransaction.Id = utils.FormatIntToString(billId)
	for i := range newTransaction.BillDetails {
		newTransaction.BillDetails[i].Id = utils.FormatIntToString(billDetailsId[i])
		newTransaction.BillDetails[i].Bill_Id = utils.FormatIntToString(billId)
		newTransaction.BillDetails[i].ProductPrice = productPrices[i]
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Successfully Created New Transaction", "data": newTransaction})
}

func createBill(bill entity.Bill, c *gin.Context, tx *sql.Tx) int {
	var maxBillId int
	queryMaxId := "SELECT MAX(id) FROM trx_bill;"

	err := tx.QueryRow(queryMaxId).Scan(&maxBillId)
	utils.Validate(err, "Getting New ID for New Transaction", c, tx)
	if err != nil {
		panic(err.Error())
	}

	newBillId := maxBillId + 1

	queryInsert := "INSERT INTO trx_bill (id, bill_date, entry_date, finish_date, employee_id, customer_id) VALUES ($1, $2, $3, $4, $5, $6);"

	_, err = tx.Exec(queryInsert, newBillId, bill.BillDate, bill.EntryDate, bill.FinishDate, bill.Employee_Id, bill.Customer_Id)
	utils.Validate(err, fmt.Sprintf("Inserted New Bill ID '%d'", newBillId), c, tx)
	if err != nil {
		panic(err.Error())
	}
	return newBillId
}

func createBillDetail(billDetails []entity.BillDetail, billId int, c *gin.Context, tx *sql.Tx) []int {
	var newBillDetailsid []int
	for _, billDetail := range billDetails {
		var maxBillDetailId int
		queryMaxId := "SELECT MAX(id) FROM trx_bill_detail;"

		err := tx.QueryRow(queryMaxId).Scan(&maxBillDetailId)
		utils.Validate(err, "Getting New ID for New Bill Detail", c, tx)
		if err != nil {
			panic(err.Error())
		}
		newBillDetailid := maxBillDetailId + 1

		queryInsert := "INSERT INTO trx_bill_detail (id, bill_id, product_id, product_price, qty) VALUES ($1, $2, $3, $4, $5);"

		_, err = tx.Exec(queryInsert, newBillDetailid, billId, billDetail.Product_Id, billDetail.ProductPrice, billDetail.Quantity)
		utils.Validate(err, fmt.Sprintf("Inserted New Bill Details ID '%d'", newBillDetailid), c, tx)
		if err != nil {
			panic(err.Error())
		}

		newBillDetailsid = append(newBillDetailsid, newBillDetailid)
	}
	return newBillDetailsid
}

func updateProductPrice(billDetailsId []int, c *gin.Context, tx *sql.Tx) []int {
	var productPrices []int
	for _, billDetailId := range billDetailsId {
		var productPrice int

		queryProductPrice := "SELECT tbd.qty * mp.price AS product_price FROM trx_bill_detail tbd JOIN mst_product mp ON tbd.product_id = mp.id WHERE tbd.id = $1;"

		err := tx.QueryRow(queryProductPrice, billDetailId).Scan(&productPrice)
		utils.Validate(err, fmt.Sprintf("Getting Product Price Bill Detail ID '%d'", billDetailId), c, tx)
		if err != nil {
			panic(err.Error())
		}
		productPrices = append(productPrices, productPrice)

		queryUpdateProductPrice := "UPDATE trx_bill_detail SET product_price = $2 WHERE id = $1;"

		_, err = tx.Exec(queryUpdateProductPrice, billDetailId, productPrice)
		utils.Validate(err, fmt.Sprintf("Updated Product Price for Bill Detail ID '%d'", billDetailId), c, tx)
		if err != nil {
			panic(err.Error())
		}
	}
	return productPrices
}

// func GetTransactionById(c *gin.Context) {
// 	db := config.ConnectDB()
// 	defer db.Close()

// 	id := c.Param("id")

// 	billId, err := strconv.Atoi(id)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Transaction ID"})
// 		return
// 	}

// 	query := "SELECT tb.id, tb.bill_date, tb.entry_date, tb.finish_date, me."
// }
