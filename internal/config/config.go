package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	Timezone   string
	UploadPath string
	WAApiURL   string
	WAApiKey   string
	Port       string
}

var AppConfig Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	AppConfig = Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "kavling_db"),
		JWTSecret:  getEnv("JWT_SECRET", "kavling-secret-key-ganti-ini"),
		Timezone:   getEnv("TIMEZONE", "Asia/Makassar"),
		UploadPath: getEnv("UPLOAD_PATH", "./uploads"),
		WAApiURL:   getEnv("WA_API_URL", ""),
		WAApiKey:   getEnv("WA_API_KEY", ""),
		Port:       getEnv("PORT", "8080"),
	}

	// Set timezone
	loc, err := time.LoadLocation(AppConfig.Timezone)
	if err != nil {
		log.Printf("Warning: Cannot load timezone %s, using UTC", AppConfig.Timezone)
	} else {
		time.Local = loc
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		AppConfig.DBHost,
		AppConfig.DBPort,
		AppConfig.DBUser,
		AppConfig.DBPassword,
		AppConfig.DBName,
		AppConfig.Timezone,
	)

	loc, err := time.LoadLocation(AppConfig.Timezone)
	if err != nil {
		log.Printf("Warning: timezone %s tidak dikenal, fallback ke UTC", AppConfig.Timezone)
		loc = time.UTC
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().In(loc)
		},
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Database connected successfully")
	return db
}
