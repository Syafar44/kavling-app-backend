package repositories

import (
	"backend-kavling/internal/models"

	"gorm.io/gorm"
)

type LokasiKavlingRepository struct{ db *gorm.DB }

func NewLokasiKavlingRepository(db *gorm.DB) *LokasiKavlingRepository {
	return &LokasiKavlingRepository{db}
}

func (r *LokasiKavlingRepository) List() ([]models.LokasiKavling, error) {
	var list []models.LokasiKavling
	err := r.db.Order("urutan_lokasi ASC").Find(&list).Error
	return list, err
}

func (r *LokasiKavlingRepository) FindByID(id int) (*models.LokasiKavling, error) {
	var m models.LokasiKavling
	err := r.db.First(&m, id).Error
	return &m, err
}

func (r *LokasiKavlingRepository) Create(m *models.LokasiKavling) error {
	return r.db.Create(m).Error
}

func (r *LokasiKavlingRepository) Update(m *models.LokasiKavling) error {
	return r.db.Save(m).Error
}

func (r *LokasiKavlingRepository) Delete(id int) error {
	return r.db.Delete(&models.LokasiKavling{}, id).Error
}

func (r *LokasiKavlingRepository) UpdateJumlahKavling(id int) error {
	return r.db.Exec(
		"UPDATE lokasi_kavling SET jumlah_kavling = (SELECT COUNT(*) FROM kavling WHERE id_lokasi = ?) WHERE id = ?",
		id, id,
	).Error
}
