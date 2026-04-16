package repositories

import (
	"backend-kavling/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Order("id ASC").Find(&users).Error
	return users, err
}

func (r *UserRepository) FindByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id int) error {
	return r.db.Delete(&models.User{}, id).Error
}

// ─── Menu & Hak Akses ────────────────────────────────────────────────────────

func (r *UserRepository) FindAllMenu() ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.Order("urutan ASC").Find(&menus).Error
	return menus, err
}

func (r *UserRepository) FindMenusByUser(userID int) ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.
		Joins("JOIN hak_akses ha ON ha.id_menu = menu.id AND ha.id_user = ? AND ha.status_hak = 1", userID).
		Order("menu.urutan ASC").
		Find(&menus).Error
	return menus, err
}

func (r *UserRepository) FindHakAksesByUser(userID int) ([]models.HakAkses, error) {
	var akses []models.HakAkses
	err := r.db.Preload("Menu").Where("id_user = ?", userID).Find(&akses).Error
	return akses, err
}

func (r *UserRepository) ReplaceHakAkses(userID int, menuIDs []int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Hapus semua hak akses lama user
		if err := tx.Where("id_user = ?", userID).Delete(&models.HakAkses{}).Error; err != nil {
			return err
		}
		// Buat hak akses baru
		for _, menuID := range menuIDs {
			ha := models.HakAkses{
				IDUser:    userID,
				IDMenu:    menuID,
				StatusHak: 1,
			}
			if err := tx.Create(&ha).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *UserRepository) HasAccess(userID, menuID int) bool {
	var count int64
	r.db.Model(&models.HakAkses{}).
		Where("id_user = ? AND id_menu = ? AND status_hak = 1", userID, menuID).
		Count(&count)
	return count > 0
}

// ─── Throttle ─────────────────────────────────────────────────────────────────

func (r *UserRepository) GetThrottle(ip, username string) (*models.Throttle, error) {
	var t models.Throttle
	err := r.db.Where("ip_address = ? AND username = ?", ip, username).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *UserRepository) IncrementThrottle(ip, username string) error {
	result := r.db.Model(&models.Throttle{}).
		Where("ip_address = ? AND username = ?", ip, username).
		Updates(map[string]interface{}{"attempts": gorm.Expr("attempts + 1"), "last_attempt": gorm.Expr("NOW()")})
	if result.RowsAffected == 0 {
		return r.db.Create(&models.Throttle{
			IPAddress: ip,
			Username:  username,
			Attempts:  1,
		}).Error
	}
	return result.Error
}

func (r *UserRepository) ResetThrottle(ip, username string) error {
	return r.db.Where("ip_address = ? AND username = ?", ip, username).
		Delete(&models.Throttle{}).Error
}

// ─── Aktifitas ─────────────────────────────────────────────────────────────────

func (r *UserRepository) CreateAktifitas(a *models.Aktifitas) error {
	return r.db.Create(a).Error
}

func (r *UserRepository) FindAktifitas(limit int) ([]models.Aktifitas, error) {
	var list []models.Aktifitas
	err := r.db.Preload("User").Order("created_at DESC").Limit(limit).Find(&list).Error
	return list, err
}
