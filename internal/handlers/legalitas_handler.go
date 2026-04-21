package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"backend-kavling/internal/config"
	"backend-kavling/internal/models"
	"backend-kavling/internal/repositories"

	"github.com/gin-gonic/gin"
)

type LegalitasHandler struct {
	repo *repositories.LegalitasRepository
}

func NewLegalitasHandler(r *repositories.LegalitasRepository) *LegalitasHandler {
	return &LegalitasHandler{r}
}

// GET /api/legalitas
func (h *LegalitasHandler) List(c *gin.Context) {
	idLokasi, _ := strconv.Atoi(c.Query("id_lokasi"))
	list, err := h.repo.List(idLokasi, c.Query("progres"), c.Query("q"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: list})
}

// GET /api/legalitas/:id
func (h *LegalitasHandler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	m, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: m})
}

// PUT /api/legalitas/:id
func (h *LegalitasHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	existing, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "tidak ditemukan"})
		return
	}

	// Handle multipart (bisa ada file bukti_foto)
	if c.ContentType() == "application/json" {
		c.ShouldBindJSON(existing)
	} else {
		existing.AtasNamaSurat = c.PostForm("atas_nama_surat")
		existing.NoSurat = c.PostForm("no_surat")
		existing.Progres = c.PostForm("progres")
		existing.Keterangan = c.PostForm("keterangan")
		if idNotarisStr := c.PostForm("id_notaris"); idNotarisStr != "" {
			idN, _ := strconv.Atoi(idNotarisStr)
			existing.IDNotaris = &idN
		}

		// Upload bukti foto
		if file, err := c.FormFile("bukti_foto"); err == nil {
			now := time.Now()
			dir := filepath.Join(config.AppConfig.UploadPath, "legalitas", now.Format("2006"), now.Format("01"))
			os.MkdirAll(dir, 0755)
			filename := strconv.Itoa(id) + "_" + now.Format("20060102150405") + filepath.Ext(file.Filename)
			dstPath := filepath.Join(dir, filename)
			if c.SaveUploadedFile(file, dstPath) == nil {
				existing.BuktiFoto = "legalitas/" + now.Format("2006") + "/" + now.Format("01") + "/" + filename
			}
		}
	}

	existing.ID = id
	if err := h.repo.Update(existing); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "berhasil diupdate", Data: existing})
}
