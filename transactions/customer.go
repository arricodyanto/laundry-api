package transactions

import (
	"challenge-goapi/entity"
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

func GetAllCustomers(c *gin.Context) {
	query := "SELECT id, name, phone_number, address FROM mst_customer WHERE is_employee = false;"

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
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

	if len(matchedCustomers) > 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Successfully Get Customers Data", "data": matchedCustomers})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customers Not Found!"})
	}
}

func GetCustomerById(c *gin.Context) {
	id := c.Param("id")

	customerId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Customer ID"})
		return
	}

	query := "SELECT id, name, phone_number, address FROM mst_customer WHERE id = $1 AND is_employee = false;"

	var matchedCustomer resultCustomer
	err = db.QueryRow(query, customerId).Scan(&matchedCustomer.Id, &matchedCustomer.Name, &matchedCustomer.PhoneNumber, &matchedCustomer.Address)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer Not Found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully Get Customer Data", "data": matchedCustomer})
}

func CreateNewCustomer(c *gin.Context) {
	var newCustomer resultCustomer
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

	newCustomer.Id = customerId
	c.JSON(http.StatusCreated, gin.H{"message": "Successfully Created New Customer", "data": newCustomer})
}

func UpdateCustomerById(c *gin.Context) {
	id := c.Param("id")

	customerId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Customer ID"})
		return
	}

	var isEmployee bool
	queryCheckEmployee := "SELECT is_employee FROM mst_customer WHERE id = $1;"

	err = db.QueryRow(queryCheckEmployee, customerId).Scan(&isEmployee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Check Employee Status"})
		return
	}

	if isEmployee {
		c.JSON(http.StatusForbidden, gin.H{"error": "Failed to Updated This Employee, Forbidden!"})
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
	c.JSON(http.StatusOK, gin.H{"message": "Successfully Updated Customer '" + id + "'", "data": updatedCustomer})
}

func DeleteCustomerById(c *gin.Context) {
	id := c.Param("id")

	customerId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Customer ID"})
		return
	}

	var isEmployee bool
	queryCheckEmployee := "SELECT is_employee FROM mst_customer WHERE id = $1;"

	err = db.QueryRow(queryCheckEmployee, customerId).Scan(&isEmployee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Check Employee Status"})
		return
	}

	if isEmployee {
		c.JSON(http.StatusForbidden, gin.H{"error": "Failed to Deleted This Employee, Forbidden!"})
		return
	}

	query := "DELETE FROM mst_customer WHERE id = $1 RETURNING *;"

	row, _ := db.Exec(query, customerId)
	if rowAffected, _ := row.RowsAffected(); rowAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully Deleted Customer '" + id + "'", "data": "OK"})
}

func GetAllEmployees(c *gin.Context) {
	query := "SELECT id, name, phone_number, address, is_employee FROM mst_customer WHERE is_employee = true;"

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer rows.Close()

	var matchedEmployee []entity.Customer
	for rows.Next() {
		var employee entity.Customer
		err := rows.Scan(&employee.Id, &employee.Name, &employee.PhoneNumber, &employee.Address, &employee.IsEmployee)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		matchedEmployee = append(matchedEmployee, employee)
	}

	if len(matchedEmployee) > 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Successfully Get Employees Data", "data": matchedEmployee})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employees Not Found!"})
	}
}

func GetEmployeeById(c *gin.Context) {
	id := c.Param("id")

	employeeId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Employee ID"})
		return
	}

	query := "SELECT id, name, phone_number, address, is_employee FROM mst_customer WHERE id = $1 AND is_employee = true;"

	var matchedEmployee entity.Customer
	err = db.QueryRow(query, employeeId).Scan(&matchedEmployee.Id, &matchedEmployee.Name, &matchedEmployee.PhoneNumber, &matchedEmployee.Address, &matchedEmployee.IsEmployee)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee Not Found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully Get Employee Data", "data": matchedEmployee})
}

func CreateNewEmployee(c *gin.Context) {
	var newEmployee entity.Customer
	err := c.ShouldBind(&newEmployee)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var maxId int
	queryMaxId := "SELECT MAX(id) FROM mst_customer;"
	err = db.QueryRow(queryMaxId).Scan(&maxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed Getting New ID for New Employee"})
		return
	}

	employeeId := maxId + 1
	isEmployee := true

	queryInsert := "INSERT INTO mst_customer (id, name, phone_number, address, is_employee) VALUES ($1, $2, $3, $4, $5);"
	_, err = db.Exec(queryInsert, employeeId, newEmployee.Name, newEmployee.PhoneNumber, newEmployee.Address, isEmployee)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newEmployee.Id = employeeId
	newEmployee.IsEmployee = isEmployee
	c.JSON(http.StatusCreated, gin.H{"message": "Successfully Created New Employee", "data": newEmployee})
}

func UpdateEmployeeById(c *gin.Context) {
	id := c.Param("id")

	employeeId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Employee ID"})
		return
	}

	var isEmployee bool
	queryCheckEmployee := "SELECT is_employee FROM mst_customer WHERE id = $1;"

	err = db.QueryRow(queryCheckEmployee, employeeId).Scan(&isEmployee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Check Employee Status"})
		return
	}

	if !isEmployee {
		c.JSON(http.StatusForbidden, gin.H{"error": "Failed to Updated This Customer, Forbidden!"})
		return
	}

	var updatedEmployee entity.Customer
	err = c.ShouldBind(&updatedEmployee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query := "UPDATE mst_customer SET name = $2, phone_number = $3, address = $4 WHERE id = $1;"

	_, err = db.Exec(query, employeeId, updatedEmployee.Name, updatedEmployee.PhoneNumber, updatedEmployee.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedEmployee.Id = employeeId
	updatedEmployee.IsEmployee = true
	c.JSON(http.StatusOK, gin.H{"message": "Successfully Updated Employee '" + id + "'", "data": updatedEmployee})
}

func DeleteEmployeeById(c *gin.Context) {
	id := c.Param("id")

	employeeId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Employee ID"})
		return
	}

	var isEmployee bool
	queryCheckEmployee := "SELECT is_employee FROM mst_customer WHERE id = $1;"

	err = db.QueryRow(queryCheckEmployee, employeeId).Scan(&isEmployee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Check Employee Status"})
		return
	}

	if !isEmployee {
		c.JSON(http.StatusForbidden, gin.H{"error": "Failed to Deleted This Customer, Forbidden!"})
		return
	}

	query := "DELETE FROM mst_customer WHERE id = $1 RETURNING *;"

	row, _ := db.Exec(query, employeeId)
	if rowAffected, _ := row.RowsAffected(); rowAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully Deleted Employee '" + id + "'", "data": "OK"})
}
