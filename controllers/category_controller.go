package controllers

import (
	"net/http"

	"github.com/doeta/Kasir-Go/configs"
	"github.com/doeta/Kasir-Go/models"
	"github.com/gin-gonic/gin"
)

// GetCategories - Semua user logedin bisa lihat kategori
func GetCategories(c *gin.Context) {
	var categories []models.Category
	configs.DB.Find(&categories)
	c.JSON(http.StatusOK, categories)
}

type CategoryInput struct {
	Name string `json:"name" binding:"required"`
}

// CreateCategory - Hanya admin
func CreateCategory(c *gin.Context) {
	var input CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{Name: input.Name}
	if err := configs.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// UpdateCategory - Hanya admin
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := configs.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var input CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.Name = input.Name
	configs.DB.Save(&category)
	c.JSON(http.StatusOK, category)
}

// DeleteCategory - Hanya admin
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := configs.DB.Delete(&models.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}
