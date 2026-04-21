package handlers

import (
	"net/http"
	"strconv"

	"backend-kavling/internal/models"
	"backend-kavling/internal/repositories"

	"github.com/gin-gonic/gin"
)

type BerandaHandler struct {
	berandaRepo *repositories.BerandaRepository
	keuanganRepo *repositories.KeuanganRepository
}

func NewBerandaHandler(br *repositories.BerandaRepository, kr *repositories.KeuanganRepository) *BerandaHandler {
	return &BerandaHandler{br, kr}
}

// GET /api/beranda/ringkasan
func (h *BerandaHandler) Ringkasan(c *gin.Context) {
	data, err := h.berandaRepo.Ringkasan()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: data})
}

// GET /api/beranda/arus-kas?tahun=2026
func (h *BerandaHandler) ArusKas(c *gin.Context) {
	tahunStr := c.Query("tahun")
	tahun, _ := strconv.Atoi(tahunStr)
	if tahun == 0 {
		tahun = 2026
	}
	data, err := h.keuanganRepo.ArusKasBulanan(tahun)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: data})
}

// GET /api/beranda/aktifitas?limit=10
func (h *BerandaHandler) Aktifitas(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	data, err := h.berandaRepo.AktifitasTerbaru(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: data})
}
