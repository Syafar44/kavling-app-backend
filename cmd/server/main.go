// Package main Kavling Mentaya API
//
//	@title						Kavling Mentaya API
//	@version					1.0
//	@description				REST API backend untuk sistem manajemen penjualan kavling tanah. Semua endpoint (kecuali /auth/login) memerlukan Bearer Token JWT di header Authorization.
//	@termsOfService				http://swagger.io/terms/
//	@contact.name				API Support
//	@contact.email				admin@kavling-mentaya.com
//	@license.name				MIT
//	@host						localhost:8080
//	@BasePath					/api
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Masukkan token JWT dengan format: Bearer {token}. Token diperoleh dari endpoint POST /api/auth/login.
package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "backend-kavling/docs"
	"backend-kavling/internal/config"
	"backend-kavling/internal/handlers"
	"backend-kavling/internal/middleware"
	"backend-kavling/internal/repositories"
	"backend-kavling/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// 1. Load konfigurasi dari .env
	config.Load()

	// 2. Koneksi database
	db := config.ConnectDB()

	// 3. Buat folder upload jika belum ada
	ensureUploadDirs(config.AppConfig.UploadPath)

	// 4. Init semua repository
	userRepo := repositories.NewUserRepository(db)
	kavlingRepo := repositories.NewKavlingRepository(db)
	denahRepo := repositories.NewDenahRepository(db)
	marketingRepo := repositories.NewMarketingRepository(db)
	customerRepo := repositories.NewCustomerRepository(db)
	bankRepo := repositories.NewBankRepository(db)
	konfigRepo := repositories.NewKonfigurasiRepository(db)
	bookingRepo := repositories.NewBookingRepository(db)
	trxKavRepo := repositories.NewTransaksiKavlingRepository(db)
	pembayaranRepo := repositories.NewPembayaranRepository(db)
	arusKasRepo := repositories.NewArusKasRepository(db)

	// 5. Init semua service
	authService := services.NewAuthService(userRepo)
	kavlingService := services.NewKavlingService(kavlingRepo)
	denahService := services.NewDenahService(denahRepo, kavlingRepo)
	transaksiService := services.NewTransaksiService(
		kavlingRepo, bookingRepo, trxKavRepo, pembayaranRepo, arusKasRepo,
	)

	// 6. Init semua handler
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userRepo, authService)
	kavlingHandler := handlers.NewKavlingHandler(kavlingRepo, kavlingService)
	denahHandler := handlers.NewDenahHandler(denahRepo, denahService)
	marketingHandler := handlers.NewMarketingHandler(marketingRepo)
	customerHandler := handlers.NewCustomerHandler(customerRepo)
	bankHandler := handlers.NewBankHandler(bankRepo)
	konfigHandler := handlers.NewKonfigurasiHandler(konfigRepo)
	bookingHandler := handlers.NewBookingHandler(bookingRepo, transaksiService)
	transaksiHandler := handlers.NewTransaksiHandler(trxKavRepo, transaksiService)
	pembayaranHandler := handlers.NewPembayaranHandler(pembayaranRepo, trxKavRepo, transaksiService)
	keuanganHandler := handlers.NewKeuanganHandler(arusKasRepo, transaksiService)
	laporanHandler := handlers.NewLaporanHandler(arusKasRepo)

	// 7. Setup Gin router
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Recovery())

	// Serve static files (upload)
	r.Static("/uploads", config.AppConfig.UploadPath)

	// Health check endpoint (untuk Docker HEALTHCHECK & Railway)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Swagger UI — akses di http://localhost:8080/swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 8. Define routes
	api := r.Group("/api")

	// ─── Public routes ────────────────────────────────────────────────────────
	api.POST("/auth/login", authHandler.Login)

	// ─── Protected routes ─────────────────────────────────────────────────────
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Dashboard / statistik
		protected.GET("/statistik", keuanganHandler.Statistik)

		// Menu
		protected.GET("/menu", userHandler.ListMenu)

		// ─── User management (admin only) ─────────────────────────────────────
		users := protected.Group("/users")
		users.Use(middleware.AdminOnly())
		{
			users.GET("", userHandler.List)
			users.POST("", userHandler.Create)
			users.GET("/:id", userHandler.Detail)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
			users.GET("/:id/access", userHandler.GetAccess)
			users.PUT("/:id/access", userHandler.UpdateAccess)
		}

		// Activity log (admin only)
		protected.GET("/activity-log", middleware.AdminOnly(), userHandler.ActivityLog)

		// ─── Denah Kavling ───────────────────────────────────────────────────
		denah := protected.Group("/denah-kavling")
		{
			denah.GET("", denahHandler.List)
			denah.POST("", denahHandler.Create)
			denah.GET("/:id", denahHandler.Detail)
			denah.PUT("/:id", denahHandler.Update)
			denah.DELETE("/:id", denahHandler.Delete)
		}

		// ─── Kavling ──────────────────────────────────────────────────────────
		kavling := protected.Group("/kavling")
		{
			kavling.GET("", kavlingHandler.List)
			kavling.GET("/:id", kavlingHandler.Detail)
			kavling.PUT("/:id", kavlingHandler.Update)
			kavling.DELETE("/:id", kavlingHandler.Delete)
		}

		// ─── Marketing ────────────────────────────────────────────────────────
		marketing := protected.Group("/marketing")
		{
			marketing.GET("", marketingHandler.List)
			marketing.POST("", marketingHandler.Create)
			marketing.GET("/:id", marketingHandler.Detail)
			marketing.PUT("/:id", marketingHandler.Update)
			marketing.DELETE("/:id", marketingHandler.Delete)
		}

		// ─── Customer ─────────────────────────────────────────────────────────
		customers := protected.Group("/customers")
		{
			customers.GET("", customerHandler.List)
			customers.POST("", customerHandler.Create)
			customers.GET("/:id", customerHandler.Detail)
			customers.PUT("/:id", customerHandler.Update)
			customers.DELETE("/:id", customerHandler.Delete)
		}

		// ─── Bank ─────────────────────────────────────────────────────────────
		bank := protected.Group("/bank")
		{
			bank.GET("", bankHandler.List)
			bank.POST("", bankHandler.Create)
			bank.PUT("/:id", bankHandler.Update)
			bank.DELETE("/:id", bankHandler.Delete)
		}

		// ─── Booking ──────────────────────────────────────────────────────────
		booking := protected.Group("/booking")
		{
			booking.GET("", bookingHandler.List)
			booking.POST("", bookingHandler.Create)
			booking.DELETE("/:id", bookingHandler.Cancel)
			booking.POST("/:id/convert", bookingHandler.Convert)
		}

		// ─── Transaksi Kavling ────────────────────────────────────────────────
		transaksi := protected.Group("/transaksi")
		{
			transaksi.GET("", transaksiHandler.List)
			transaksi.GET("/:id", transaksiHandler.Detail)
			transaksi.POST("/cash", transaksiHandler.CreateCash)
			transaksi.POST("/kredit", transaksiHandler.CreateKredit)
		}

		// ─── Pembayaran ───────────────────────────────────────────────────────
		pembayaran := protected.Group("/pembayaran")
		{
			pembayaran.GET("", pembayaranHandler.List)
			pembayaran.GET("/:id_kavling", pembayaranHandler.DetailKavling)
			pembayaran.POST("/:id_kavling/bayar", pembayaranHandler.Bayar)
		}

		// ─── Keuangan ─────────────────────────────────────────────────────────
		keuangan := protected.Group("/keuangan")
		{
			keuangan.GET("/transaksi", keuanganHandler.ListArusKas)
			keuangan.POST("/transaksi", keuanganHandler.CreateArusKas)
			keuangan.DELETE("/transaksi/:id", keuanganHandler.DeleteArusKas)
			keuangan.GET("/rekap-kredit", keuanganHandler.RekapKredit)
			keuangan.GET("/statistik", keuanganHandler.Statistik)
		}

		// ─── Laporan ──────────────────────────────────────────────────────────
		laporan := protected.Group("/laporan")
		{
			laporan.GET("/umum", laporanHandler.Umum)
			laporan.GET("/per-customer", laporanHandler.PerCustomer)
			laporan.GET("/per-kategori", laporanHandler.PerKategori)
			laporan.GET("/export-excel", laporanHandler.ExportExcel)
			laporan.GET("/export-excel-pembayaran", laporanHandler.ExportExcelPembayaran)
		}

		// ─── Konfigurasi ──────────────────────────────────────────────────────
		konfig := protected.Group("/konfigurasi")
		{
			konfig.GET("", konfigHandler.Get)
			konfig.PUT("", konfigHandler.Update)
			konfig.GET("/wa", konfigHandler.GetWA)
			konfig.PUT("/wa", konfigHandler.UpdateWA)
			konfig.GET("/template", konfigHandler.ListTemplate)
			konfig.POST("/template", konfigHandler.CreateTemplate)
			konfig.PUT("/template/:id", konfigHandler.UpdateTemplate)
			konfig.POST("/reset", konfigHandler.ResetData)
		}
	}

	// 9. Jalankan server
	port := config.AppConfig.Port
	log.Printf("Server berjalan di port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}

func ensureUploadDirs(base string) {
	dirs := []string{
		"lampiran_customer",
		"bukti_trx/bayar_angsuran",
		"lampiran_trx",
		"konfigurasi",
	}
	for _, d := range dirs {
		path := filepath.Join(base, d)
		if err := os.MkdirAll(path, 0755); err != nil {
			log.Printf("Warning: gagal membuat folder %s: %v", path, err)
		}
	}
}
