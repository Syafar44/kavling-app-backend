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

// ─── Notaris Handler ──────────────────────────────────────────────────────────

type NotarisHandler struct {
	repo *repositories.NotarisRepository
}

func NewNotarisHandler(r *repositories.NotarisRepository) *NotarisHandler {
	return &NotarisHandler{r}
}

func (h *NotarisHandler) List(c *gin.Context) {
	list, _ := h.repo.List()
	c.JSON(http.StatusOK, models.Response{Success: true, Data: list})
}

func (h *NotarisHandler) Create(c *gin.Context) {
	var m models.Notaris
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}
	if m.Keterangan == "" {
		m.Keterangan = "NOTARIS"
	}
	h.repo.Create(&m)
	c.JSON(http.StatusCreated, models.Response{Success: true, Data: m})
}

func (h *NotarisHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	existing, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "tidak ditemukan"})
		return
	}
	c.ShouldBindJSON(existing)
	existing.ID = id
	h.repo.Update(existing)
	c.JSON(http.StatusOK, models.Response{Success: true, Data: existing})
}

func (h *NotarisHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.repo.Delete(id)
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "dihapus"})
}

// ─── List Penjualan Handler ───────────────────────────────────────────────────

type ListPenjualanHandler struct {
	repo *repositories.ListPenjualanRepository
}

func NewListPenjualanHandler(r *repositories.ListPenjualanRepository) *ListPenjualanHandler {
	return &ListPenjualanHandler{r}
}

func (h *ListPenjualanHandler) List(c *gin.Context) {
	list, _ := h.repo.List()
	c.JSON(http.StatusOK, models.Response{Success: true, Data: list})
}

func (h *ListPenjualanHandler) Create(c *gin.Context) {
	var m models.ListPenjualan
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}
	h.repo.Create(&m)
	c.JSON(http.StatusCreated, models.Response{Success: true, Data: m})
}

func (h *ListPenjualanHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	existing, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "tidak ditemukan"})
		return
	}
	c.ShouldBindJSON(existing)
	existing.ID = id
	h.repo.Update(existing)
	c.JSON(http.StatusOK, models.Response{Success: true, Data: existing})
}

func (h *ListPenjualanHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.repo.Delete(id)
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "dihapus"})
}

// ─── Media Handler ────────────────────────────────────────────────────────────

type MediaHandler struct {
	repo *repositories.MediaRepository
}

func NewMediaHandler(r *repositories.MediaRepository) *MediaHandler {
	return &MediaHandler{r}
}

func (h *MediaHandler) List(c *gin.Context) {
	list, _ := h.repo.List()
	c.JSON(http.StatusOK, models.Response{Success: true, Data: list})
}

func (h *MediaHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	existing, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "tidak ditemukan"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: "file wajib diupload"})
		return
	}

	now := time.Now()
	dir := filepath.Join(config.AppConfig.UploadPath, "media")
	os.MkdirAll(dir, 0755)
	filename := strconv.Itoa(id) + "_" + now.Format("20060102150405") + filepath.Ext(file.Filename)
	dstPath := filepath.Join(dir, filename)
	if err := c.SaveUploadedFile(file, dstPath); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: "gagal menyimpan file"})
		return
	}

	existing.NamaFile = file.Filename
	existing.PathFile = "media/" + filename
	h.repo.Update(existing)
	c.JSON(http.StatusOK, models.Response{Success: true, Data: existing})
}

// ─── Landing Handler ──────────────────────────────────────────────────────────

type LandingHandler struct {
	repo *repositories.LandingRepository
}

func NewLandingHandler(r *repositories.LandingRepository) *LandingHandler {
	return &LandingHandler{r}
}

func (h *LandingHandler) List(c *gin.Context) {
	list, _ := h.repo.List(c.Query("item"))
	c.JSON(http.StatusOK, models.Response{Success: true, Data: list})
}

func (h *LandingHandler) Create(c *gin.Context) {
	var m models.LandingKonten
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}
	h.repo.Create(&m)
	c.JSON(http.StatusCreated, models.Response{Success: true, Data: m})
}

func (h *LandingHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	existing, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "tidak ditemukan"})
		return
	}
	c.ShouldBindJSON(existing)
	existing.ID = id
	h.repo.Update(existing)
	c.JSON(http.StatusOK, models.Response{Success: true, Data: existing})
}

func (h *LandingHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.repo.Delete(id)
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "dihapus"})
}

// ─── Kavling V2 Handler ───────────────────────────────────────────────────────

type KavlingV2Handler struct {
	repo       *repositories.KavlingV2Repository
	lokasiRepo *repositories.LokasiKavlingRepository
}

func NewKavlingV2Handler(r *repositories.KavlingV2Repository, l *repositories.LokasiKavlingRepository) *KavlingV2Handler {
	return &KavlingV2Handler{r, l}
}

// GET /api/kavling?id_lokasi=1&status=0
func (h *KavlingV2Handler) List(c *gin.Context) {
	idLokasi, _ := strconv.Atoi(c.Query("id_lokasi"))
	status := -1
	if s := c.Query("status"); s != "" {
		status, _ = strconv.Atoi(s)
	}
	list, err := h.repo.List(idLokasi, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: list})
}

// GET /api/kavling/:id
func (h *KavlingV2Handler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	m, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: m})
}

// PUT /api/kavling/:id — update dimensi, harga, sertipikat
func (h *KavlingV2Handler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	existing, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "tidak ditemukan"})
		return
	}
	var req struct {
		PanjangKanan  float64 `json:"panjang_kanan"`
		PanjangKiri   float64 `json:"panjang_kiri"`
		LebarDepan    float64 `json:"lebar_depan"`
		LebarBelakang float64 `json:"lebar_belakang"`
		LuasTanah     float64 `json:"luas_tanah"`
		HargaPerMeter float64 `json:"harga_per_meter"`
		HargaJualCash float64 `json:"harga_jual_cash"`
		NoSertipikat  string  `json:"no_sertipikat"`
		Keterangan    string  `json:"keterangan"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}
	if req.PanjangKanan > 0 { existing.PanjangKanan = req.PanjangKanan }
	if req.PanjangKiri > 0  { existing.PanjangKiri = req.PanjangKiri }
	if req.LebarDepan > 0   { existing.LebarDepan = req.LebarDepan }
	if req.LebarBelakang > 0 { existing.LebarBelakang = req.LebarBelakang }
	if req.LuasTanah > 0    { existing.LuasTanah = req.LuasTanah }
	if req.HargaPerMeter > 0 { existing.HargaPerMeter = req.HargaPerMeter }
	if req.HargaJualCash > 0 { existing.HargaJualCash = req.HargaJualCash }
	if req.NoSertipikat != "" { existing.NoSertipikat = req.NoSertipikat }
	if req.Keterangan != ""  { existing.Keterangan = req.Keterangan }

	if err := h.repo.Update(existing); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: existing})
}

// POST /api/kavling/bulk — bulk create kavlings untuk satu lokasi
func (h *KavlingV2Handler) BulkCreate(c *gin.Context) {
	var req struct {
		IDLokasi int              `json:"id_lokasi" binding:"required"`
		Kavlings []models.Kavling `json:"kavlings" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}
	for i := range req.Kavlings {
		req.Kavlings[i].IDLokasi = req.IDLokasi
	}
	if err := h.repo.BulkCreate(req.Kavlings); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	h.lokasiRepo.UpdateJumlahKavling(req.IDLokasi)
	c.JSON(http.StatusCreated, models.Response{Success: true, Message: "kavling berhasil dibuat", Data: len(req.Kavlings)})
}

// ─── Jatuh Tempo Handler ──────────────────────────────────────────────────────

type JatuhTempoHandler struct {
	repo *repositories.JatuhTempoRepository
}

func NewJatuhTempoHandler(r *repositories.JatuhTempoRepository) *JatuhTempoHandler {
	return &JatuhTempoHandler{r}
}

// GET /api/jatuh-tempo
func (h *JatuhTempoHandler) List(c *gin.Context) {
	idLokasi, _ := strconv.Atoi(c.Query("id_lokasi"))
	list, err := h.repo.List(idLokasi, c.Query("jenis"), c.Query("q"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: list})
}
