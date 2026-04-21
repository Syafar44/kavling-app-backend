package handlers

import (
	"net/http"
	"strconv"

	"backend-kavling/internal/models"
	"backend-kavling/internal/repositories"

	"github.com/gin-gonic/gin"
)

type SiteplanHandler struct {
	lokasiRepo  *repositories.LokasiKavlingRepository
	kavlingRepo *repositories.KavlingV2Repository
}

func NewSiteplanHandler(lr *repositories.LokasiKavlingRepository, kr *repositories.KavlingV2Repository) *SiteplanHandler {
	return &SiteplanHandler{lr, kr}
}

// GET /api/siteplan — list semua lokasi + kavlings per lokasi
func (h *SiteplanHandler) List(c *gin.Context) {
	lokasis, err := h.lokasiRepo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: lokasis})
}

// GET /api/siteplan/:id_lokasi — satu lokasi + array kavling dengan status
func (h *SiteplanHandler) Detail(c *gin.Context) {
	idLokasi, _ := strconv.Atoi(c.Param("id_lokasi"))
	lokasi, err := h.lokasiRepo.FindByID(idLokasi)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "lokasi tidak ditemukan"})
		return
	}
	kavlings, err := h.kavlingRepo.ListByLokasi(idLokasi)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data: gin.H{
			"lokasi":   lokasi,
			"kavlings": kavlings,
		},
	})
}

// GET /api/siteplan/:id_lokasi/kavling/:kode — detail 1 kavling + customer
func (h *SiteplanHandler) KavlingDetail(c *gin.Context) {
	idLokasi, _ := strconv.Atoi(c.Param("id_lokasi"))
	kode := c.Param("kode")

	kavling, err := h.kavlingRepo.FindByKode(idLokasi, kode)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "kavling tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data: gin.H{
			"kavling":  kavling,
			"lokasi":   kavling.Lokasi,
			"customer": kavling.Customer,
		},
	})
}
