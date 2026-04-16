// Seed: jalankan sekali untuk membuat user admin pertama
// Usage: go run cmd/seed/main.go
package main

import (
	"fmt"
	"log"

	"backend-kavling/internal/config"
	"backend-kavling/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	config.Load()
	db := config.ConnectDB()

	// Buat user admin default
	password := "admin123"
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	admin := models.User{
		Username: "admin",
		Password: string(hashed),
		Nama:     "Administrator",
		IsAdmin:  1,
		Status:   "AKTIF",
	}

	result := db.Where("username = ?", "admin").FirstOrCreate(&admin)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	if result.RowsAffected > 0 {
		fmt.Println("User admin berhasil dibuat:")
		fmt.Printf("  Username: admin\n")
		fmt.Printf("  Password: %s\n", password)
		fmt.Println("  SEGERA ganti password setelah pertama login!")
	} else {
		fmt.Println("User admin sudah ada, seed diabaikan.")
	}
}
