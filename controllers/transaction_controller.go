package controllers

import (
	"net/http"
	"time"

	"github.com/doeta/Kasir-Go/configs"
	"github.com/doeta/Kasir-Go/models"
	"github.com/gin-gonic/gin"
)

type TransactionInput struct {
	TotalAmount     int                        `json:"total_amount"`
	PaymentMethodID uint                       `json:"payment_method_id"`
	Details         []models.TransactionDetail `json:"details"`
}

func CreateTransaction(c *gin.Context) {
	var input TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buat Transaksi Baru
	// Ambil User ID dari Context (Middleware)
	userID, _ := c.Get("user_id")
	
	transaction := models.Transaction{
		TotalAmount:     input.TotalAmount,
		TransactionDate: time.Now(),
		UserID:          uint(userID.(float64)), // Casting dari JWT Claims
		PaymentMethodID: input.PaymentMethodID,
	}

	// Mulai Transaksi Database (Atomic)
	tx := configs.DB.Begin()

	// 1. Simpan Header Transaksi
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal buat transaksi"})
		return
	}

	// 2. Simpan Detail & Kurangi Stok
	for _, item := range input.Details {
		item.TransactionID = transaction.ID // Sambungkan ke header

		// Cek Stok Dulu
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Produk tidak ditemukan"})
			return
		}

		if product.Stock < item.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stok tidak cukup untuk produk: " + product.Name})
			return
		}

		// Kurangi Stok
		product.Stock = product.Stock - item.Quantity
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update stok"})
			return
		}
		
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan detail"})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Transaksi berhasil", "data": transaction})
}


func GetTransactions(c *gin.Context) {
	var transactions []models.Transaction
	configs.DB.Preload("Details").Preload("User").Preload("PaymentMethod").Find(&transactions)
	c.JSON(http.StatusOK, transactions)
}
