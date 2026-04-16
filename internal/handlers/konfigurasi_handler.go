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

type KonfigurasiHandler struct {
	repo *repositories.KonfigurasiRepository
}

func NewKonfigurasiHandler(repo *repositories.KonfigurasiRepository) *KonfigurasiHandler {
	return &KonfigurasiHandler{repo: repo}
}

// Get godoc
//
//	@Summary		Ambil konfigurasi sistem
//	@Description	Mengambil data konfigurasi umum sistem (nama perusahaan, alamat, logo, tanda tangan, dll)
//	@Tags			Konfigurasi
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object} object
//	@Failure		404	{object} object
//	@Router			/konfigurasi [get]
func (h *KonfigurasiHandler) Get(c *gin.Context) {
	k, err := h.repo.Get()
	if err != nil {
		helpers.NotFound(c, "Konfigurasi belum diatur")
		return
	}
	helpers.OK(c, "OK", k)
}

// Update godoc
//
//	@Summary		Update konfigurasi sistem
//	@Description	Mengubah data konfigurasi umum. Mendukung upload logo dan tanda tangan digital (multipart/form-data).
//	@Tags			Konfigurasi
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			nama_perusahaan	formData	string	false	"Nama perusahaan"
//	@Param			alamat			formData	string	false	"Alamat"
//	@Param			telepon			formData	string	false	"Telepon"
//	@Param			email			formData	string	false	"Email"
//	@Param			website			formData	string	false	"Website"
//	@Param			nama_ttd		formData	string	false	"Nama penanda tangan"
//	@Param			jabatan_ttd		formData	string	false	"Jabatan penanda tangan"
//	@Param			logo			formData	file	false	"File logo perusahaan"
//	@Param			ttd_digital		formData	file	false	"File tanda tangan digital"
//	@Success		200				{object} object
//	@Failure		404				{object} object
//	@Router			/konfigurasi [put]
func (h *KonfigurasiHandler) Update(c *gin.Context) {
	k, err := h.repo.Get()
	if err != nil {
		helpers.NotFound(c, "Konfigurasi tidak ditemukan")
		return
	}

	if v := c.PostForm("nama_perusahaan"); v != "" {
		k.NamaPerusahaan = v
	}
	if v := c.PostForm("alamat"); v != "" {
		k.Alamat = v
	}
	if v := c.PostForm("telepon"); v != "" {
		k.Telepon = v
	}
	if v := c.PostForm("email"); v != "" {
		k.Email = v
	}
	if v := c.PostForm("website"); v != "" {
		k.Website = v
	}
	if v := c.PostForm("nama_ttd"); v != "" {
		k.NamaTtd = v
	}
	if v := c.PostForm("jabatan_ttd"); v != "" {
		k.JabatanTtd = v
	}

	uploadDir := filepath.Join(config.AppConfig.UploadPath, "konfigurasi")
	if file, err := c.FormFile("logo"); err == nil {
		filename := "logo" + filepath.Ext(file.Filename)
		if err := c.SaveUploadedFile(file, filepath.Join(uploadDir, filename)); err == nil {
			k.Logo = filename
		}
	}
	if file, err := c.FormFile("ttd_digital"); err == nil {
		filename := "ttd_digital" + filepath.Ext(file.Filename)
		if err := c.SaveUploadedFile(file, filepath.Join(uploadDir, filename)); err == nil {
			k.TtdDigital = filename
		}
	}

	k.UpdatedAt = time.Now()
	if err := h.repo.Update(k); err != nil {
		helpers.InternalError(c, "Gagal update konfigurasi")
		return
	}

	helpers.OK(c, "Konfigurasi berhasil diupdate", k)
}

// GetWA godoc
//
//	@Summary		Ambil konfigurasi WhatsApp
//	@Description	Mengambil konfigurasi gateway WhatsApp yang aktif
//	@Tags			Konfigurasi
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object} object
//	@Failure		404	{object} object
//	@Router			/konfigurasi/wa [get]
func (h *KonfigurasiHandler) GetWA(c *gin.Context) {
	wa, err := h.repo.GetWA()
	if err != nil {
		helpers.NotFound(c, "Konfigurasi WhatsApp belum diatur")
		return
	}
	helpers.OK(c, "OK", wa)
}

type waRequest struct {
	IDDevice string `json:"id_device" binding:"required" example:"device-001"`
	NoTelp   string `json:"no_telp" binding:"required" example:"628123456789"`
	ApiURL   string `json:"api_url" binding:"required" example:"https://wa.api.example.com"`
	ApiKey   string `json:"api_key" example:"secret-key"`
}

// UpdateWA godoc
//
//	@Summary		Update konfigurasi WhatsApp
//	@Description	Mengubah atau membuat konfigurasi gateway WhatsApp
//	@Tags			Konfigurasi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		waRequest	true	"Konfigurasi WhatsApp"
//	@Success		200		{object} object
//	@Failure		400		{object} object
//	@Router			/konfigurasi/wa [put]
func (h *KonfigurasiHandler) UpdateWA(c *gin.Context) {
	var req waRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak lengkap", err.Error())
		return
	}

	wa, err := h.repo.GetWA()
	if err != nil {
		wa = &models.KonfigurasiWA{
			IDDevice: req.IDDevice,
			NoTelp:   req.NoTelp,
			ApiURL:   req.ApiURL,
			ApiKey:   req.ApiKey,
			IsAktif:  1,
		}
	} else {
		wa.IDDevice = req.IDDevice
		wa.NoTelp = req.NoTelp
		wa.ApiURL = req.ApiURL
		wa.ApiKey = req.ApiKey
		wa.UpdatedAt = time.Now()
	}

	if err := h.repo.UpsertWA(wa); err != nil {
		helpers.InternalError(c, "Gagal update konfigurasi WhatsApp")
		return
	}

	helpers.OK(c, "Konfigurasi WhatsApp berhasil diupdate", wa)
}

// ListTemplate godoc
//
//	@Summary		List template pesan WhatsApp
//	@Description	Mengambil daftar semua template pesan WhatsApp
//	@Tags			Konfigurasi
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object} object
//	@Router			/konfigurasi/template [get]
func (h *KonfigurasiHandler) ListTemplate(c *gin.Context) {
	list, err := h.repo.FindAllTemplate()
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil template")
		return
	}
	helpers.OK(c, "OK", list)
}

type templateRequest struct {
	Nama    string `json:"nama" binding:"required" example:"Kwitansi Pembayaran"`
	Tipe    string `json:"tipe" binding:"required" example:"kwitansi"`
	Isi     string `json:"isi" binding:"required" example:"Halo {nama_customer}, pembayaran ke-{ke} sebesar Rp {nominal} telah kami terima."`
	IsAktif int    `json:"is_aktif" example:"1"`
}

// CreateTemplate godoc
//
//	@Summary		Tambah template pesan WhatsApp
//	@Description	Membuat template pesan WhatsApp baru. Variabel template menggunakan format {nama_variabel}.
//	@Tags			Konfigurasi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		templateRequest	true	"Data template"
//	@Success		201		{object} object
//	@Failure		400		{object} object
//	@Router			/konfigurasi/template [post]
func (h *KonfigurasiHandler) CreateTemplate(c *gin.Context) {
	var req templateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak lengkap", err.Error())
		return
	}

	isAktif := req.IsAktif
	if isAktif == 0 {
		isAktif = 1
	}

	t := &models.Template{
		Nama:    req.Nama,
		Tipe:    req.Tipe,
		Isi:     req.Isi,
		IsAktif: isAktif,
	}

	if err := h.repo.SaveTemplate(t); err != nil {
		helpers.InternalError(c, "Gagal menyimpan template")
		return
	}

	helpers.Created(c, "Template berhasil dibuat", t)
}

type resetRequest struct {
	Items []string `json:"items"` // pembayaran, transaksi, keuangan, kavling, customer, marketing
}

// ResetData godoc
//
//	@Summary		Reset data terpilih
//	@Description	Menghapus data berdasarkan kategori yang dipilih. Jika items kosong, semua data direset.
//	@Tags			Konfigurasi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		resetRequest	false	"Kategori yang akan direset"
//	@Success		200		{object} object
//	@Failure		500		{object} object
//	@Router			/konfigurasi/reset [post]
func (h *KonfigurasiHandler) ResetData(c *gin.Context) {
	var req resetRequest
	_ = c.ShouldBindJSON(&req)

	var err error
	if len(req.Items) == 0 {
		err = h.repo.ResetAllData()
	} else {
		err = h.repo.ResetSelected(req.Items)
	}
	if err != nil {
		helpers.InternalError(c, "Gagal mereset data: "+err.Error())
		return
	}
	helpers.OK(c, "Data berhasil direset", nil)
}

// UpdateTemplate godoc
//
//	@Summary		Update template pesan WhatsApp
//	@Description	Mengubah template pesan WhatsApp berdasarkan ID
//	@Tags			Konfigurasi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int				true	"Template ID"
//	@Param			body	body		templateRequest	true	"Data template"
//	@Success		200		{object} object
//	@Failure		400		{object} object
//	@Router			/konfigurasi/template/{id} [put]
func (h *KonfigurasiHandler) UpdateTemplate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req templateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	t := &models.Template{
		ID:      id,
		Nama:    req.Nama,
		Tipe:    req.Tipe,
		Isi:     req.Isi,
		IsAktif: req.IsAktif,
	}

	if err := h.repo.SaveTemplate(t); err != nil {
		helpers.InternalError(c, "Gagal update template")
		return
	}

	helpers.OK(c, "Template berhasil diupdate", t)
}
