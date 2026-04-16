package repositories

import (
	"time"

	"backend-kavling/internal/models"

	"gorm.io/gorm"
)

type KavlingRepository struct {
	db *gorm.DB
}

func NewKavlingRepository(db *gorm.DB) *KavlingRepository {
	return &KavlingRepository{db: db}
}

// FindAll returns kavlings, optionally filtered by denah_kavling_id
func (r *KavlingRepository) FindAll(denahID int) ([]models.KavlingPeta, error) {
	var list []models.KavlingPeta
	q := r.db.Order("kode_kavling ASC")
	if denahID > 0 {
		q = q.Where("denah_kavling_id = ?", denahID)
	}
	err := q.Find(&list).Error
	return list, err
}

func (r *KavlingRepository) FindByID(id int) (*models.KavlingPeta, error) {
	var k models.KavlingPeta
	err := r.db.First(&k, id).Error
	if err != nil {
		return nil, err
	}
	return &k, nil
}

func (r *KavlingRepository) FindByKode(kode string) (*models.KavlingPeta, error) {
	var k models.KavlingPeta
	err := r.db.Where("kode_kavling = ?", kode).First(&k).Error
	if err != nil {
		return nil, err
	}
	return &k, nil
}

func (r *KavlingRepository) Create(k *models.KavlingPeta) error {
	return r.db.Create(k).Error
}

// CreateBatch digunakan oleh DenahService saat parsing SVG
func (r *KavlingRepository) CreateBatch(list []models.KavlingPeta) error {
	return r.db.CreateInBatches(list, 100).Error
}

func (r *KavlingRepository) Update(k *models.KavlingPeta) error {
	return r.db.Save(k).Error
}

func (r *KavlingRepository) Delete(id int) error {
	return r.db.Delete(&models.KavlingPeta{}, id).Error
}

func (r *KavlingRepository) UpdateStatus(id, status int) error {
	return r.db.Model(&models.KavlingPeta{}).Where("id = ?", id).
		Update("status", status).Error
}

func (r *KavlingRepository) UpdateTglJatuhTempo(id int, tgl *time.Time) error {
	return r.db.Model(&models.KavlingPeta{}).Where("id = ?", id).
		Update("tgl_jatuh_tempo", tgl).Error
}

func (r *KavlingRepository) HasTransaksi(id int) bool {
	var count int64
	r.db.Model(&models.TransaksiKavling{}).Where("id_kavling = ?", id).Count(&count)
	if count > 0 {
		return true
	}
	r.db.Model(&models.TransaksiBooking{}).Where("id_kavling = ? AND status = 1", id).Count(&count)
	return count > 0
}

func (r *KavlingRepository) ExistsByKode(kode string, excludeID int) bool {
	var count int64
	q := r.db.Model(&models.KavlingPeta{}).Where("kode_kavling = ?", kode)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	q.Count(&count)
	return count > 0
}
