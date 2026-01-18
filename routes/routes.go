package routes

import (
	"net/http"

	"github.com/doeta/Kasir-Go/controllers"
	"github.com/doeta/Kasir-Go/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Root Endpoint (Health Check)
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Kasir-Go API is running",
			"docs":    "/swagger/index.html",
		})
	})

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API Routes
	api := r.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Public Protected (Bisa diakses Admin & Kasir)
			protected.GET("/products", controllers.GetProducts)
			protected.GET("/categories", controllers.GetCategories)
			protected.GET("/payments", controllers.GetPayments)
			
			// Transactions (Hanya Kasir yang bisa buat transaksi)
			protected.POST("/transactions", middleware.RoleMiddleware("kasir"), controllers.CreateTransaction)
			protected.GET("/transactions", controllers.GetTransactions) // History bisa dilihat semua user login

			// Admin Only (Manage Produk, Payment, User)
			admin := protected.Group("/")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				// Product Management
				admin.POST("/products", controllers.CreateProduct)
				admin.PUT("/products/:id", controllers.UpdateProduct)
				admin.DELETE("/products/:id", controllers.DeleteProduct)

				// Category Management
				admin.POST("/categories", controllers.CreateCategory)
				admin.PUT("/categories/:id", controllers.UpdateCategory)
				admin.DELETE("/categories/:id", controllers.DeleteCategory)

				// Payment Management
				admin.POST("/payments", controllers.CreatePayment)
				admin.PUT("/payments/:id", controllers.UpdatePayment)
				admin.DELETE("/payments/:id", controllers.DeletePayment)

				// User Management
				// URL: /api/admin/users
				admin.GET("/admin/users", controllers.GetUsers) 
				admin.POST("/admin/users", controllers.CreateUser)
				admin.PUT("/admin/users/:id", controllers.UpdateUser)
				admin.DELETE("/admin/users/:id", controllers.DeleteUser)
			}
		}
	}

	return r
}
