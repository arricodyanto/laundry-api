package main

import (
	"challenge-goapi/transactions"

	"github.com/gin-gonic/gin"
)

func main() {
	// Tulis kode kamu disini
	router := gin.Default()

	// Products
	router.GET("/products", transactions.GetAllProducts)
	router.GET("/products/:id", transactions.GetProductById)
	router.POST("/products", transactions.CreateNewProduct)
	router.PUT("/products/:id", transactions.UpdateProductById)
	router.DELETE("/products/:id", transactions.DeleteProductById)

	router.Run(":8080")
}
