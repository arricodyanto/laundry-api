package transactions

import (
	"challenge-goapi/config"
	"challenge-goapi/entity"
	"challenge-goapi/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllCustomers(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	query := "SELECT id, name, phone_number, address FROM mst_customer;"

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var matchedCustomers []entity.Customer
	for rows.Next() {
		var customer entity.Customer
		err := rows.Scan(&customer.Id, &customer.Name, &customer.PhoneNumber, &customer.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		matchedCustomers = append(matchedCustomers, customer)
	}
	if len(matchedCustomers) > 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Successfully Get Customers Data", "data": matchedCustomers})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customers Not Found!"})
	}
}

func GetCustomerById(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	id := c.Param("id")

	_, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Customer ID"})
		return
	}

	query := "SELECT id, name, phone_number, address FROM mst_customer WHERE id = $1;"

	var matchedCustomer entity.Customer
	err = db.QueryRow(query, id).Scan(&matchedCustomer.Id, &matchedCustomer.Name, &matchedCustomer.PhoneNumber, &matchedCustomer.Address)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer Not Found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully Get Customer Data", "data": matchedCustomer})
}

func CreateNewCustomer(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	var newCustomer entity.Customer
	err := c.ShouldBind(&newCustomer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var maxId int
	queryMaxId := "SELECT MAX(id) FROM mst_customer;"
	err = db.QueryRow(queryMaxId).Scan(&maxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed Getting New ID for New Customer"})
		return
	}

	customerId := maxId + 1

	queryInsert := "INSERT INTO mst_customer (id, name, phone_number, address) VALUES ($1, $2, $3, $4);"
	_, err = db.Exec(queryInsert, customerId, newCustomer.Name, newCustomer.PhoneNumber, newCustomer.Address)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newCustomer.Id = utils.FormatIntToString(customerId)
	c.JSON(http.StatusCreated, gin.H{"message": "Successfully Created New Customer", "data": newCustomer})
}

func UpdateCustomerById(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	id := c.Param("id")

	_, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Customer ID"})
		return
	}

	var updatedCustomer entity.Customer
	err = c.ShouldBind(&updatedCustomer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query := "UPDATE mst_customer SET name = $2, phone_number = $3, address = $4 WHERE id = $1 RETURNING *;"

	row, err := db.Exec(query, id, updatedCustomer.Name, updatedCustomer.PhoneNumber, updatedCustomer.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rowAffected, _ := row.RowsAffected(); rowAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer ID Not Found!"})
		return
	}

	updatedCustomer.Id = id
	c.JSON(http.StatusOK, gin.H{"message": "Successfully Updated Customer '" + id + "'", "data": updatedCustomer})
}

func DeleteCustomerById(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	id := c.Param("id")

	customerId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Customer ID"})
		return
	}

	query := "DELETE FROM mst_customer WHERE id = $1 RETURNING *;"

	row, err := db.Exec(query, customerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rowAffected, _ := row.RowsAffected(); rowAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer ID Not Found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully Deleted Customer '" + id + "'", "data": "OK"})
}
