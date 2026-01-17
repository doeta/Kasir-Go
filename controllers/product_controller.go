package controllers

import (
	"net/http"

	"github.com/doeta/Kasir-Go/configs"
	"github.com/doeta/Kasir-Go/models"
	"github.com/gin-gonic/gin"
)


func CreateProduct(c *gin.Context) {
	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	configs.DB.Create(&input)
	c.JSON(http.StatusOK, input)
}


func GetProducts(c *gin.Context) {
	var products []models.Product
	configs.DB.Find(&products)
	c.JSON(http.StatusOK, products)
}


func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := configs.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	configs.DB.Model(&product).Updates(input)
	c.JSON(http.StatusOK, product)
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	configs.DB.Delete(&models.Product{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil dihapus"})
}
