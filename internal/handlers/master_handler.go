package handlers

import (
	"path/filepath"
	"strconv"
	"time"

	"backend-kavling/internal/config"
	"backend-kavling/internal/helpers"
	"backend-kavling/internal/models"
	"backend-kavling/internal/repositories"

	"github.com/gin-gonic/gin"
)

// ─── Marketing Handler ────────────────────────────────────────────────────────

type MarketingHandler struct {
	repo *repositories.MarketingRepository
}

func NewMarketingHandler(repo *repositories.MarketingRepository) *MarketingHandler {
	return &MarketingHandler{repo: repo}
}

// List godoc
//
//	@Summary		List marketing
//	@Description	Mengambil daftar semua marketing beserta jumlah transaksi OPEN dan CLOSED
//	@Tags			Marketing
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object} object
//	@Router			/marketing [get]
func (h *MarketingHandler) List(c *gin.Context) {
	list, err := h.repo.FindAll()
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil data marketing")
		return
	}
	for i := range list {
		open, closed := h.repo.CountOpenClosed(list[i].ID)
		list[i].JumlahOpen = open
		list[i].JumlahClosed = closed
	}
	helpers.OK(c, "OK", list)
}

// Detail godoc
//
//	@Summary		Detail marketing
//	@Description	Mengambil detail satu marketing beserta jumlah transaksi OPEN dan CLOSED
//	@Tags			Marketing
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Marketing ID"
//	@Success		200	{object} object
//	@Failure		404	{object} object
//	@Router			/marketing/{id} [get]
func (h *MarketingHandler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	m, err := h.repo.FindByID(id)
	if err != nil {
		helpers.NotFound(c, "Marketing tidak ditemukan")
		return
	}
	open, closed := h.repo.CountOpenClosed(id)
	m.JumlahOpen = open
	m.JumlahClosed = closed
	helpers.OK(c, "OK", m)
}

type marketingRequest struct {
	Nama             string  `json:"nama" binding:"required" example:"Budi Santoso"`
	NoTelp           string  `json:"no_telp" example:"08123456789"`
	Alamat           string  `json:"alamat" example:"Jl. Contoh No. 1, Sampit"`
	Email            string  `json:"email" example:"budi@email.com"`
	PersentaseKomisi float64 `json:"persentase_komisi" example:"2.5"`
	Status           int     `json:"status" example:"1"`
}

// Create godoc
//
//	@Summary		Tambah marketing
//	@Description	Menambahkan data marketing (sales) baru
//	@Tags			Marketing
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		marketingRequest	true	"Data marketing"
//	@Success		201		{object} object
//	@Failure		400		{object} object
//	@Router			/marketing [post]
func (h *MarketingHandler) Create(c *gin.Context) {
	var req marketingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak lengkap", err.Error())
		return
	}

	status := req.Status
	if status == 0 {
		status = 1
	}

	m := &models.Marketing{
		Nama:             req.Nama,
		NoTelp:           req.NoTelp,
		Alamat:           req.Alamat,
		Email:            req.Email,
		PersentaseKomisi: req.PersentaseKomisi,
		Status:           status,
	}

	if err := h.repo.Create(m); err != nil {
		helpers.InternalError(c, "Gagal menyimpan marketing")
		return
	}

	helpers.Created(c, "Marketing berhasil dibuat", m)
}

// Update godoc
//
//	@Summary		Update marketing
//	@Description	Mengubah data marketing berdasarkan ID
//	@Tags			Marketing
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int					true	"Marketing ID"
//	@Param			body	body		marketingRequest	true	"Data update marketing"
//	@Success		200		{object} object
//	@Failure		404		{object} object
//	@Router			/marketing/{id} [put]
func (h *MarketingHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	m, err := h.repo.FindByID(id)
	if err != nil {
		helpers.NotFound(c, "Marketing tidak ditemukan")
		return
	}

	var req marketingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	m.Nama = req.Nama
	m.NoTelp = req.NoTelp
	m.Alamat = req.Alamat
	m.Email = req.Email
	m.PersentaseKomisi = req.PersentaseKomisi
	if req.Status > 0 {
		m.Status = req.Status
	}
	m.UpdatedAt = time.Now()

	if err := h.repo.Update(m); err != nil {
		helpers.InternalError(c, "Gagal update marketing")
		return
	}

	helpers.OK(c, "Marketing berhasil diupdate", m)
}

// Delete godoc
//
//	@Summary		Hapus marketing
//	@Description	Menghapus marketing. Ditolak jika marketing sudah memiliki transaksi.
//	@Tags			Marketing
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Marketing ID"
//	@Success		200	{object} object
//	@Failure		400	{object} object	"Sudah ada transaksi"
//	@Failure		404	{object} object
//	@Router			/marketing/{id} [delete]
func (h *MarketingHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if _, err := h.repo.FindByID(id); err != nil {
		helpers.NotFound(c, "Marketing tidak ditemukan")
		return
	}

	if h.repo.HasTransaksi(id) {
		helpers.BadRequest(c, "Marketing sudah memiliki transaksi, tidak bisa dihapus", nil)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		helpers.InternalError(c, "Gagal menghapus marketing")
		return
	}

	helpers.OK(c, "Marketing berhasil dihapus", nil)
}

// ─── Customer Handler ─────────────────────────────────────────────────────────

type CustomerHandler struct {
	repo *repositories.CustomerRepository
}

func NewCustomerHandler(repo *repositories.CustomerRepository) *CustomerHandler {
	return &CustomerHandler{repo: repo}
}

// List godoc
//
//	@Summary		List customer
//	@Description	Mengambil daftar semua customer (pembeli)
//	@Tags			Customer
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object} object
//	@Router			/customers [get]
func (h *CustomerHandler) List(c *gin.Context) {
	list, err := h.repo.FindAll()
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil data customer")
		return
	}
	helpers.OK(c, "OK", list)
}

// Detail godoc
//
//	@Summary		Detail customer
//	@Description	Mengambil detail satu customer berdasarkan ID
//	@Tags			Customer
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Customer ID"
//	@Success		200	{object} object
//	@Failure		404	{object} object
//	@Router			/customers/{id} [get]
func (h *CustomerHandler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	cust, err := h.repo.FindByID(id)
	if err != nil {
		helpers.NotFound(c, "Customer tidak ditemukan")
		return
	}
	helpers.OK(c, "OK", cust)
}

// Create godoc
//
//	@Summary		Tambah customer
//	@Description	Membuat data customer baru. Mendukung upload file foto KTP dan KK (multipart/form-data).
//	@Tags			Customer
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			nama		formData	string	true	"Nama lengkap"
//	@Param			no_telp		formData	string	false	"No. telepon"
//	@Param			no_ktp		formData	string	false	"Nomor KTP"
//	@Param			alamat		formData	string	false	"Alamat"
//	@Param			pekerjaan	formData	string	false	"Pekerjaan"
//	@Param			foto_ktp	formData	file	false	"Foto KTP (jpg/png)"
//	@Param			foto_kk		formData	file	false	"Foto KK (jpg/png)"
//	@Success		201			{object} object
//	@Failure		400			{object} object
//	@Router			/customers [post]
func (h *CustomerHandler) Create(c *gin.Context) {
	nama := c.PostForm("nama")
	if nama == "" {
		helpers.BadRequest(c, "Nama customer wajib diisi", nil)
		return
	}

	cust := &models.Customer{
		Nama:      nama,
		NoTelp:    c.PostForm("no_telp"),
		NoKTP:     c.PostForm("no_ktp"),
		Alamat:    c.PostForm("alamat"),
		Pekerjaan: c.PostForm("pekerjaan"),
	}

	uploadDir := filepath.Join(config.AppConfig.UploadPath, "lampiran_customer")
	if ktpFile, err := c.FormFile("foto_ktp"); err == nil {
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_ktp" + filepath.Ext(ktpFile.Filename)
		if err := c.SaveUploadedFile(ktpFile, filepath.Join(uploadDir, filename)); err == nil {
			cust.FotoKTP = filename
		}
	}
	if kkFile, err := c.FormFile("foto_kk"); err == nil {
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_kk" + filepath.Ext(kkFile.Filename)
		if err := c.SaveUploadedFile(kkFile, filepath.Join(uploadDir, filename)); err == nil {
			cust.FotoKK = filename
		}
	}

	if err := h.repo.Create(cust); err != nil {
		helpers.InternalError(c, "Gagal menyimpan customer")
		return
	}

	helpers.Created(c, "Customer berhasil dibuat", cust)
}

// Update godoc
//
//	@Summary		Update customer
//	@Description	Mengubah data customer. Mendukung re-upload foto KTP/KK.
//	@Tags			Customer
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			path		int		true	"Customer ID"
//	@Param			nama		formData	string	false	"Nama lengkap"
//	@Param			no_telp		formData	string	false	"No. telepon"
//	@Param			no_ktp		formData	string	false	"Nomor KTP"
//	@Param			alamat		formData	string	false	"Alamat"
//	@Param			pekerjaan	formData	string	false	"Pekerjaan"
//	@Param			foto_ktp	formData	file	false	"Foto KTP baru"
//	@Param			foto_kk		formData	file	false	"Foto KK baru"
//	@Success		200			{object} object
//	@Failure		404			{object} object
//	@Router			/customers/{id} [put]
func (h *CustomerHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	cust, err := h.repo.FindByID(id)
	if err != nil {
		helpers.NotFound(c, "Customer tidak ditemukan")
		return
	}

	if nama := c.PostForm("nama"); nama != "" {
		cust.Nama = nama
	}
	if v := c.PostForm("no_telp"); v != "" {
		cust.NoTelp = v
	}
	if v := c.PostForm("no_ktp"); v != "" {
		cust.NoKTP = v
	}
	if v := c.PostForm("alamat"); v != "" {
		cust.Alamat = v
	}
	if v := c.PostForm("pekerjaan"); v != "" {
		cust.Pekerjaan = v
	}

	uploadDir := filepath.Join(config.AppConfig.UploadPath, "lampiran_customer")
	if ktpFile, err := c.FormFile("foto_ktp"); err == nil {
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_ktp" + filepath.Ext(ktpFile.Filename)
		if err := c.SaveUploadedFile(ktpFile, filepath.Join(uploadDir, filename)); err == nil {
			cust.FotoKTP = filename
		}
	}
	if kkFile, err := c.FormFile("foto_kk"); err == nil {
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_kk" + filepath.Ext(kkFile.Filename)
		if err := c.SaveUploadedFile(kkFile, filepath.Join(uploadDir, filename)); err == nil {
			cust.FotoKK = filename
		}
	}

	cust.UpdatedAt = time.Now()
	if err := h.repo.Update(cust); err != nil {
		helpers.InternalError(c, "Gagal update customer")
		return
	}

	helpers.OK(c, "Customer berhasil diupdate", cust)
}

// Delete godoc
//
//	@Summary		Hapus customer
//	@Description	Menghapus customer berdasarkan ID
//	@Tags			Customer
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Customer ID"
//	@Success		200	{object} object
//	@Failure		404	{object} object
//	@Router			/customers/{id} [delete]
func (h *CustomerHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if _, err := h.repo.FindByID(id); err != nil {
		helpers.NotFound(c, "Customer tidak ditemukan")
		return
	}

	if err := h.repo.Delete(id); err != nil {
		helpers.InternalError(c, "Gagal menghapus customer")
		return
	}

	helpers.OK(c, "Customer berhasil dihapus", nil)
}

// ─── Bank Handler ─────────────────────────────────────────────────────────────

type BankHandler struct {
	repo *repositories.BankRepository
}

func NewBankHandler(repo *repositories.BankRepository) *BankHandler {
	return &BankHandler{repo: repo}
}

// List godoc
//
//	@Summary		List rekening bank/kas
//	@Description	Mengambil daftar semua rekening bank dan kas yang aktif
//	@Tags			Bank
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object} object
//	@Router			/bank [get]
func (h *BankHandler) List(c *gin.Context) {
	list, err := h.repo.FindAll()
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil data bank")
		return
	}
	helpers.OK(c, "OK", list)
}

type bankRequest struct {
	NamaBank     string `json:"nama_bank" binding:"required" example:"BRI"`
	NamaRekening string `json:"nama_rekening" binding:"required" example:"PT Kavling Mentaya"`
	NoRekening   string `json:"no_rekening" binding:"required" example:"1234567890"`
	IsKas        int    `json:"is_kas" example:"0"`
}

// Create godoc
//
//	@Summary		Tambah rekening bank/kas
//	@Description	Menambahkan rekening bank atau kas tunai baru
//	@Tags			Bank
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		bankRequest	true	"Data bank"
//	@Success		201		{object} object
//	@Failure		400		{object} object
//	@Router			/bank [post]
func (h *BankHandler) Create(c *gin.Context) {
	var req bankRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak lengkap", err.Error())
		return
	}

	b := &models.Bank{
		NamaBank:     req.NamaBank,
		NamaRekening: req.NamaRekening,
		NoRekening:   req.NoRekening,
		IsKas:        req.IsKas,
		Status:       1,
	}

	if err := h.repo.Create(b); err != nil {
		helpers.InternalError(c, "Gagal menyimpan bank")
		return
	}

	helpers.Created(c, "Bank berhasil ditambahkan", b)
}

// Update godoc
//
//	@Summary		Update rekening bank/kas
//	@Description	Mengubah data rekening bank berdasarkan ID
//	@Tags			Bank
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int			true	"Bank ID"
//	@Param			body	body		bankRequest	true	"Data update bank"
//	@Success		200		{object} object
//	@Failure		404		{object} object
//	@Router			/bank/{id} [put]
func (h *BankHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	b, err := h.repo.FindByID(id)
	if err != nil {
		helpers.NotFound(c, "Bank tidak ditemukan")
		return
	}

	var req bankRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	b.NamaBank = req.NamaBank
	b.NamaRekening = req.NamaRekening
	b.NoRekening = req.NoRekening
	b.IsKas = req.IsKas
	b.UpdatedAt = time.Now()

	if err := h.repo.Update(b); err != nil {
		helpers.InternalError(c, "Gagal update bank")
		return
	}

	helpers.OK(c, "Bank berhasil diupdate", b)
}

// Delete godoc
//
//	@Summary		Hapus rekening bank/kas
//	@Description	Menghapus rekening bank berdasarkan ID
//	@Tags			Bank
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Bank ID"
//	@Success		200	{object} object
//	@Failure		404	{object} object
//	@Router			/bank/{id} [delete]
func (h *BankHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if _, err := h.repo.FindByID(id); err != nil {
		helpers.NotFound(c, "Bank tidak ditemukan")
		return
	}
	if err := h.repo.Delete(id); err != nil {
		helpers.InternalError(c, "Gagal menghapus bank")
		return
	}
	helpers.OK(c, "Bank berhasil dihapus", nil)
}
