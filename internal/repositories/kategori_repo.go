package repositories

import (
	"backend-kavling/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type KategoriRepository struct{ db *gorm.DB }

func NewKategoriRepository(db *gorm.DB) *KategoriRepository {
	return &KategoriRepository{db}
}

func (r *KategoriRepository) List(jenis string) ([]models.KategoriTransaksi, error) {
	q := r.db.Model(&models.KategoriTransaksi{})
	if jenis != "" {
		q = q.Where("jenis = ?", jenis)
	}
	var list []models.KategoriTransaksi
	err := q.Order("kode ASC").Find(&list).Error
	return list, err
}

func (r *KategoriRepository) FindByID(id int) (*models.KategoriTransaksi, error) {
	var m models.KategoriTransaksi
	err := r.db.First(&m, id).Error
	return &m, err
}

func (r *KategoriRepository) NextKode() (string, error) {
	var maxKode string
	r.db.Model(&models.KategoriTransaksi{}).Select("COALESCE(MAX(kode), '000')").Scan(&maxKode)
	var n int
	fmt.Sscanf(maxKode, "%d", &n)
	return fmt.Sprintf("%03d", n+1), nil
}

func (r *KategoriRepository) Create(m *models.KategoriTransaksi) error {
	return r.db.Create(m).Error
}

func (r *KategoriRepository) Update(m *models.KategoriTransaksi) error {
	return r.db.Save(m).Error
}

func (r *KategoriRepository) Delete(id int) error {
	var m models.KategoriTransaksi
	if err := r.db.First(&m, id).Error; err != nil {
		return err
	}
	if m.IsSystem {
		return fmt.Errorf("kategori sistem tidak dapat dihapus")
	}
	return r.db.Delete(&m).Error
}
