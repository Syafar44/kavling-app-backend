package services

import (
	"errors"
	"time"

	"backend-kavling/internal/middleware"
	"backend-kavling/internal/models"
	"backend-kavling/internal/repositories"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

const maxLoginAttempts = 5

type LoginResult struct {
	Token string
	User  *models.User
}

func (s *AuthService) Login(username, password, ip string) (*LoginResult, bool, error) {
	// Cek throttle
	throttle, _ := s.userRepo.GetThrottle(ip, username)
	needCaptcha := throttle != nil && throttle.Attempts >= maxLoginAttempts

	// Cari user
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = s.userRepo.IncrementThrottle(ip, username)
			return nil, needCaptcha, errors.New("username tidak ditemukan")
		}
		return nil, needCaptcha, err
	}

	// Cek status blokir
	if user.Status == "BLOKIR" {
		return nil, needCaptcha, errors.New("akun anda diblokir, hubungi administrator")
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		_ = s.userRepo.IncrementThrottle(ip, username)
		newThrottle, _ := s.userRepo.GetThrottle(ip, username)
		needCaptcha = newThrottle != nil && newThrottle.Attempts >= maxLoginAttempts
		return nil, needCaptcha, errors.New("password salah")
	}

	// Reset throttle setelah login berhasil
	_ = s.userRepo.ResetThrottle(ip, username)

	// Generate JWT
	token, err := middleware.GenerateToken(user.ID, user.Username, user.IsAdmin)
	if err != nil {
		return nil, false, err
	}

	// Catat aktifitas
	userID := user.ID
	_ = s.userRepo.CreateAktifitas(&models.Aktifitas{
		IDUser:    &userID,
		Aksi:      "LOGIN",
		Keterangan: "User berhasil login",
		IPAddress: ip,
		CreatedAt: time.Now(),
	})

	return &LoginResult{Token: token, User: user}, false, nil
}

func (s *AuthService) NeedCaptcha(ip, username string) bool {
	throttle, err := s.userRepo.GetThrottle(ip, username)
	if err != nil {
		return false
	}
	return throttle.Attempts >= maxLoginAttempts
}

func (s *AuthService) HashPassword(plain string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
