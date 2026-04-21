package handlers

import (
	"net/http"
	"strconv"
	"time"

	"backend-kavling/internal/models"
	"backend-kavling/internal/repositories"

	"github.com/gin-gonic/gin"
)

type KeuanganV2Handler struct {
	repo         *repositories.KeuanganRepository
	kategoriRepo *repositories.KategoriRepository
}

func NewKeuanganV2Handler(r *repositories.KeuanganRepository, k *repositories.KategoriRepository) *KeuanganV2Handler {
	return &KeuanganV2Handler{r, k}
}

// ─── Kategori Transaksi ───────────────────────────────────────────────────────

func (h *KeuanganV2Handler) ListKategori(c *gin.Context) {
	jenis := c.Query("jenis")
	list, err := h.kategoriRepo.List(jenis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: list})
}

func (h *KeuanganV2Handler) CreateKategori(c *gin.Context) {
	var m models.KategoriTransaksi
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}
	if m.Kode == "" {
		kode, _ := h.kategoriRepo.NextKode()
		m.Kode = kode
	}
	if err := h.kategoriRepo.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, models.Response{Success: true, Data: m})
}

func (h *KeuanganV2Handler) UpdateKategori(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	existing, err := h.kategoriRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "tidak ditemukan"})
		return
	}
	if err := c.ShouldBindJSON(existing); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}
	existing.ID = id
	if err := h.kategoriRepo.Update(existing); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: existing})
}

func (h *KeuanganV2Handler) DeleteKategori(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.kategoriRepo.Delete(id); err != nil {
		c.JSON(http.StatusUnprocessableEntity, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "berhasil dihapus"})
}

// ─── Pemasukan ───────────────────────────────────────────────────────────────

func (h *KeuanganV2Handler) ListPemasukan(c *gin.Context) {
	q := c.Query("q")
	tahun, _ := strconv.Atoi(c.Query("tahun"))
	bulan, _ := strconv.Atoi(c.Query("bulan"))
	idBank, _ := strconv.Atoi(c.Query("id_bank"))
	idKat, _ := strconv.Atoi(c.Query("id_kategori"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	list, total, err := h.repo.ListTransaksi("Pemasukan", q, tahun, bulan, idBank, idKat, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.PaginatedResponse{Success: true, Data: list, Total: total, Page: page, PerPage: perPage})
}

func (h *KeuanganV2Handler) CreatePemasukan(c *gin.Context) {
	var req struct {
		Tanggal    string  `json:"tanggal" binding:"required"`
		IDKategori int     `json:"id_kategori" binding:"required"`
		IDBank     *int    `json:"id_bank"`
		Nominal    float64 `json:"nominal" binding:"required,gt=0"`
		Keterangan string  `json:"keterangan"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}
	kat, err := h.kategoriRepo.FindByID(req.IDKategori)
	if err != nil || kat.Jenis != "PEMASUKAN" {
		c.JSON(http.StatusBadRequest, models.Response{Message: "kategori bukan pemasukan"})
		return
	}
	tgl, _ := time.Parse("2006-01-02", req.Tanggal)
	m := &models.Transaksi{
		NoTransaksi:   h.repo.GenerateNoTransaksi("IN"),
		Jenis:         "Pemasukan",
		Kategori:      kat.Kategori,
		Nominal:       req.Nominal,
		IDBank:        req.IDBank,
		Tanggal:       tgl,
		Keterangan:    req.Keterangan,
		ReferensiTipe: "manual",
	}
	if err := h.repo.CreateTransaksi(m); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	if req.IDBank != nil {
		h.repo.AddSaldo(*req.IDBank, req.Nominal)
	}
	c.JSON(http.StatusCreated, models.Response{Success: true, Data: m})
}

func (h *KeuanganV2Handler) DeletePemasukan(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	trx, err := h.repo.FindTransaksiByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "tidak ditemukan"})
		return
	}
	if trx.IDBank != nil {
		h.repo.SubtractSaldo(*trx.IDBank, trx.Nominal)
	}
	h.repo.DeleteTransaksi(id)
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "dihapus"})
}

// ─── Pengeluaran ──────────────────────────────────────────────────────────────

func (h *KeuanganV2Handler) ListPengeluaran(c *gin.Context) {
	q := c.Query("q")
	tahun, _ := strconv.Atoi(c.Query("tahun"))
	bulan, _ := strconv.Atoi(c.Query("bulan"))
	idBank, _ := strconv.Atoi(c.Query("id_bank"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	list, total, err := h.repo.ListTransaksi("Pengeluaran", q, tahun, bulan, idBank, 0, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.PaginatedResponse{Success: true, Data: list, Total: total, Page: page, PerPage: perPage})
}

func (h *KeuanganV2Handler) CreatePengeluaran(c *gin.Context) {
	var req struct {
		Tanggal    string  `json:"tanggal" binding:"required"`
		IDKategori int     `json:"id_kategori" binding:"required"`
		IDBank     *int    `json:"id_bank"`
		Nominal    float64 `json:"nominal" binding:"required,gt=0"`
		Keterangan string  `json:"keterangan"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}
	kat, err := h.kategoriRepo.FindByID(req.IDKategori)
	if err != nil || kat.Jenis != "PENGELUARAN" {
		c.JSON(http.StatusBadRequest, models.Response{Message: "kategori bukan pengeluaran"})
		return
	}
	tgl, _ := time.Parse("2006-01-02", req.Tanggal)
	m := &models.Transaksi{
		NoTransaksi:   h.repo.GenerateNoTransaksi("OUT"),
		Jenis:         "Pengeluaran",
		Kategori:      kat.Kategori,
		Nominal:       req.Nominal,
		IDBank:        req.IDBank,
		Tanggal:       tgl,
		Keterangan:    req.Keterangan,
		ReferensiTipe: "manual",
	}
	if err := h.repo.CreateTransaksi(m); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	if req.IDBank != nil {
		h.repo.SubtractSaldo(*req.IDBank, req.Nominal)
	}
	c.JSON(http.StatusCreated, models.Response{Success: true, Data: m})
}

func (h *KeuanganV2Handler) DeletePengeluaran(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	trx, err := h.repo.FindTransaksiByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "tidak ditemukan"})
		return
	}
	if trx.IDBank != nil {
		h.repo.AddSaldo(*trx.IDBank, trx.Nominal)
	}
	h.repo.DeleteTransaksi(id)
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "dihapus"})
}

// ─── Hutang ──────────────────────────────────────────────────────────────────

func (h *KeuanganV2Handler) ListHutang(c *gin.Context) {
	list, err := h.repo.ListHutang(c.Query("q"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: list})
}

func (h *KeuanganV2Handler) CreateHutang(c *gin.Context) {
	var m models.Hutang
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}
	if err := h.repo.CreateHutang(&m); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, models.Response{Success: true, Data: m})
}

func (h *KeuanganV2Handler) UpdateHutang(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	existing, err := h.repo.FindHutangByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "tidak ditemukan"})
		return
	}
	c.ShouldBindJSON(existing)
	existing.ID = id
	h.repo.UpdateHutang(existing)
	c.JSON(http.StatusOK, models.Response{Success: true, Data: existing})
}

func (h *KeuanganV2Handler) DeleteHutang(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.repo.DeleteHutang(id)
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "dihapus"})
}

func (h *KeuanganV2Handler) BayarHutang(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Nominal float64 `json:"nominal" binding:"required,gt=0"`
		IDBank  *int    `json:"id_bank"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}
	if err := h.repo.BayarHutang(id, req.Nominal); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}

	// Ambil data hutang untuk mendapatkan id_bank bawaan jika tidak dikirim
	hutang, err := h.repo.FindHutangByID(id)
	if err == nil && req.IDBank == nil {
		req.IDBank = hutang.IDBank
	}
	// Catat ke arus kas
	tgl := time.Now()
	trx := &models.Transaksi{
		NoTransaksi:   h.repo.GenerateNoTransaksi("OUT"),
		Jenis:         "Pengeluaran",
		Kategori:      "Pembayaran Hutang",
		Nominal:       req.Nominal,
		IDBank:        req.IDBank,
		Tanggal:       tgl,
		ReferensiTipe: "hutang",
		ReferensiID:   &id,
	}
	h.repo.CreateTransaksi(trx)
	if req.IDBank != nil {
		h.repo.SubtractSaldo(*req.IDBank, req.Nominal)
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "hutang dibayar"})
}

// ─── Piutang ─────────────────────────────────────────────────────────────────

func (h *KeuanganV2Handler) ListPiutang(c *gin.Context) {
	list, err := h.repo.ListPiutang(c.Query("q"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: list})
}

func (h *KeuanganV2Handler) CreatePiutang(c *gin.Context) {
	var m models.Piutang
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}
	if err := h.repo.CreatePiutang(&m); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, models.Response{Success: true, Data: m})
}

func (h *KeuanganV2Handler) DeletePiutang(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.repo.DeletePiutang(id)
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "dihapus"})
}

func (h *KeuanganV2Handler) BayarPiutang(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Nominal float64 `json:"nominal" binding:"required,gt=0"`
		IDBank  *int    `json:"id_bank"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}
	h.repo.BayarPiutang(id, req.Nominal)

	// Ambil data piutang untuk mendapatkan id_bank bawaan jika tidak dikirim
	piutang, err := h.repo.FindPiutangByID(id)
	if err == nil && req.IDBank == nil {
		req.IDBank = piutang.IDBank
	}
	tgl := time.Now()
	trx := &models.Transaksi{
		NoTransaksi:   h.repo.GenerateNoTransaksi("IN"),
		Jenis:         "Pemasukan",
		Kategori:      "Pencairan Piutang",
		Nominal:       req.Nominal,
		IDBank:        req.IDBank,
		Tanggal:       tgl,
		ReferensiTipe: "piutang",
		ReferensiID:   &id,
	}
	h.repo.CreateTransaksi(trx)
	if req.IDBank != nil {
		h.repo.AddSaldo(*req.IDBank, req.Nominal)
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "piutang diterima"})
}

// ─── Mutasi Saldo ─────────────────────────────────────────────────────────────

func (h *KeuanganV2Handler) ListMutasi(c *gin.Context) {
	list, err := h.repo.ListMutasi()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: list})
}

func (h *KeuanganV2Handler) CreateMutasi(c *gin.Context) {
	var m models.MutasiSaldo
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}
	if m.IDBankAsal == m.IDBankTujuan {
		c.JSON(http.StatusBadRequest, models.Response{Message: "rekening asal dan tujuan tidak boleh sama"})
		return
	}
	if err := h.repo.CreateMutasi(&m); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	// Catat 2 transaksi arus kas
	tgl := m.Tanggal
	trxKeluar := &models.Transaksi{
		NoTransaksi:   h.repo.GenerateNoTransaksi("OUT"),
		Jenis:         "Pengeluaran",
		Kategori:      "Mutasi Saldo",
		Nominal:       m.Nominal,
		IDBank:        &m.IDBankAsal,
		Tanggal:       tgl,
		ReferensiTipe: "mutasi",
		ReferensiID:   &m.ID,
		Keterangan:    m.Keterangan,
	}
	trxMasuk := &models.Transaksi{
		NoTransaksi:   h.repo.GenerateNoTransaksi("IN"),
		Jenis:         "Pemasukan",
		Kategori:      "Terima Saldo",
		Nominal:       m.Nominal,
		IDBank:        &m.IDBankTujuan,
		Tanggal:       tgl,
		ReferensiTipe: "mutasi",
		ReferensiID:   &m.ID,
		Keterangan:    m.Keterangan,
	}
	h.repo.CreateTransaksi(trxKeluar)
	h.repo.CreateTransaksi(trxMasuk)
	h.repo.SubtractSaldo(m.IDBankAsal, m.Nominal)
	h.repo.AddSaldo(m.IDBankTujuan, m.Nominal)
	c.JSON(http.StatusCreated, models.Response{Success: true, Data: m})
}

func (h *KeuanganV2Handler) DeleteMutasi(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.repo.DeleteMutasi(id)
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "dihapus"})
}

// ─── Laporan Arus Kas ─────────────────────────────────────────────────────────

func (h *KeuanganV2Handler) LaporanArusKas(c *gin.Context) {
	tahun, _ := strconv.Atoi(c.Query("tahun"))
	bulan, _ := strconv.Atoi(c.Query("bulan"))
	idBank, _ := strconv.Atoi(c.Query("id_bank"))
	idLokasi, _ := strconv.Atoi(c.Query("id_lokasi"))
	if tahun == 0 {
		c.JSON(http.StatusBadRequest, models.Response{Message: "parameter tahun wajib diisi"})
		return
	}
	list, err := h.repo.ArusKas(tahun, bulan, idBank, idLokasi)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	var totalDebit, totalKredit float64
	for _, t := range list {
		if t.Jenis == "Pemasukan" {
			totalDebit += t.Nominal
		} else {
			totalKredit += t.Nominal
		}
	}
	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data: gin.H{
			"list":         list,
			"total_debit":  totalDebit,
			"total_kredit": totalKredit,
			"saldo_akhir":  totalDebit - totalKredit,
		},
	})
}
