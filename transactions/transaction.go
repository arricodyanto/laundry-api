package transactions

import (
	"challenge-goapi/config"
	"challenge-goapi/entity"
	"challenge-goapi/utils"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

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

type billDetails struct {
	Id           string         `json:"id"`
	Bill_Id      string         `json:"billId"`
	Product      entity.Product `json:"product"`
	ProductPrice int            `json:"productPrice"`
	Quantity     int            `json:"qty"`
}
type listTransaction struct {
	Id          string          `json:"id"`
	BillDate    string          `json:"billDate"`
	EntryDate   string          `json:"entryDate"`
	FinishDate  string          `json:"finishDate"`
	Employee    entity.Employee `json:"employee"`
	Customer    entity.Customer `json:"customer"`
	BillDetails []billDetails   `json:"billDetails"`
	TotalBill   int             `json:"totalBill"`
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
		return
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

func GetTransactionById(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()
	defer utils.ErrorRecover(c)

	id := c.Param("id_bill")

	billId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Transaction ID"})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	transaction := getBillById(billId, c, tx)
	billDetails := getBillDetailsById(billId, c, tx)
	totalBill := getTotalBill(billId, c, tx)

	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	transaction.BillDate = utils.FormatTimeStringToString(transaction.BillDate)
	transaction.EntryDate = utils.FormatTimeStringToString(transaction.EntryDate)
	transaction.FinishDate = utils.FormatTimeStringToString(transaction.FinishDate)
	transaction.BillDetails = billDetails
	transaction.TotalBill = totalBill

	c.JSON(http.StatusOK, gin.H{"message": "Successfully Get Transaction", "data": transaction})
}

func getBillById(billId int, c *gin.Context, tx *sql.Tx) listTransaction {
	query := "SELECT tb.id, tb.bill_date, tb.entry_date, tb.finish_date, me.id, me.name, me.phone_number, me.address, mc.id, mc.name, mc.phone_number, mc.address FROM trx_bill tb JOIN mst_employee me ON tb.employee_id = me.id JOIN mst_customer mc ON tb.customer_id = mc.id WHERE tb.id = $1;"

	var matchedBill listTransaction
	err := tx.QueryRow(query, billId).Scan(&matchedBill.Id, &matchedBill.BillDate, &matchedBill.EntryDate, &matchedBill.FinishDate, &matchedBill.Employee.Id, &matchedBill.Employee.Name, &matchedBill.Employee.PhoneNumber, &matchedBill.Employee.Address, &matchedBill.Customer.Id, &matchedBill.Customer.Name, &matchedBill.Customer.PhoneNumber, &matchedBill.Employee.Address)
	utils.Validate(err, fmt.Sprintf("Getting Bill for ID '%d'", billId), c, tx)
	if err != nil {
		panic("Bill ID Not Found!")
	}
	return matchedBill
}

func getBillDetailsById(billId int, c *gin.Context, tx *sql.Tx) []billDetails {
	query := "SELECT tbd.id, tbd.bill_id, tp.id, tp.name, tp.price, tp.unit, tbd.product_price, tbd.qty FROM trx_bill_detail tbd JOIN mst_product tp ON tbd.product_id = tp.id WHERE tbd.bill_id = $1;"

	rows, err := tx.Query(query, billId)
	utils.Validate(err, fmt.Sprintf("Getting Bill Details for Bill ID '%d'", billId), c, tx)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var matchedBillDetails []billDetails
	for rows.Next() {
		var matchedBillDetail billDetails
		err := rows.Scan(&matchedBillDetail.Id, &matchedBillDetail.Bill_Id, &matchedBillDetail.Product.Id, &matchedBillDetail.Product.Name, &matchedBillDetail.Product.Price, &matchedBillDetail.Product.Unit, &matchedBillDetail.ProductPrice, &matchedBillDetail.Quantity)
		if err != nil {
			panic("Failed to Get Bill Details")
		}
		matchedBillDetails = append(matchedBillDetails, matchedBillDetail)
	}

	return matchedBillDetails
}

func getTotalBill(billId int, c *gin.Context, tx *sql.Tx) int {
	query := "SELECT SUM(product_price) productPrice FROM trx_bill_detail WHERE bill_id = $1;"

	var totalBill int
	err := tx.QueryRow(query, billId).Scan(&totalBill)
	utils.Validate(err, fmt.Sprintf("Getting Total Price for Bill ID '%d'", billId), c, tx)
	if err != nil {
		panic(fmt.Sprintf("Total Price for Bill ID '%d' Not Found", billId))
	}

	return totalBill
}

func GetAllTransactions(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()
	defer utils.ErrorRecover(c)

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	query := "SELECT tb.id FROM trx_bill tb JOIN mst_employee me ON tb.employee_id = me.id JOIN mst_customer mc ON tb.customer_id = mc.id WHERE 1 = 1"

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	var rows *sql.Rows
	if startDate != "" && endDate != "" {
		query += " AND tb.bill_date >= $1 AND tb.bill_date <= $2"
		rows, err = tx.Query(query, startDate, endDate)
	} else if startDate != "" && endDate == "" {
		query += " AND tb.bill_date >= $1"
		rows, err = tx.Query(query, startDate)
	} else if startDate == "" && endDate != "" {
		query += " AND tb.bill_date <= $1"
		rows, err = tx.Query(query, endDate)
	} else {
		rows, err = tx.Query(query)
	}

	utils.Validate(err, "Getting All Transactions", c, tx)
	if err != nil {
		panic("Failed to Get All Transactions")
	}
	defer rows.Close()

	var billIds []int
	for rows.Next() {
		var billId int
		err := rows.Scan(&billId)
		if err != nil {
			panic("Failed to Scan All Transactions")
		}
		billIds = append(billIds, billId)
	}

	var transactions []listTransaction
	for _, billId := range billIds {
		bill := getBillById(billId, c, tx)
		billDetails := getBillDetailsById(billId, c, tx)
		totalBill := getTotalBill(billId, c, tx)

		bill.BillDate = utils.FormatTimeStringToString(bill.BillDate)
		bill.EntryDate = utils.FormatTimeStringToString(bill.EntryDate)
		bill.FinishDate = utils.FormatTimeStringToString(bill.FinishDate)
		bill.BillDetails = billDetails
		bill.TotalBill = totalBill

		transactions = append(transactions, bill)
	}

	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully Get All Transaction", "data": transactions})
}
