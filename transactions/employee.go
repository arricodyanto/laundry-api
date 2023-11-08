package transactions

import (
	"challenge-goapi/config"
	"challenge-goapi/entity"
	"challenge-goapi/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllEmployees(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	query := "SELECT id, name, phone_number, address FROM mst_employee;"

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var matchedEmployees []entity.Employee
	for rows.Next() {
		var employee entity.Employee
		err := rows.Scan(&employee.Id, &employee.Name, &employee.PhoneNumber, &employee.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		matchedEmployees = append(matchedEmployees, employee)
	}
	if len(matchedEmployees) > 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Successfully Get Employees Data", "data": matchedEmployees})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employees Not Found!"})
	}
}

func GetEmployeeById(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	id := c.Param("id")

	_, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Employee ID"})
		return
	}

	query := "SELECT id, name, phone_number, address FROM mst_employee WHERE id = $1;"

	var matchedEmployee entity.Employee
	err = db.QueryRow(query, id).Scan(&matchedEmployee.Id, &matchedEmployee.Name, &matchedEmployee.PhoneNumber, &matchedEmployee.Address)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee Not Found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully Get Employee Data", "data": matchedEmployee})
}

func CreateNewEmployee(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	var newEmployee entity.Employee
	err := c.ShouldBind(&newEmployee)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var maxId int
	queryMaxId := "SELECT MAX(id) FROM mst_employee;"
	err = db.QueryRow(queryMaxId).Scan(&maxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed Getting New ID for New Employee"})
		return
	}

	employeeId := maxId + 1

	queryInsert := "INSERT INTO mst_employee (id, name, phone_number, address) VALUES ($1, $2, $3, $4);"
	_, err = db.Exec(queryInsert, employeeId, newEmployee.Name, newEmployee.PhoneNumber, newEmployee.Address)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newEmployee.Id = utils.FormatIntToString(employeeId)
	c.JSON(http.StatusCreated, gin.H{"message": "Successfully Created New Employee", "data": newEmployee})
}

func UpdateEmployeeById(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	id := c.Param("id")

	_, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Employee ID"})
		return
	}

	var updatedEmployee entity.Employee
	err = c.ShouldBind(&updatedEmployee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query := "UPDATE mst_employee SET name = $2, phone_number = $3, address = $4 WHERE id = $1 RETURNING *;"

	row, err := db.Exec(query, id, updatedEmployee.Name, updatedEmployee.PhoneNumber, updatedEmployee.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rowAffected, _ := row.RowsAffected(); rowAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee ID Not Found!"})
		return
	}

	updatedEmployee.Id = id
	c.JSON(http.StatusOK, gin.H{"message": "Successfully Updated Employee '" + id + "'", "data": updatedEmployee})
}

func DeleteEmployeeById(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	id := c.Param("id")

	employeeId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Employee ID"})
		return
	}

	query := "DELETE FROM mst_employee WHERE id = $1 RETURNING *;"

	row, err := db.Exec(query, employeeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rowAffected, _ := row.RowsAffected(); rowAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee ID Not Found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully Deleted Employee '" + id + "'", "data": "OK"})
}
