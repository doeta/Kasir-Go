package controllers

import (
	"net/http"

	"github.com/doeta/Kasir-Go/configs"
	"github.com/doeta/Kasir-Go/models"
	"github.com/gin-gonic/gin"
)


func CreatePayment(c *gin.Context) {
	var payment models.PaymentMethod
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	configs.DB.Create(&payment)
	c.JSON(http.StatusOK, payment)
}

func GetPayments(c *gin.Context) {
	var payments []models.PaymentMethod
	configs.DB.Find(&payments)
	c.JSON(http.StatusOK, payments)
}


func UpdatePayment(c *gin.Context) {
	id := c.Param("id")
	var payment models.PaymentMethod
	if err := configs.DB.First(&payment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}
	
	var input models.PaymentMethod
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	configs.DB.Model(&payment).Updates(input)
	c.JSON(http.StatusOK, payment)
}

func DeletePayment(c *gin.Context) {
	id := c.Param("id")
	configs.DB.Delete(&models.PaymentMethod{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Data dihapus"})
}
