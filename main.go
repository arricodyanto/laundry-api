package main

import (
	"challenge-goapi/transactions"

	"github.com/gin-gonic/gin"
)

func main() {
	// Tulis kode kamu disini
	router := gin.Default()

	// Customers
	router.GET("/customers/", func(c *gin.Context) {
		transactions.GetAllCustomers(c, false)
	})
	router.GET("/customers/:id", func(c *gin.Context) {
		transactions.GetCustomerById(c, false)
	})
	router.POST("/customers", func(c *gin.Context) {
		transactions.CreateNewCustomer(c, false)
	})
	router.PUT("/customers/:id", func(c *gin.Context) {
		transactions.UpdateCustomerById(c, false)
	})
	router.DELETE("/customers/:id", func(c *gin.Context) {
		transactions.DeleteCustomerById(c, false)
	})

	// Employees
	router.GET("/employees/", func(c *gin.Context) {
		transactions.GetAllCustomers(c, true)
	})
	router.GET("/employees/:id", func(c *gin.Context) {
		transactions.GetCustomerById(c, true)
	})
	router.POST("/employees", func(c *gin.Context) {
		transactions.CreateNewCustomer(c, true)
	})
	router.PUT("/employees/:id", func(c *gin.Context) {
		transactions.UpdateCustomerById(c, true)
	})
	router.DELETE("/employees/:id", func(c *gin.Context) {
		transactions.DeleteCustomerById(c, true)
	})

	// Products
	router.GET("/products", transactions.GetAllProducts)
	router.GET("/products/:id", transactions.GetProductById)
	router.POST("/products", transactions.CreateNewProduct)
	router.PUT("/products/:id", transactions.UpdateProductById)
	router.DELETE("/products/:id", transactions.DeleteProductById)

	router.Run(":8080")
}
