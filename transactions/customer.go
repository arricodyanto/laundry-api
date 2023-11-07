package transactions

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type resultCustomer struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
}

var message string

func GetAllCustomers(c *gin.Context, isEmployee bool) {
	query := "SELECT id, name, phone_number, address FROM mst_customer WHERE is_employee = $1;"

	rows, err := db.Query(query, isEmployee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var matchedCustomers []resultCustomer
	for rows.Next() {
		var customer resultCustomer
		err := rows.Scan(&customer.Id, &customer.Name, &customer.PhoneNumber, &customer.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		matchedCustomers = append(matchedCustomers, customer)
	}

	if !isEmployee {
		message = "Customers"
	} else {
		message = "Employees"
	}
	if len(matchedCustomers) > 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Successfully Get " + message + " Data", "data": matchedCustomers})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": message + " Not Found!"})
	}
}

func GetCustomerById(c *gin.Context, isEmployee bool) {
	id := c.Param("id")

	if !isEmployee {
		message = "Customer"
	} else {
		message = "Employee"
	}

	customerId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid" + message + " ID"})
		return
	}

	strIsEmployee := strconv.FormatBool(isEmployee)

	query := "SELECT id, name, phone_number, address FROM mst_customer WHERE id = $1 AND is_employee = " + strIsEmployee + ";"

	var matchedCustomer resultCustomer
	err = db.QueryRow(query, customerId).Scan(&matchedCustomer.Id, &matchedCustomer.Name, &matchedCustomer.PhoneNumber, &matchedCustomer.Address)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": message + " Not Found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully Get " + message + " Data", "data": matchedCustomer})
}

func CreateNewCustomer(c *gin.Context, isEmployee bool) {
	var newCustomer resultCustomer
	err := c.ShouldBind(&newCustomer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isEmployee {
		message = "Customer"
	} else {
		message = "Employee"
	}

	var maxId int
	queryMaxId := "SELECT MAX(id) FROM mst_customer;"
	err = db.QueryRow(queryMaxId).Scan(&maxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed Getting New ID for New " + message})
		return
	}

	customerId := maxId + 1

	queryInsert := "INSERT INTO mst_customer (id, name, phone_number, address, is_employee) VALUES ($1, $2, $3, $4, $5);"
	_, err = db.Exec(queryInsert, customerId, newCustomer.Name, newCustomer.PhoneNumber, newCustomer.Address, isEmployee)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newCustomer.Id = customerId
	c.JSON(http.StatusCreated, gin.H{"message": "Successfully Created New " + message, "data": newCustomer})
}

func UpdateCustomerById(c *gin.Context, isEmployee bool) {
	id := c.Param("id")

	if !isEmployee {
		message = "Customer"
	} else {
		message = "Employee"
	}

	customerId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid " + message + " ID"})
		return
	}

	queryCheckEmployee := "SELECT is_employee FROM mst_customer WHERE id = $1;"

	var isCustomer bool
	err = db.QueryRow(queryCheckEmployee, customerId).Scan(&isCustomer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Check " + message + " Status. ID Not Found"})
		return
	}

	if !isEmployee && isCustomer {
		c.JSON(http.StatusForbidden, gin.H{"error": "Failed to Updated This Employee, Forbidden!"})
		return
	} else if isEmployee && !isCustomer {
		c.JSON(http.StatusForbidden, gin.H{"error": "Failed to Updated This Customer, Forbidden!"})
		return
	}

	var updatedCustomer resultCustomer
	err = c.ShouldBind(&updatedCustomer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query := "UPDATE mst_customer SET name = $2, phone_number = $3, address = $4 WHERE id = $1;"

	_, err = db.Exec(query, customerId, updatedCustomer.Name, updatedCustomer.PhoneNumber, updatedCustomer.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedCustomer.Id = customerId
	c.JSON(http.StatusOK, gin.H{"message": "Successfully Updated " + message + " '" + id + "'", "data": updatedCustomer})
}

func DeleteCustomerById(c *gin.Context, isEmployee bool) {
	id := c.Param("id")

	if !isEmployee {
		message = "Customer"
	} else {
		message = "Employee"
	}

	customerId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid " + message + " ID"})
		return
	}

	var isCustomer bool
	queryCheckEmployee := "SELECT is_employee FROM mst_customer WHERE id = $1;"

	err = db.QueryRow(queryCheckEmployee, customerId).Scan(&isCustomer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Check " + message + " Status. ID Not Found"})
		return
	}

	if !isEmployee && isCustomer {
		c.JSON(http.StatusForbidden, gin.H{"error": "Failed to Deleted This Employee, Forbidden!"})
		return
	} else if isEmployee && !isCustomer {
		c.JSON(http.StatusForbidden, gin.H{"error": "Failed to Deleted This Customer, Forbidden!"})
		return
	}

	query := "DELETE FROM mst_customer WHERE id = $1 RETURNING *;"

	row, _ := db.Exec(query, customerId)
	if rowAffected, _ := row.RowsAffected(); rowAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully Deleted " + message + " '" + id + "'", "data": "OK"})
}
