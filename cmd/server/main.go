// Package main Kavling Mentaya API
//
//	@title						Kavling Mentaya API
//	@version					2.0
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
	_ "time/tzdata"

	"backend-kavling/docs"
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
	config.Load()
	db := config.ConnectDB()
	ensureUploadDirs(config.AppConfig.UploadPath)

	// ─── Repositories ─────────────────────────────────────────────────────────
	userRepo         := repositories.NewUserRepository(db)
	kavlingRepo      := repositories.NewKavlingRepository(db)
	denahRepo        := repositories.NewDenahRepository(db)
	marketingRepo    := repositories.NewMarketingRepository(db)
	customerRepo     := repositories.NewCustomerRepository(db)
	bankRepo         := repositories.NewBankRepository(db)
	konfigRepo       := repositories.NewKonfigurasiRepository(db)
	bookingRepo      := repositories.NewBookingRepository(db)
	trxKavRepo       := repositories.NewTransaksiKavlingRepository(db)
	pembayaranRepo   := repositories.NewPembayaranRepository(db)
	arusKasRepo      := repositories.NewArusKasRepository(db)

	// ─── New Repositories (v2) ────────────────────────────────────────────────
	lokasiRepo       := repositories.NewLokasiKavlingRepository(db)
	kavlingV2Repo    := repositories.NewKavlingV2Repository(db)
	kategoriRepo     := repositories.NewKategoriRepository(db)
	tagihanRepo      := repositories.NewTagihanRepository(db)
	keuanganRepo     := repositories.NewKeuanganRepository(db)
	customerV2Repo   := repositories.NewCustomerV2Repository(db)
	legalitasRepo    := repositories.NewLegalitasRepository(db)
	notarisRepo      := repositories.NewNotarisRepository(db)
	listPenjualanRepo := repositories.NewListPenjualanRepository(db)
	mediaRepo        := repositories.NewMediaRepository(db)
	landingRepo      := repositories.NewLandingRepository(db)
	berandaRepo      := repositories.NewBerandaRepository(db)
	pembayaranV2Repo := repositories.NewPembayaranV2Repository(db)
	jatuhTempoRepo   := repositories.NewJatuhTempoRepository(db)


	// NOTE: BackfillFromCustomers, FixTagihanBF, BackfillBFPembayaran DIHAPUS.
	// Semua logika pembayaran sekarang di customer_handler_v2.go → generatePembayaran()


	// ─── Services ─────────────────────────────────────────────────────────────
	authService      := services.NewAuthService(userRepo)
	kavlingService   := services.NewKavlingService(kavlingRepo)
	denahService     := services.NewDenahService(denahRepo, kavlingRepo)
	transaksiService := services.NewTransaksiService(
		kavlingRepo, bookingRepo, trxKavRepo, pembayaranRepo, arusKasRepo,
	)

	// ─── Handlers (legacy) ────────────────────────────────────────────────────
	authHandler      := handlers.NewAuthHandler(authService)
	userHandler      := handlers.NewUserHandler(userRepo, authService)
	kavlingHandler   := handlers.NewKavlingHandler(kavlingRepo, kavlingService)
	denahHandler     := handlers.NewDenahHandler(denahRepo, denahService)
	marketingHandler := handlers.NewMarketingHandler(marketingRepo)
	customerHandler  := handlers.NewCustomerHandler(customerRepo)
	bankHandler      := handlers.NewBankHandler(bankRepo)
	konfigHandler    := handlers.NewKonfigurasiHandler(konfigRepo)
	bookingHandler   := handlers.NewBookingHandler(bookingRepo, transaksiService)
	transaksiHandler := handlers.NewTransaksiHandler(trxKavRepo, transaksiService)
	pembayaranHandler := handlers.NewPembayaranHandler(pembayaranRepo, trxKavRepo, transaksiService)
	keuanganHandler  := handlers.NewKeuanganHandler(arusKasRepo, transaksiService)
	laporanHandler   := handlers.NewLaporanHandler(arusKasRepo)

	// ─── Handlers (v2 — new modules) ──────────────────────────────────────────
	lokasiHandler    := handlers.NewLokasiKavlingHandler(lokasiRepo, kavlingV2Repo)
	siteplanHandler  := handlers.NewSiteplanHandler(lokasiRepo, kavlingV2Repo)
	berandaHandler   := handlers.NewBerandaHandler(berandaRepo, keuanganRepo)
	kavlingV2Handler := handlers.NewKavlingV2Handler(kavlingV2Repo, lokasiRepo)
	tagihanHandler   := handlers.NewTagihanHandler(tagihanRepo, pembayaranV2Repo, kavlingV2Repo, keuanganRepo, trxKavRepo)
	keuanganV2Handler := handlers.NewKeuanganV2Handler(keuanganRepo, kategoriRepo)
	customerV2Handler := handlers.NewCustomerV2Handler(customerV2Repo, kavlingV2Repo, lokasiRepo, tagihanRepo, trxKavRepo, keuanganRepo)
	legalitasHandler := handlers.NewLegalitasHandler(legalitasRepo)
	notarisHandler   := handlers.NewNotarisHandler(notarisRepo)
	listPenjualanHandler := handlers.NewListPenjualanHandler(listPenjualanRepo)
	mediaHandler     := handlers.NewMediaHandler(mediaRepo)
	landingHandler   := handlers.NewLandingHandler(landingRepo)
	hakAksesHandler  := handlers.NewHakAksesExtHandler(userRepo, db)
	jatuhTempoHandler := handlers.NewJatuhTempoHandler(jatuhTempoRepo)

	// ─── Router setup ─────────────────────────────────────────────────────────
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Recovery())

	r.Static("/uploads", config.AppConfig.UploadPath)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().Format(time.RFC3339)})
	})

	docs.SwaggerInfo.Host = os.Getenv("SWAGGER_HOST")
	docs.SwaggerInfo.Schemes = []string{"https", "http"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	api.POST("/auth/login", authHandler.Login)

	// Public: landing page
	api.GET("/public/landing", landingHandler.List)

	// ─── Protected ────────────────────────────────────────────────────────────
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// ── Beranda / Dashboard ──────────────────────────────────────────────
		protected.GET("/beranda/ringkasan",  berandaHandler.Ringkasan)
		protected.GET("/beranda/arus-kas",   berandaHandler.ArusKas)
		protected.GET("/beranda/aktifitas",  berandaHandler.Aktifitas)
		// Legacy statistik endpoint
		protected.GET("/statistik",          keuanganHandler.Statistik)

		// ── Menu ──────────────────────────────────────────────────────────────
		protected.GET("/menu", userHandler.ListMenu)

		// ── Siteplan ──────────────────────────────────────────────────────────
		protected.GET("/siteplan",                                   siteplanHandler.List)
		protected.GET("/siteplan/:id_lokasi",                        siteplanHandler.Detail)
		protected.GET("/siteplan/:id_lokasi/kavling/:kode",          siteplanHandler.KavlingDetail)

		// ── Lokasi Kavling (master) ────────────────────────────────────────────
		lokasi := protected.Group("/lokasi-kavling")
		{
			lokasi.GET("",      lokasiHandler.List)
			lokasi.GET("/:id",  lokasiHandler.Detail)
			lokasi.POST("",     lokasiHandler.Create)
			lokasi.PUT("/:id",  lokasiHandler.Update)
			lokasi.DELETE("/:id", lokasiHandler.Delete)
		}

		// ── Kavling V2 ────────────────────────────────────────────────────────
		kavlingV2 := protected.Group("/kavling")
		{
			kavlingV2.GET("",            kavlingV2Handler.List)
			kavlingV2.GET("/:id",        kavlingV2Handler.Detail)
			kavlingV2.PUT("/:id",        kavlingV2Handler.Update)
			kavlingV2.POST("/bulk",      kavlingV2Handler.BulkCreate)
			// Legacy routes (kavling_peta based)
			kavlingV2.DELETE("/:id",     kavlingHandler.Delete)
		}

		// ── User management ────────────────────────────────────────────────────
		users := protected.Group("/users")
		users.Use(middleware.AdminOnly())
		{
			users.GET("",           userHandler.List)
			users.POST("",          userHandler.Create)
			users.GET("/:id",       userHandler.Detail)
			users.PUT("/:id",       userHandler.Update)
			users.DELETE("/:id",    userHandler.Delete)
			users.GET("/:id/access", userHandler.GetAccess)
			users.PUT("/:id/access", userHandler.UpdateAccess)
		}
		protected.GET("/activity-log", middleware.AdminOnly(), userHandler.ActivityLog)

		// ── Hak Akses (matrix) ────────────────────────────────────────────────
		hakAkses := protected.Group("/hak-akses")
		{
			hakAkses.GET("",       hakAksesHandler.GetMatrix)
			hakAkses.PUT("/bulk",  hakAksesHandler.BulkUpdate)
		}

		// ── Denah Kavling (legacy) ────────────────────────────────────────────
		denah := protected.Group("/denah-kavling")
		{
			denah.GET("",      denahHandler.List)
			denah.POST("",     denahHandler.Create)
			denah.GET("/:id",  denahHandler.Detail)
			denah.PUT("/:id",  denahHandler.Update)
			denah.DELETE("/:id", denahHandler.Delete)
		}

		// ── Marketing ─────────────────────────────────────────────────────────
		marketing := protected.Group("/marketing")
		{
			marketing.GET("",      marketingHandler.List)
			marketing.POST("",     marketingHandler.Create)
			marketing.GET("/:id",  marketingHandler.Detail)
			marketing.PUT("/:id",  marketingHandler.Update)
			marketing.DELETE("/:id", marketingHandler.Delete)
		}

		// ── Customer V2 ────────────────────────────────────────────────────────
		customers := protected.Group("/customer")
		{
			customers.GET("",      customerV2Handler.List)
			customers.POST("",     customerV2Handler.Create)
			customers.GET("/:id",  customerV2Handler.Detail)
			customers.PUT("/:id",  customerV2Handler.Update)
			customers.DELETE("/:id", customerV2Handler.Delete)
			// Files
			customers.GET("/:id/file",                 customerV2Handler.ListFile)
			customers.POST("/:id/file",                customerV2Handler.UploadFile)
			customers.DELETE("/:id/file/:id_file",     customerV2Handler.DeleteFile)
		}
		// Legacy /customers endpoint
		legacyCust := protected.Group("/customers")
		{
			legacyCust.GET("",      customerHandler.List)
			legacyCust.POST("",     customerHandler.Create)
			legacyCust.GET("/:id",  customerHandler.Detail)
			legacyCust.PUT("/:id",  customerHandler.Update)
			legacyCust.DELETE("/:id", customerHandler.Delete)
		}

		// ── Arsip Customer ────────────────────────────────────────────────────
		arsip := protected.Group("/customer/arsip")
		{
			arsip.GET("",      customerV2Handler.ListArsip)
			arsip.DELETE("/:id", customerV2Handler.DeleteArsip)
		}

		// ── Prospek ───────────────────────────────────────────────────────────
		prospek := protected.Group("/prospek")
		{
			prospek.GET("",          customerV2Handler.ListProspek)
			prospek.POST("",         customerV2Handler.CreateProspek)
			prospek.PUT("/:id",      customerV2Handler.UpdateProspek)
			prospek.DELETE("/:id",   customerV2Handler.DeleteProspek)
			prospek.POST("/:id/convert", customerV2Handler.ConvertProspek)
		}

		// ── Legalitas ─────────────────────────────────────────────────────────
		legalitas := protected.Group("/legalitas")
		{
			legalitas.GET("",      legalitasHandler.List)
			legalitas.GET("/:id",  legalitasHandler.Detail)
			legalitas.PUT("/:id",  legalitasHandler.Update)
		}

		// ── Notaris ───────────────────────────────────────────────────────────
		notaris := protected.Group("/notaris")
		{
			notaris.GET("",      notarisHandler.List)
			notaris.POST("",     notarisHandler.Create)
			notaris.PUT("/:id",  notarisHandler.Update)
			notaris.DELETE("/:id", notarisHandler.Delete)
		}

		// ── Bank ──────────────────────────────────────────────────────────────
		bank := protected.Group("/bank")
		{
			bank.GET("",      bankHandler.List)
			bank.POST("",     bankHandler.Create)
			bank.PUT("/:id",  bankHandler.Update)
			bank.DELETE("/:id", bankHandler.Delete)
		}

		// ── Booking (legacy) ──────────────────────────────────────────────────
		booking := protected.Group("/booking")
		{
			booking.GET("",           bookingHandler.List)
			booking.POST("",          bookingHandler.Create)
			booking.DELETE("/:id",    bookingHandler.Cancel)
			booking.POST("/:id/convert", bookingHandler.Convert)
		}

		// ── Transaksi Kavling (legacy) ─────────────────────────────────────────
		transaksi := protected.Group("/transaksi")
		{
			transaksi.GET("",         transaksiHandler.List)
			transaksi.GET("/:id",     transaksiHandler.Detail)
			transaksi.POST("/cash",   transaksiHandler.CreateCash)
			transaksi.POST("/kredit", transaksiHandler.CreateKredit)
		}

		// ── Pembayaran V2 (tagihan + pemasukan) ───────────────────────────────
		pembayaran := protected.Group("/pembayaran")
		{
			pembayaran.GET("",                                     tagihanHandler.List)
			pembayaran.GET("/debug",                               tagihanHandler.Debug)
			pembayaran.GET("/:id_transaksi",                       tagihanHandler.DetailPembayaran)
			pembayaran.POST("/:id_transaksi/tagihan",              tagihanHandler.TambahTagihan)
			pembayaran.PUT("/:id_transaksi/tagihan/:id",           tagihanHandler.UpdateTagihan)
			pembayaran.DELETE("/:id_transaksi/tagihan/:id",        tagihanHandler.HapusTagihan)
			pembayaran.POST("/:id_transaksi/pemasukan",            tagihanHandler.TambahPemasukan)
			pembayaran.DELETE("/:id_transaksi/pemasukan/:id",      tagihanHandler.HapusPemasukan)
		}
		// Legacy bayar (separate to avoid wildcard conflict)
		protected.POST("/pembayaran-legacy/:id_kavling/bayar", pembayaranHandler.Bayar)

		// ── Jatuh Tempo ───────────────────────────────────────────────────────
		protected.GET("/jatuh-tempo", jatuhTempoHandler.List)

		// ── Keuangan V2 ────────────────────────────────────────────────────────
		// Kategori Transaksi
		kategori := protected.Group("/kategori-transaksi")
		{
			kategori.GET("",      keuanganV2Handler.ListKategori)
			kategori.POST("",     keuanganV2Handler.CreateKategori)
			kategori.PUT("/:id",  keuanganV2Handler.UpdateKategori)
			kategori.DELETE("/:id", keuanganV2Handler.DeleteKategori)
		}

		// Pemasukan
		pemasukan := protected.Group("/pemasukan")
		{
			pemasukan.GET("",      keuanganV2Handler.ListPemasukan)
			pemasukan.POST("",     keuanganV2Handler.CreatePemasukan)
			pemasukan.DELETE("/:id", keuanganV2Handler.DeletePemasukan)
		}

		// Pengeluaran
		pengeluaran := protected.Group("/pengeluaran")
		{
			pengeluaran.GET("",      keuanganV2Handler.ListPengeluaran)
			pengeluaran.POST("",     keuanganV2Handler.CreatePengeluaran)
			pengeluaran.DELETE("/:id", keuanganV2Handler.DeletePengeluaran)
		}

		// Hutang
		hutang := protected.Group("/hutang")
		{
			hutang.GET("",          keuanganV2Handler.ListHutang)
			hutang.POST("",         keuanganV2Handler.CreateHutang)
			hutang.PUT("/:id",      keuanganV2Handler.UpdateHutang)
			hutang.DELETE("/:id",   keuanganV2Handler.DeleteHutang)
			hutang.POST("/:id/bayar", keuanganV2Handler.BayarHutang)
		}

		// Piutang
		piutang := protected.Group("/piutang")
		{
			piutang.GET("",          keuanganV2Handler.ListPiutang)
			piutang.POST("",         keuanganV2Handler.CreatePiutang)
			piutang.DELETE("/:id",   keuanganV2Handler.DeletePiutang)
			piutang.POST("/:id/bayar", keuanganV2Handler.BayarPiutang)
		}

		// Mutasi Saldo
		mutasi := protected.Group("/mutasi-saldo")
		{
			mutasi.GET("",      keuanganV2Handler.ListMutasi)
			mutasi.POST("",     keuanganV2Handler.CreateMutasi)
			mutasi.DELETE("/:id", keuanganV2Handler.DeleteMutasi)
		}

		// Laporan Arus Kas
		protected.GET("/laporan/arus-kas", keuanganV2Handler.LaporanArusKas)

		// Legacy keuangan routes
		keuangan := protected.Group("/keuangan")
		{
			keuangan.GET("/transaksi",     keuanganHandler.ListArusKas)
			keuangan.POST("/transaksi",    keuanganHandler.CreateArusKas)
			keuangan.DELETE("/transaksi/:id", keuanganHandler.DeleteArusKas)
			keuangan.GET("/rekap-kredit",  keuanganHandler.RekapKredit)
			keuangan.GET("/statistik",     keuanganHandler.Statistik)
		}

		// Legacy laporan
		laporan := protected.Group("/laporan")
		{
			laporan.GET("/umum",                    laporanHandler.Umum)
			laporan.GET("/per-customer",            laporanHandler.PerCustomer)
			laporan.GET("/per-kategori",            laporanHandler.PerKategori)
			laporan.GET("/export-excel",            laporanHandler.ExportExcel)
			laporan.GET("/export-excel-pembayaran", laporanHandler.ExportExcelPembayaran)
		}

		// ── List Penjualan ────────────────────────────────────────────────────
		listPenjualan := protected.Group("/list-penjualan")
		{
			listPenjualan.GET("",      listPenjualanHandler.List)
			listPenjualan.POST("",     listPenjualanHandler.Create)
			listPenjualan.PUT("/:id",  listPenjualanHandler.Update)
			listPenjualan.DELETE("/:id", listPenjualanHandler.Delete)
		}

		// ── Media ─────────────────────────────────────────────────────────────
		media := protected.Group("/media")
		{
			media.GET("",      mediaHandler.List)
			media.PUT("/:id",  mediaHandler.Update)
		}

		// ── Landing ───────────────────────────────────────────────────────────
		landing := protected.Group("/landing")
		{
			landing.GET("",      landingHandler.List)
			landing.POST("",     landingHandler.Create)
			landing.PUT("/:id",  landingHandler.Update)
			landing.DELETE("/:id", landingHandler.Delete)
		}

		// ── Konfigurasi ───────────────────────────────────────────────────────
		konfig := protected.Group("/konfigurasi")
		{
			konfig.GET("",              konfigHandler.Get)
			konfig.PUT("",              konfigHandler.Update)
			konfig.GET("/wa",           konfigHandler.GetWA)
			konfig.PUT("/wa",           konfigHandler.UpdateWA)
			konfig.GET("/template",     konfigHandler.ListTemplate)
			konfig.POST("/template",    konfigHandler.CreateTemplate)
			konfig.PUT("/template/:id", konfigHandler.UpdateTemplate)
			konfig.POST("/reset",       konfigHandler.ResetData)
		}
	}

	port := config.AppConfig.Port
	log.Printf("Server berjalan di port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}

func ensureUploadDirs(base string) {
	dirs := []string{
		"customer",
		"legalitas",
		"pembayaran",
		"hutang",
		"piutang",
		"mutasi",
		"media",
		"marketing",
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
