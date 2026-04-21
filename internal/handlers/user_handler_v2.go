package handlers

import (
	"net/http"
	"strconv"

	"backend-kavling/internal/models"
	"backend-kavling/internal/repositories"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HakAksesExtHandler manages extended permission matrix (lihat/tambah/edit/hapus/beranda per menu)
type HakAksesExtHandler struct {
	userRepo *repositories.UserRepository
	db       *gorm.DB
}

func NewHakAksesExtHandler(ur *repositories.UserRepository, db *gorm.DB) *HakAksesExtHandler {
	return &HakAksesExtHandler{userRepo: ur, db: db}
}

// GET /api/hak-akses?id_user=X
func (h *HakAksesExtHandler) GetMatrix(c *gin.Context) {
	idUser, _ := strconv.Atoi(c.Query("id_user"))
	if idUser == 0 {
		c.JSON(http.StatusBadRequest, models.Response{Message: "id_user wajib diisi"})
		return
	}

	var allMenu []models.Menu
	h.db.Order("urutan ASC").Find(&allMenu)

	var hakAkses []models.HakAkses
	h.db.Where("id_user = ?", idUser).Preload("Menu").Find(&hakAkses)

	aksesByMenu := make(map[int]models.HakAkses)
	for _, ha := range hakAkses {
		aksesByMenu[ha.IDMenu] = ha
	}

	type MatrixRow struct {
		Menu    models.Menu       `json:"menu"`
		HakAkses models.HakAkses `json:"hak_akses"`
	}
	var matrix []MatrixRow
	for _, m := range allMenu {
		ha, ok := aksesByMenu[m.ID]
		if !ok {
			ha = models.HakAkses{IDUser: idUser, IDMenu: m.ID}
		}
		matrix = append(matrix, MatrixRow{Menu: m, HakAkses: ha})
	}

	c.JSON(http.StatusOK, models.Response{Success: true, Data: matrix})
}

// PUT /api/hak-akses/bulk
func (h *HakAksesExtHandler) BulkUpdate(c *gin.Context) {
	var req struct {
		IDUser      int `json:"id_user" binding:"required"`
		Permissions []struct {
			IDMenu  int  `json:"id_menu"`
			Lihat   bool `json:"lihat"`
			Beranda bool `json:"beranda"`
			Tambah  bool `json:"tambah"`
			Edit    bool `json:"edit"`
			Hapus   bool `json:"hapus"`
		} `json:"permissions" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		tx.Where("id_user = ?", req.IDUser).Delete(&models.HakAkses{})
		for _, p := range req.Permissions {
			ha := models.HakAkses{
				IDUser:    req.IDUser,
				IDMenu:    p.IDMenu,
				StatusHak: 1,
				Lihat:     p.Lihat,
				Beranda:   p.Beranda,
				Tambah:    p.Tambah,
				Edit:      p.Edit,
				Hapus:     p.Hapus,
			}
			if err := tx.Create(&ha).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "hak akses berhasil diupdate"})
}
