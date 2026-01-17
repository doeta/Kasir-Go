package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/doeta/Kasir-Go/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var dsn string
	if os.Getenv("DATABASE_URL") != "" {
		dsn = os.Getenv("DATABASE_URL")
	} else {
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")
		sslmode := os.Getenv("DB_SSLMODE")
		if sslmode == "" {
			sslmode = "disable"
		}

		// String koneksi
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta", host, user, password, dbname, port, sslmode)
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi database: ", err)
	}

	// Auto Migrate (Otomatis bikin tabel di database)
	fmt.Println("Migrasi Database berjalan...")
	DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Category{}, &models.Transaction{}, &models.TransactionDetail{}, &models.PaymentMethod{})
	
	seedAdmin()
}

// Fungsi bikin admin otomatis kalau belum ada
func seedAdmin() {
    var user models.User
    // Cek apakah sudah ada user di database?
    if err := DB.First(&user).Error; err != nil {
        // Kalau error (berarti kosong), kita buat admin baru
        fmt.Println("User kosong, membuat admin otomatis...")
        
        hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
        
        admin := models.User{
            Username: "admin",
            Password: string(hashedPassword),
            Role:     "admin",
        }
        
        DB.Create(&admin)
        fmt.Println("Admin dibuat! Username: 'admin', Password: 'admin123'")
    }
}
