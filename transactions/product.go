package transactions

import (
	"challenge-goapi/config"
	"challenge-goapi/entity"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllProducts(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	productName := c.Query("name")

	query := "SELECT id, name, price, unit FROM mst_product"

	var rows *sql.Rows
	var err error

	if productName != "" {
		query += " WHERE name ILIKE '%' || $1 || '%';"
		rows, err = db.Query(query, productName)
	} else {
		rows, err = db.Query(query)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer rows.Close()

	var matchedProducts []entity.Product
	for rows.Next() {
		var product entity.Product
		err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Unit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		matchedProducts = append(matchedProducts, product)
	}

	if len(matchedProducts) > 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Successfully Get Books Data", "data": matchedProducts})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Products Not Found!"})
	}
}

func GetProductById(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	id := c.Param("id")

	productId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID"})
		return
	}

	query := "SELECT id, name, price, unit FROM mst_product WHERE id = $1;"

	var matchedProduct entity.Product
	err = db.QueryRow(query, productId).Scan(&matchedProduct.Id, &matchedProduct.Name, &matchedProduct.Price, &matchedProduct.Unit)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product Founded!", "data": matchedProduct})
}

func CreateNewProduct(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	var newProduct entity.Product
	err := c.ShouldBind(&newProduct)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var maxId int
	queryMaxId := "SELECT MAX(id) FROM mst_product;"
	err = db.QueryRow(queryMaxId).Scan(&maxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed Getting New ID for New Product"})
		return
	}

	productId := maxId + 1
	queryInsert := "INSERT INTO mst_product (id, name, price, unit) VALUES ($1, $2, $3, $4);"
	_, err = db.Exec(queryInsert, productId, newProduct.Name, newProduct.Price, newProduct.Unit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	newProduct.Id = productId
	c.JSON(http.StatusCreated, gin.H{"message": "Successfully Created New Product", "data": newProduct})
}

func UpdateProductById(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	id := c.Param("id")

	productId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID"})
		return
	}

	var updatedProduct entity.Product
	err = c.ShouldBind(&updatedProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	query := "UPDATE mst_product SET name = $2, price = $3, unit = $4 WHERE id = $1;"

	_, err = db.Exec(query, productId, updatedProduct.Name, updatedProduct.Price, updatedProduct.Unit)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product Not Found"})
		return
	}
	updatedProduct.Id = productId
	c.JSON(http.StatusOK, gin.H{"message": "Successfully Updated Product ID '" + id + "'", "data": updatedProduct})
}

func DeleteProductById(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	id := c.Param("id")
	productId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID"})
		return
	}

	query := "DELETE FROM mst_product WHERE id = $1 RETURNING *;"

	row, _ := db.Exec(query, productId)
	if rowAffected, _ := row.RowsAffected(); rowAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product Not Found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully Deleted Product", "data": "OK"})
}
