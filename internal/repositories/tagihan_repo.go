package repositories

import (
	"backend-kavling/internal/models"

	"gorm.io/gorm"
)

type TagihanRepository struct{ db *gorm.DB }

func NewTagihanRepository(db *gorm.DB) *TagihanRepository {
	return &TagihanRepository{db}
}

func (r *TagihanRepository) ListByTransaksi(idTransaksi int) ([]models.Tagihan, error) {
	var list []models.Tagihan
	err := r.db.Where("id_transaksi = ?", idTransaksi).
		Preload("Kategori").
		Order("id ASC").
		Find(&list).Error
	return list, err
}

func (r *TagihanRepository) FindByID(id int) (*models.Tagihan, error) {
	var m models.Tagihan
	err := r.db.Preload("Kategori").First(&m, id).Error
	return &m, err
}

func (r *TagihanRepository) Create(m *models.Tagihan) error {
	return r.db.Create(m).Error
}

func (r *TagihanRepository) Update(m *models.Tagihan) error {
	return r.db.Save(m).Error
}

func (r *TagihanRepository) Delete(id int) error {
	return r.db.Delete(&models.Tagihan{}, id).Error
}

func (r *TagihanRepository) FindKategoriByKode(kode string) int {
	var m models.KategoriTransaksi
	r.db.Where("kode = ?", kode).First(&m)
	return m.ID
}

func (r *TagihanRepository) SumByTransaksi(idTransaksi int) (float64, error) {
	var total float64
	err := r.db.Model(&models.Tagihan{}).
		Where("id_transaksi = ?", idTransaksi).
		Select("COALESCE(SUM(nominal), 0)").
		Scan(&total).Error
	return total, err
}
