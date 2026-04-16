package handlers

import (
	"path/filepath"
	"strconv"
	"time"

	"backend-kavling/internal/config"
	"backend-kavling/internal/helpers"
	"backend-kavling/internal/repositories"
	"backend-kavling/internal/services"

	"github.com/gin-gonic/gin"
)

// ─── Booking Handler ──────────────────────────────────────────────────────────

type BookingHandler struct {
	repo    *repositories.BookingRepository
	service *services.TransaksiService
}

func NewBookingHandler(repo *repositories.BookingRepository, service *services.TransaksiService) *BookingHandler {
	return &BookingHandler{repo: repo, service: service}
}

// List godoc
//
//	@Summary		List booking aktif
//	@Description	Mengambil daftar semua booking yang masih aktif (status=1)
//	@Tags			Booking
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object} object
//	@Router			/booking [get]
func (h *BookingHandler) List(c *gin.Context) {
	list, err := h.repo.FindAll()
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil data booking")
		return
	}
	helpers.OK(c, "OK", list)
}

// Create godoc
//
//	@Summary		Buat booking baru
//	@Description	Membuat booking kavling. Kavling harus berstatus Kosong (0). Status kavling berubah menjadi Booking (1).
//	@Tags			Booking
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		services.BookingInput	true	"Data booking"
//	@Success		201		{object} object
//	@Failure		400		{object} object
//	@Router			/booking [post]
func (h *BookingHandler) Create(c *gin.Context) {
	var input services.BookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.BadRequest(c, "Data tidak lengkap", err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(int)
	input.IDUser = &uid

	booking, err := h.service.CreateBooking(input)
	if err != nil {
		helpers.BadRequest(c, err.Error(), nil)
		return
	}

	helpers.Created(c, "Booking berhasil dibuat", booking)
}

// Cancel godoc
//
//	@Summary		Batalkan booking
//	@Description	Membatalkan booking aktif. Status kavling dikembalikan ke Kosong (0).
//	@Tags			Booking
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Booking ID"
//	@Success		200	{object} object
//	@Failure		400	{object} object
//	@Router			/booking/{id} [delete]
func (h *BookingHandler) Cancel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.CancelBooking(id); err != nil {
		helpers.BadRequest(c, err.Error(), nil)
		return
	}

	helpers.OK(c, "Booking berhasil dibatalkan", nil)
}

// Convert godoc
//
//	@Summary		Konversi booking ke transaksi
//	@Description	Mengubah booking menjadi transaksi cash atau kredit. Jika kredit, cicilan per bulan otomatis dihitung = (harga_jual - uang_muka) / lama_cicilan.
//	@Tags			Booking
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"Booking ID"
//	@Param			body	body		services.ConvertInput	true	"Data konversi (jenis: cash|kredit)"
//	@Success		201		{object} object
//	@Failure		400		{object} object
//	@Router			/booking/{id}/convert [post]
func (h *BookingHandler) Convert(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var input services.ConvertInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.BadRequest(c, "Data tidak lengkap", err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(int)
	input.IDUser = &uid

	trx, err := h.service.ConvertBooking(id, input)
	if err != nil {
		helpers.BadRequest(c, err.Error(), nil)
		return
	}

	helpers.Created(c, "Booking berhasil dikonversi", trx)
}

// ─── Transaksi Cash/Kredit Handler ────────────────────────────────────────────

type TransaksiHandler struct {
	repo    *repositories.TransaksiKavlingRepository
	service *services.TransaksiService
}

func NewTransaksiHandler(repo *repositories.TransaksiKavlingRepository, service *services.TransaksiService) *TransaksiHandler {
	return &TransaksiHandler{repo: repo, service: service}
}

// List godoc
//
//	@Summary		List semua transaksi kavling
//	@Description	Mengambil daftar semua transaksi jual-beli kavling. Gunakan query ?jenis=kredit untuk filter kredit saja.
//	@Tags			Transaksi Kavling
//	@Produce		json
//	@Security		BearerAuth
//	@Param			jenis	query	string	false	"Filter jenis: kredit"
//	@Success		200		{object} object
//	@Router			/transaksi [get]
func (h *TransaksiHandler) List(c *gin.Context) {
	if c.Query("jenis") == "kredit" {
		list, err := h.repo.FindAllKredit()
		if err != nil {
			helpers.InternalError(c, "Gagal mengambil data transaksi kredit")
			return
		}
		helpers.OK(c, "OK", list)
		return
	}
	list, err := h.repo.FindAll()
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil data transaksi")
		return
	}
	helpers.OK(c, "OK", list)
}

// Detail godoc
//
//	@Summary		Detail transaksi kavling
//	@Description	Mengambil detail satu transaksi kavling berdasarkan ID
//	@Tags			Transaksi Kavling
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Transaksi ID"
//	@Success		200	{object} object
//	@Failure		404	{object} object
//	@Router			/transaksi/{id} [get]
func (h *TransaksiHandler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	trx, err := h.repo.FindByID(id)
	if err != nil {
		helpers.NotFound(c, "Transaksi tidak ditemukan")
		return
	}
	helpers.OK(c, "OK", trx)
}

// CreateCash godoc
//
//	@Summary		Transaksi cash langsung
//	@Description	Membuat transaksi pembelian cash langsung dari peta kavling (tanpa booking). Kavling harus berstatus Kosong (0).
//	@Tags			Transaksi Kavling
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		services.CashInput	true	"Data transaksi cash"
//	@Success		201		{object} object
//	@Failure		400		{object} object
//	@Router			/transaksi/cash [post]
func (h *TransaksiHandler) CreateCash(c *gin.Context) {
	var input services.CashInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.BadRequest(c, "Data tidak lengkap", err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(int)
	input.IDUser = &uid

	trx, err := h.service.CreateCash(input)
	if err != nil {
		helpers.BadRequest(c, err.Error(), nil)
		return
	}

	helpers.Created(c, "Transaksi cash berhasil dibuat", trx)
}

// CreateKredit godoc
//
//	@Summary		Transaksi kredit langsung
//	@Description	Membuat transaksi pembelian kredit langsung dari peta kavling. Cicilan per bulan = (harga_jual - uang_muka) / lama_cicilan. Tgl jatuh tempo default tanggal 10 bulan berikutnya.
//	@Tags			Transaksi Kavling
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		services.KreditInput	true	"Data transaksi kredit"
//	@Success		201		{object} object
//	@Failure		400		{object} object
//	@Router			/transaksi/kredit [post]
func (h *TransaksiHandler) CreateKredit(c *gin.Context) {
	var input services.KreditInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.BadRequest(c, "Data tidak lengkap", err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(int)
	input.IDUser = &uid

	trx, err := h.service.CreateKredit(input)
	if err != nil {
		helpers.BadRequest(c, err.Error(), nil)
		return
	}

	helpers.Created(c, "Transaksi kredit berhasil dibuat", trx)
}

// ─── Pembayaran Handler ───────────────────────────────────────────────────────

type PembayaranHandler struct {
	repo    *repositories.PembayaranRepository
	trxRepo *repositories.TransaksiKavlingRepository
	service *services.TransaksiService
}

func NewPembayaranHandler(
	repo *repositories.PembayaranRepository,
	trxRepo *repositories.TransaksiKavlingRepository,
	service *services.TransaksiService,
) *PembayaranHandler {
	return &PembayaranHandler{repo: repo, trxRepo: trxRepo, service: service}
}

// List godoc
//
//	@Summary		List semua pembayaran cicilan
//	@Description	Mengambil semua record pembayaran cicilan beserta info kavling dan customer
//	@Tags			Pembayaran
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object} object
//	@Router			/pembayaran [get]
func (h *PembayaranHandler) List(c *gin.Context) {
	list, err := h.repo.FindAll()
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil data pembayaran")
		return
	}
	helpers.OK(c, "OK", list)
}

// DetailKavling godoc
//
//	@Summary		Detail pembayaran per kavling
//	@Description	Mengambil riwayat pembayaran cicilan untuk kavling tertentu beserta data transaksinya
//	@Tags			Pembayaran
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id_kavling	path		int	true	"Kavling ID"
//	@Success		200			{object} object
//	@Failure		500			{object} object
//	@Router			/pembayaran/{id_kavling} [get]
func (h *PembayaranHandler) DetailKavling(c *gin.Context) {
	idKavling, _ := strconv.Atoi(c.Param("id_kavling"))

	list, err := h.repo.FindByKavling(idKavling)
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil data pembayaran")
		return
	}

	trx, _ := h.trxRepo.FindByKavling(idKavling)

	helpers.OK(c, "OK", gin.H{
		"transaksi":   trx,
		"pembayaran":  list,
		"total_bayar": len(list),
	})
}

// Bayar godoc
//
//	@Summary		Bayar cicilan
//	@Description	Mencatat pembayaran cicilan untuk kavling kredit. Otomatis mencatat ke arus kas. Jika sudah mencapai lama_cicilan, status kavling berubah menjadi CASH/Lunas.
//	@Tags			Pembayaran
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id_kavling			path		int		true	"Kavling ID"
//	@Param			jumlah_bayar		formData	number	true	"Jumlah pembayaran"
//	@Param			id_bank				formData	int		false	"ID rekening bank tujuan"
//	@Param			tanggal				formData	string	false	"Tanggal bayar (YYYY-MM-DD, default: hari ini)"
//	@Param			keterangan			formData	string	false	"Keterangan"
//	@Param			bukti_pembayaran	formData	file	false	"Bukti transfer (gambar)"
//	@Success		201					{object} object
//	@Failure		400					{object} object
//	@Router			/pembayaran/{id_kavling}/bayar [post]
func (h *PembayaranHandler) Bayar(c *gin.Context) {
	idKavling, _ := strconv.Atoi(c.Param("id_kavling"))

	jumlahBayarStr := c.PostForm("jumlah_bayar")
	jumlahBayar, err := strconv.ParseFloat(jumlahBayarStr, 64)
	if err != nil || jumlahBayar <= 0 {
		helpers.BadRequest(c, "Jumlah bayar tidak valid", nil)
		return
	}

	var idBank *int
	if idBankStr := c.PostForm("id_bank"); idBankStr != "" {
		id, _ := strconv.Atoi(idBankStr)
		idBank = &id
	}

	userID, _ := c.Get("user_id")
	uid := userID.(int)

	input := services.BayarInput{
		IDKavling:   idKavling,
		IDBank:      idBank,
		Tanggal:     c.PostForm("tanggal"),
		JumlahBayar: jumlahBayar,
		Keterangan:  c.PostForm("keterangan"),
		IDUser:      &uid,
	}

	uploadDir := filepath.Join(config.AppConfig.UploadPath, "bukti_trx", "bayar_angsuran")
	if file, err := c.FormFile("bukti_pembayaran"); err == nil {
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + filepath.Ext(file.Filename)
		if err := c.SaveUploadedFile(file, filepath.Join(uploadDir, filename)); err == nil {
			input.BuktiPembayaran = filename
		}
	}

	pem, err := h.service.BayarCicilan(input)
	if err != nil {
		helpers.BadRequest(c, err.Error(), nil)
		return
	}

	helpers.Created(c, "Pembayaran berhasil dicatat", pem)
}

// ─── Keuangan (Arus Kas) Handler ──────────────────────────────────────────────

type KeuanganHandler struct {
	repo    *repositories.ArusKasRepository
	service *services.TransaksiService
}

func NewKeuanganHandler(repo *repositories.ArusKasRepository, service *services.TransaksiService) *KeuanganHandler {
	return &KeuanganHandler{repo: repo, service: service}
}

// ListArusKas godoc
//
//	@Summary		List arus kas
//	@Description	Mengambil daftar transaksi arus kas dengan filter tanggal. Response juga mencakup total pemasukan, pengeluaran, dan saldo.
//	@Tags			Keuangan
//	@Produce		json
//	@Security		BearerAuth
//	@Param			dari	query		string	false	"Dari tanggal (YYYY-MM-DD)"
//	@Param			sampai	query		string	false	"Sampai tanggal (YYYY-MM-DD)"
//	@Success		200		{object} object
//	@Router			/keuangan/transaksi [get]
func (h *KeuanganHandler) ListArusKas(c *gin.Context) {
	dari := c.Query("dari")
	sampai := c.Query("sampai")

	list, err := h.repo.FindAll(dari, sampai)
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil data arus kas")
		return
	}

	pemasukan, pengeluaran := h.repo.SumByJenis(dari, sampai)

	helpers.OK(c, "OK", gin.H{
		"transaksi":   list,
		"pemasukan":   pemasukan,
		"pengeluaran": pengeluaran,
		"saldo":       pemasukan - pengeluaran,
	})
}

type arusKasRequest struct {
	Jenis      string  `json:"jenis" binding:"required,oneof=Pemasukan Pengeluaran" example:"Pengeluaran"`
	Kategori   string  `json:"kategori" binding:"required" example:"Biaya Operasional"`
	Keterangan string  `json:"keterangan" example:"Pembelian ATK"`
	Nominal    float64 `json:"nominal" binding:"required,gt=0" example:"500000"`
	IDBank     *int    `json:"id_bank" example:"1"`
	Tanggal    string  `json:"tanggal" example:"2026-04-14"`
}

// CreateArusKas godoc
//
//	@Summary		Tambah transaksi arus kas manual
//	@Description	Mencatat pemasukan atau pengeluaran secara manual (bukan dari cicilan). Saldo bank otomatis diupdate.
//	@Tags			Keuangan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		arusKasRequest	true	"Data transaksi"
//	@Success		201		{object} object
//	@Failure		400		{object} object
//	@Router			/keuangan/transaksi [post]
func (h *KeuanganHandler) CreateArusKas(c *gin.Context) {
	var req arusKasRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak lengkap", err.Error())
		return
	}

	tanggal := time.Now()
	if req.Tanggal != "" {
		t, err := time.Parse("2006-01-02", req.Tanggal)
		if err != nil {
			helpers.BadRequest(c, "Format tanggal tidak valid", nil)
			return
		}
		tanggal = t
	}

	userID, _ := c.Get("user_id")
	uid := userID.(int)

	delta := req.Nominal
	if req.Jenis == "Pengeluaran" {
		delta = -req.Nominal
	}

	t := &repositories.ArusKasEntry{
		NoTransaksi:   helpers.GenerateNoArusKas(),
		Jenis:         req.Jenis,
		Kategori:      req.Kategori,
		Keterangan:    req.Keterangan,
		Nominal:       req.Nominal,
		IDBank:        req.IDBank,
		Tanggal:       tanggal,
		ReferensiTipe: "manual",
		IDUser:        &uid,
	}

	if err := h.repo.CreateEntry(t); err != nil {
		helpers.InternalError(c, "Gagal menyimpan transaksi")
		return
	}

	if req.IDBank != nil {
		_ = h.repo.UpdateBankSaldo(*req.IDBank, delta)
	}

	helpers.Created(c, "Transaksi berhasil dicatat", t)
}

// DeleteArusKas godoc
//
//	@Summary		Hapus transaksi arus kas
//	@Description	Menghapus transaksi arus kas manual. Saldo bank otomatis dikembalikan.
//	@Tags			Keuangan
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Transaksi ID"
//	@Success		200	{object} object
//	@Failure		404	{object} object
//	@Router			/keuangan/transaksi/{id} [delete]
func (h *KeuanganHandler) DeleteArusKas(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	trx, err := h.repo.FindByID(id)
	if err != nil {
		helpers.NotFound(c, "Transaksi tidak ditemukan")
		return
	}

	if trx.IDBank != nil {
		delta := -trx.Nominal
		if trx.Jenis == "Pengeluaran" {
			delta = trx.Nominal
		}
		_ = h.repo.UpdateBankSaldo(*trx.IDBank, delta)
	}

	if err := h.repo.Delete(id); err != nil {
		helpers.InternalError(c, "Gagal menghapus transaksi")
		return
	}

	helpers.OK(c, "Transaksi berhasil dihapus", nil)
}

// RekapKredit godoc
//
//	@Summary		Rekap kredit semua kavling
//	@Description	Menghitung rekap cicilan semua kavling kredit: bulan berjalan, jumlah pembayaran, tunggakan, nominal tunggakan, dan status bulan ini.
//	@Tags			Keuangan
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object} object
//	@Router			/keuangan/rekap-kredit [get]
func (h *KeuanganHandler) RekapKredit(c *gin.Context) {
	rekap, err := h.service.GetRekapKredit()
	if err != nil {
		helpers.InternalError(c, "Gagal menghitung rekap kredit")
		return
	}
	helpers.OK(c, "OK", rekap)
}

// Statistik godoc
//
//	@Summary		Statistik arus kas per bulan
//	@Description	Mengambil rekap pemasukan dan pengeluaran per bulan untuk tahun tertentu
//	@Tags			Keuangan
//	@Produce		json
//	@Security		BearerAuth
//	@Param			year	query		int	false	"Tahun (default: tahun sekarang)"
//	@Success		200		{object} object
//	@Router			/keuangan/statistik [get]
func (h *KeuanganHandler) Statistik(c *gin.Context) {
	year := time.Now().Year()
	if y := c.Query("year"); y != "" {
		if v, err := strconv.Atoi(y); err == nil {
			year = v
		}
	}

	data, err := h.repo.GetStatistik(year)
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil statistik")
		return
	}
	helpers.OK(c, "OK", data)
}
