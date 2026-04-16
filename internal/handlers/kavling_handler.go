package handlers

import (
	"strconv"

	"backend-kavling/internal/helpers"
	"backend-kavling/internal/repositories"
	"backend-kavling/internal/services"

	"github.com/gin-gonic/gin"
)

type KavlingHandler struct {
	repo    *repositories.KavlingRepository
	service *services.KavlingService
}

func NewKavlingHandler(repo *repositories.KavlingRepository, service *services.KavlingService) *KavlingHandler {
	return &KavlingHandler{repo: repo, service: service}
}

// List godoc
//
//	@Summary		List kavling
//	@Description	Mengambil daftar kavling. Bisa difilter berdasarkan denah_kavling_id.
//	@Tags			Kavling
//	@Produce		json
//	@Security		BearerAuth
//	@Param			denah_kavling_id	query		int	false	"Filter berdasarkan denah kavling ID"
//	@Success		200					{object}	object
//	@Router			/kavling [get]
func (h *KavlingHandler) List(c *gin.Context) {
	denahID, _ := strconv.Atoi(c.Query("denah_kavling_id"))
	list, err := h.repo.FindAll(denahID)
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil data kavling")
		return
	}
	helpers.OK(c, "OK", list)
}

// Detail godoc
//
//	@Summary		Detail kavling
//	@Description	Mengambil detail satu kavling berdasarkan ID
//	@Tags			Kavling
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Kavling ID"
//	@Success		200	{object} object
//	@Failure		404	{object} object
//	@Router			/kavling/{id} [get]
func (h *KavlingHandler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	k, err := h.repo.FindByID(id)
	if err != nil {
		helpers.NotFound(c, "Kavling tidak ditemukan")
		return
	}
	helpers.OK(c, "OK", k)
}

type updateKavlingRequest struct {
	PanjangKanan  float64 `json:"panjang_kanan" example:"10.5"`
	PanjangKiri   float64 `json:"panjang_kiri" example:"10.5"`
	LebarDepan    float64 `json:"lebar_depan" example:"6.0"`
	LebarBelakang float64 `json:"lebar_belakang" example:"6.0"`
	LuasTanah     float64 `json:"luas_tanah" example:"63.0"`
	HargaPerMeter float64 `json:"harga_per_meter" example:"2500000"`
	HargaJualCash float64 `json:"harga_jual_cash" example:"157500000"`
	Status        *int    `json:"status" example:"0"`
}

// Update godoc
//
//	@Summary		Update kavling
//	@Description	Mengubah data dimensi, harga, dan status kavling. kode_kavling dan kode_map bersifat immutable (tidak bisa diubah).
//	@Tags			Kavling
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"Kavling ID"
//	@Param			body	body		updateKavlingRequest	true	"Data update kavling"
//	@Success		200		{object} object
//	@Failure		400		{object} object
//	@Failure		404		{object} object
//	@Router			/kavling/{id} [put]
func (h *KavlingHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	k, err := h.repo.FindByID(id)
	if err != nil {
		helpers.NotFound(c, "Kavling tidak ditemukan")
		return
	}

	var req updateKavlingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	k.PanjangKanan = req.PanjangKanan
	k.PanjangKiri = req.PanjangKiri
	k.LebarDepan = req.LebarDepan
	k.LebarBelakang = req.LebarBelakang
	k.LuasTanah = req.LuasTanah
	k.HargaPerMeter = req.HargaPerMeter
	k.HargaJualCash = req.HargaJualCash
	if req.Status != nil {
		k.Status = *req.Status
	}

	if err := h.repo.Update(k); err != nil {
		helpers.InternalError(c, "Gagal update kavling")
		return
	}

	helpers.OK(c, "Kavling berhasil diupdate", k)
}

// Delete godoc
//
//	@Summary		Hapus kavling
//	@Description	Menghapus kavling. Hanya kavling dengan status Kosong (0) yang bisa dihapus.
//	@Tags			Kavling
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Kavling ID"
//	@Success		200	{object} object
//	@Failure		400	{object} object	"Kavling sudah ada transaksi"
//	@Failure		404	{object} object
//	@Router			/kavling/{id} [delete]
func (h *KavlingHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.ValidateDelete(id); err != nil {
		helpers.BadRequest(c, err.Error(), nil)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		helpers.InternalError(c, "Gagal menghapus kavling")
		return
	}

	helpers.OK(c, "Kavling berhasil dihapus", nil)
}
