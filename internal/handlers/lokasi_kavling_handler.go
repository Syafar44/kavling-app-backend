package handlers

import (
	"net/http"
	"regexp"
	"strconv"

	"backend-kavling/internal/models"
	"backend-kavling/internal/repositories"

	"github.com/gin-gonic/gin"
)

type LokasiKavlingHandler struct {
	repo        *repositories.LokasiKavlingRepository
	kavlingRepo *repositories.KavlingV2Repository
}

func NewLokasiKavlingHandler(repo *repositories.LokasiKavlingRepository, kv *repositories.KavlingV2Repository) *LokasiKavlingHandler {
	return &LokasiKavlingHandler{repo, kv}
}

// extractPathIDs parses SVG content and returns all <path id="..."> values.
var rePath = regexp.MustCompile(`<path[^>]*\bid="([^"]+)"`)

func extractPathIDs(svg string) []string {
	matches := rePath.FindAllStringSubmatch(svg, -1)
	ids := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) > 1 && m[1] != "" {
			ids = append(ids, m[1])
		}
	}
	return ids
}

// GET /api/lokasi-kavling
func (h *LokasiKavlingHandler) List(c *gin.Context) {
	list, err := h.repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "ok", Data: list})
}

// GET /api/lokasi-kavling/:id
func (h *LokasiKavlingHandler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	m, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: m})
}

// POST /api/lokasi-kavling
func (h *LokasiKavlingHandler) Create(c *gin.Context) {
	var m models.LokasiKavling
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "invalid body", Errors: err.Error()})
		return
	}
	if err := h.repo.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}

	// Auto-generate kavling dari path IDs di SVG
	if m.SvgContent != "" {
		pathIDs := extractPathIDs(m.SvgContent)
		if len(pathIDs) > 0 {
			kavlings := make([]models.Kavling, 0, len(pathIDs))
			for _, id := range pathIDs {
				kavlings = append(kavlings, models.Kavling{
					IDLokasi:    m.ID,
					KodeKavling: id,
				})
			}
			_ = h.kavlingRepo.BulkCreate(kavlings)
			h.repo.UpdateJumlahKavling(m.ID)
		}
	}

	c.JSON(http.StatusCreated, models.Response{Success: true, Message: "lokasi berhasil dibuat", Data: m})
}

// PUT /api/lokasi-kavling/:id
func (h *LokasiKavlingHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	existing, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "tidak ditemukan"})
		return
	}
	if err := c.ShouldBindJSON(existing); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "invalid body"})
		return
	}
	existing.ID = id
	if err := h.repo.Update(existing); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "berhasil diupdate", Data: existing})
}

// DELETE /api/lokasi-kavling/:id
func (h *LokasiKavlingHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "berhasil dihapus"})
}
