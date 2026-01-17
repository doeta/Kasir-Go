package controllers

import (
	"net/http"

	"github.com/doeta/Kasir-Go/configs"
	"github.com/doeta/Kasir-Go/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)


func GetUsers(c *gin.Context) {
	var users []models.User
	configs.DB.Find(&users)
	c.JSON(http.StatusOK, users)
}


func CreateUser(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi Role
	if input.Role != "admin" && input.Role != "kasir" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role harus 'admin' atau 'kasir'"})
		return
	}

	// Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hashing password"})
		return
	}
	input.Password = string(hashedPassword)

	// Simpan ke DB
	if err := configs.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal buat user, username mungkin sudah ada"})
		return
	}

	input.Password = "" 
	c.JSON(http.StatusOK, gin.H{"message": "User berhasil dibuat", "data": input})
}


func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := configs.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi Role jika diupdate
	if input.Role != "" && input.Role != "admin" && input.Role != "kasir" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role harus 'admin' atau 'kasir'"})
		return
	}

	// Jika password diisi, hash ulang
	if input.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		input.Password = string(hashedPassword)
	} else {
		input.Password = user.Password // Keep old password
	}

	configs.DB.Model(&user).Updates(input)
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	configs.DB.Delete(&models.User{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "User berhasil dihapus"})
}
