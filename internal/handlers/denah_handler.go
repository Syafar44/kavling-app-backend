package handlers

import (
	"strconv"

	"backend-kavling/internal/helpers"
	"backend-kavling/internal/repositories"
	"backend-kavling/internal/services"

	"github.com/gin-gonic/gin"
)

type DenahHandler struct {
	repo    *repositories.DenahRepository
	service *services.DenahService
}

func NewDenahHandler(repo *repositories.DenahRepository, service *services.DenahService) *DenahHandler {
	return &DenahHandler{repo: repo, service: service}
}

type createDenahRequest struct {
	Nama       string `json:"nama" binding:"required" example:"Blok A"`
	SvgContent string `json:"svg_content" binding:"required" example:"<svg viewBox='0 0 800 600'>...</svg>"`
}

type updateDenahRequest struct {
	Nama       string `json:"nama" binding:"required" example:"Blok A"`
	SvgContent string `json:"svg_content" example:"<svg viewBox='0 0 800 600'>...</svg>"`
}

// ListDenah godoc
//
//	@Summary		List semua denah kavling
//	@Description	Mengambil daftar semua denah kavling beserta ringkasan jumlah kavling
//	@Tags			Denah Kavling
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object} object
//	@Router			/denah-kavling [get]
func (h *DenahHandler) List(c *gin.Context) {
	list, err := h.repo.FindAllWithSummary()
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil data denah kavling")
		return
	}
	helpers.OK(c, "OK", list)
}

// DetailDenah godoc
//
//	@Summary		Detail denah kavling
//	@Description	Mengambil detail satu denah kavling beserta seluruh kavling di dalamnya
//	@Tags			Denah Kavling
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Denah ID"
//	@Success		200	{object} object
//	@Failure		404	{object} object
//	@Router			/denah-kavling/{id} [get]
func (h *DenahHandler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	denah, err := h.repo.FindByID(id)
	if err != nil {
		helpers.NotFound(c, "Denah kavling tidak ditemukan")
		return
	}
	helpers.OK(c, "OK", denah)
}

// CreateDenah godoc
//
//	@Summary		Buat denah kavling baru dari SVG
//	@Description	Menerima nama dan kode SVG. Sistem mem-parse SVG, mengekstrak semua <path id="..."> sebagai kavling, lalu menyimpan denah beserta kavling-kavlingnya secara otomatis.
//	@Tags			Denah Kavling
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		createDenahRequest	true	"Data denah kavling"
//	@Success		201		{object} object
//	@Failure		400		{object} object
//	@Router			/denah-kavling [post]
func (h *DenahHandler) Create(c *gin.Context) {
	var req createDenahRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak lengkap", err.Error())
		return
	}

	denah, err := h.service.CreateDenah(req.Nama, req.SvgContent)
	if err != nil {
		helpers.BadRequest(c, err.Error(), nil)
		return
	}

	helpers.Created(c, "Denah kavling berhasil dibuat", denah)
}

// UpdateDenah godoc
//
//	@Summary		Update denah kavling
//	@Description	Mengubah nama dan/atau SVG content denah. Jika svg_content diisi, viewbox akan diperbarui. Kavling yang sudah ada tidak diubah.
//	@Tags			Denah Kavling
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int					true	"Denah ID"
//	@Param			body	body		updateDenahRequest	true	"Data update denah"
//	@Success		200		{object} object
//	@Failure		400		{object} object
//	@Failure		404		{object} object
//	@Router			/denah-kavling/{id} [put]
func (h *DenahHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req updateDenahRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	denah, err := h.service.UpdateSVG(id, req.Nama, req.SvgContent)
	if err != nil {
		if err.Error() == "denah tidak ditemukan" {
			helpers.NotFound(c, err.Error())
			return
		}
		helpers.BadRequest(c, err.Error(), nil)
		return
	}

	helpers.OK(c, "Denah kavling berhasil diupdate", denah)
}

// DeleteDenah godoc
//
//	@Summary		Hapus denah kavling
//	@Description	Menghapus denah beserta semua kavling di dalamnya. Gagal jika ada kavling dengan transaksi aktif.
//	@Tags			Denah Kavling
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Denah ID"
//	@Success		200	{object} object
//	@Failure		400	{object} object	"Ada transaksi aktif"
//	@Failure		404	{object} object
//	@Router			/denah-kavling/{id} [delete]
func (h *DenahHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// Check exists first
	if _, err := h.repo.FindByID(id); err != nil {
		helpers.NotFound(c, "Denah kavling tidak ditemukan")
		return
	}

	if err := h.service.DeleteDenah(id); err != nil {
		helpers.BadRequest(c, err.Error(), nil)
		return
	}

	helpers.OK(c, "Denah kavling berhasil dihapus", nil)
}
