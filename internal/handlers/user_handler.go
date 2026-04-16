package handlers

import (
	"strconv"

	"backend-kavling/internal/helpers"
	"backend-kavling/internal/models"
	"backend-kavling/internal/repositories"
	"backend-kavling/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo    *repositories.UserRepository
	authService *services.AuthService
}

func NewUserHandler(userRepo *repositories.UserRepository, authService *services.AuthService) *UserHandler {
	return &UserHandler{userRepo: userRepo, authService: authService}
}

// List godoc
//
//	@Summary		List semua user
//	@Description	Mengambil daftar semua user (admin only)
//	@Tags			Users
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object} object
//	@Failure		401	{object} object
//	@Failure		403	{object} object
//	@Router			/users [get]
func (h *UserHandler) List(c *gin.Context) {
	users, err := h.userRepo.FindAll()
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil data user")
		return
	}
	helpers.OK(c, "OK", users)
}

// Detail godoc
//
//	@Summary		Detail user
//	@Description	Mengambil detail satu user berdasarkan ID
//	@Tags			Users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object} object
//	@Failure		404	{object} object
//	@Router			/users/{id} [get]
func (h *UserHandler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.userRepo.FindByID(id)
	if err != nil {
		helpers.NotFound(c, "User tidak ditemukan")
		return
	}
	helpers.OK(c, "OK", user)
}

type createUserRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
	Nama     string `json:"nama" binding:"required" example:"John Doe"`
	IsAdmin  int    `json:"is_admin" example:"0"`
}

// Create godoc
//
//	@Summary		Tambah user baru
//	@Description	Membuat user baru (admin only)
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		createUserRequest	true	"Data user baru"
//	@Success		201		{object} object
//	@Failure		400		{object} object
//	@Router			/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak lengkap", err.Error())
		return
	}

	if existing, _ := h.userRepo.FindByUsername(req.Username); existing != nil {
		helpers.BadRequest(c, "Username sudah digunakan", nil)
		return
	}

	hashed, err := h.authService.HashPassword(req.Password)
	if err != nil {
		helpers.InternalError(c, "Gagal membuat password")
		return
	}

	user := &models.User{
		Username: req.Username,
		Password: hashed,
		Nama:     req.Nama,
		IsAdmin:  req.IsAdmin,
		Status:   "AKTIF",
	}

	if err := h.userRepo.Create(user); err != nil {
		helpers.InternalError(c, "Gagal membuat user")
		return
	}

	helpers.Created(c, "User berhasil dibuat", user)
}

type updateUserRequest struct {
	Nama     string `json:"nama" example:"John Doe"`
	Password string `json:"password" example:"newpassword"`
	IsAdmin  *int   `json:"is_admin" example:"0"`
	Status   string `json:"status" example:"AKTIF" enums:"AKTIF,BLOKIR"`
}

// Update godoc
//
//	@Summary		Update user
//	@Description	Mengubah data user (admin only)
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int					true	"User ID"
//	@Param			body	body		updateUserRequest	true	"Data update user"
//	@Success		200		{object} object
//	@Failure		404		{object} object
//	@Router			/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.userRepo.FindByID(id)
	if err != nil {
		helpers.NotFound(c, "User tidak ditemukan")
		return
	}

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	if req.Nama != "" {
		user.Nama = req.Nama
	}
	if req.Password != "" {
		hashed, err := h.authService.HashPassword(req.Password)
		if err != nil {
			helpers.InternalError(c, "Gagal update password")
			return
		}
		user.Password = hashed
	}
	if req.IsAdmin != nil {
		user.IsAdmin = *req.IsAdmin
	}
	if req.Status != "" {
		user.Status = req.Status
	}

	if err := h.userRepo.Update(user); err != nil {
		helpers.InternalError(c, "Gagal update user")
		return
	}

	helpers.OK(c, "User berhasil diupdate", user)
}

// Delete godoc
//
//	@Summary		Hapus user
//	@Description	Menghapus user berdasarkan ID (admin only, tidak bisa hapus diri sendiri)
//	@Tags			Users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object} object
//	@Failure		400	{object} object
//	@Failure		404	{object} object
//	@Router			/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	currentUserID, _ := c.Get("user_id")
	if currentUserID.(int) == id {
		helpers.BadRequest(c, "Tidak bisa menghapus akun sendiri", nil)
		return
	}

	if _, err := h.userRepo.FindByID(id); err != nil {
		helpers.NotFound(c, "User tidak ditemukan")
		return
	}

	if err := h.userRepo.Delete(id); err != nil {
		helpers.InternalError(c, "Gagal menghapus user")
		return
	}

	helpers.OK(c, "User berhasil dihapus", nil)
}

// GetAccess godoc
//
//	@Summary		Hak akses user
//	@Description	Mengambil daftar hak akses menu untuk user tertentu
//	@Tags			Users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object} object
//	@Router			/users/{id}/access [get]
func (h *UserHandler) GetAccess(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	akses, err := h.userRepo.FindHakAksesByUser(id)
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil data hak akses")
		return
	}
	helpers.OK(c, "OK", akses)
}

type updateAccessRequest struct {
	MenuIDs []int `json:"menu_ids" example:"1,2,3,4"`
}

// UpdateAccess godoc
//
//	@Summary		Update hak akses user
//	@Description	Mengganti seluruh hak akses menu untuk user (admin only)
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int					true	"User ID"
//	@Param			body	body		updateAccessRequest	true	"Array id_menu yang dicentang"
//	@Success		200		{object} object
//	@Router			/users/{id}/access [put]
func (h *UserHandler) UpdateAccess(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if _, err := h.userRepo.FindByID(id); err != nil {
		helpers.NotFound(c, "User tidak ditemukan")
		return
	}

	var req updateAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	if err := h.userRepo.ReplaceHakAkses(id, req.MenuIDs); err != nil {
		helpers.InternalError(c, "Gagal update hak akses")
		return
	}

	helpers.OK(c, "Hak akses berhasil diupdate", nil)
}

// ActivityLog godoc
//
//	@Summary		Log aktivitas user
//	@Description	Mengambil log aktivitas semua user (admin only)
//	@Tags			Users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			limit	query		int	false	"Jumlah record (default: 100)"
//	@Success		200		{object} object
//	@Router			/activity-log [get]
func (h *UserHandler) ActivityLog(c *gin.Context) {
	limit := 100
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	list, err := h.userRepo.FindAktifitas(limit)
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil log aktivitas")
		return
	}

	helpers.OK(c, "OK", list)
}

// ListMenu godoc
//
//	@Summary		List menu aplikasi
//	@Description	Mengambil daftar semua menu aplikasi
//	@Tags			Users
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object} object
//	@Router			/menu [get]
func (h *UserHandler) ListMenu(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(int)
	user, err := h.userRepo.FindByID(uid)
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil data user")
		return
	}

	var menus []models.Menu
	if user.IsAdmin == 1 {
		menus, err = h.userRepo.FindAllMenu()
	} else {
		menus, err = h.userRepo.FindMenusByUser(uid)
	}
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil menu")
		return
	}
	helpers.OK(c, "OK", menus)
}
