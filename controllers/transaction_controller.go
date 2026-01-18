package controllers

import (
	"net/http"
	"time"

	"github.com/doeta/Kasir-Go/configs"
	"github.com/doeta/Kasir-Go/models"
	"github.com/gin-gonic/gin"
)

type TransactionDetailInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

type CreateTransactionInput struct {
	PaymentMethodID uint                     `json:"payment_method_id" binding:"required"`
	Details         []TransactionDetailInput `json:"details" binding:"required,dive"`
}

func CreateTransaction(c *gin.Context) {
	var input CreateTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil User ID dari Context (Middleware)
	userID, _ := c.Get("user_id")

	// Mulai Transaksi Database (Atomic)
	tx := configs.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var totalAmount int
	var transactionDetails []models.TransactionDetail

	// 1. Validasi Stok & Hitung Total (Looping Inputs)
	for _, itemInput := range input.Details {
		var product models.Product
		// Lock row with clause to prevent race condition if high concurrency (Optional but good)
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&product, itemInput.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Produk tidak ditemukan (ID: " + string(rune(itemInput.ProductID)) + ")"})
			return
		}

		if product.Stock < itemInput.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stok tidak cukup untuk produk: " + product.Name})
			return
		}

		// Calculate Subtotal
		subtotal := product.Price * itemInput.Quantity
		totalAmount += subtotal

		// Kurangi Stok
		product.Stock -= itemInput.Quantity
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update stok"})
			return
		}

		// Prepare Detail Data
		transactionDetails = append(transactionDetails, models.TransactionDetail{
			ProductID: itemInput.ProductID,
			Quantity:  itemInput.Quantity,
			Subtotal:  subtotal,
		})
	}

	// 2. Buat Header Transaksi
	transaction := models.Transaction{
		UserID:          uint(userID.(float64)),
		PaymentMethodID: input.PaymentMethodID,
		TransactionDate: time.Now(),
		TotalAmount:     totalAmount,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal buat transaksi"})
		return
	}

	// 3. Simpan Details dengan TransactionID yang baru dibuat
	for i := range transactionDetails {
		transactionDetails[i].TransactionID = transaction.ID
		if err := tx.Create(&transactionDetails[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan detail transaksi"})
			return
		}
	}

	// Commit Transaksi
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal commit transaksi"})
		return
	}

	// Load Data Lengkap untuk Response
	// (Optional: Bisa skip jika ingin response cepat, tapi user biasanya butuh struk)
	configs.DB.Preload("Details").First(&transaction, transaction.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaksi berhasil",
		"data":    transaction,
	})
}


func GetTransactions(c *gin.Context) {
	var transactions []models.Transaction
	configs.DB.Preload("Details").Preload("User").Preload("PaymentMethod").Find(&transactions)
	c.JSON(http.StatusOK, transactions)
}
