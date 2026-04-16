package handlers

import (
	"backend-kavling/internal/helpers"
	"backend-kavling/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type loginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"admin123"`
	Captcha  string `json:"captcha" example:""`
}

type loginResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Login berhasil"`
	Data    struct {
		Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
		User  struct {
			ID       int    `json:"id" example:"1"`
			Username string `json:"username" example:"admin"`
			Nama     string `json:"nama" example:"Administrator"`
			IsAdmin  int    `json:"is_admin" example:"1"`
		} `json:"user"`
	} `json:"data"`
}

// Login godoc
//
//	@Summary		Login user
//	@Description	Autentikasi user dan mendapatkan JWT token. Setelah 5x gagal login, sistem meminta captcha.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		loginRequest	true	"Kredensial login"
//	@Success		200		{object}	loginResponse
//	@Failure		400		{object} object
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Username dan password wajib diisi", nil)
		return
	}

	ip := helpers.GetClientIP(c)

	if h.authService.NeedCaptcha(ip, req.Username) && req.Captcha == "" {
		c.JSON(400, gin.H{
			"success":      false,
			"message":      "Terlalu banyak percobaan login, harap isi captcha",
			"need_captcha": true,
		})
		return
	}

	result, needCaptcha, err := h.authService.Login(req.Username, req.Password, ip)
	if err != nil {
		c.JSON(400, gin.H{
			"success":      false,
			"message":      err.Error(),
			"need_captcha": needCaptcha,
		})
		return
	}

	helpers.OK(c, "Login berhasil", gin.H{
		"token": result.Token,
		"user":  result.User,
	})
}
