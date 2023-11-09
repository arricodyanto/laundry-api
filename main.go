package main

import (
	"challenge-goapi/transactions"

	"github.com/gin-gonic/gin"
)

func main() {
	// Tulis kode kamu disini
	router := gin.Default()

	// Customers
	router.GET("/customers/", transactions.GetAllCustomers)
	router.GET("/customers/:id", transactions.GetCustomerById)
	router.POST("/customers", transactions.CreateNewCustomer)
	router.PUT("/customers/:id", transactions.UpdateCustomerById)
	router.DELETE("/customers/:id", transactions.DeleteCustomerById)

	// // Employees
	router.GET("/employees/", transactions.GetAllEmployees)
	router.GET("/employees/:id", transactions.GetEmployeeById)
	router.POST("/employees", transactions.CreateNewEmployee)
	router.PUT("/employees/:id", transactions.UpdateEmployeeById)
	router.DELETE("/employees/:id", transactions.DeleteEmployeeById)

	// Products
	router.GET("/products", transactions.GetAllProducts)
	router.GET("/products/:id", transactions.GetProductById)
	router.POST("/products", transactions.CreateNewProduct)
	router.PUT("/products/:id", transactions.UpdateProductById)
	router.DELETE("/products/:id", transactions.DeleteProductById)

	// Transactions
	router.GET("/transactions/:id_bill", transactions.GetTransactionById)
	router.POST("/transactions", transactions.CreateNewTransaction)

	router.Run(":8080")
}
