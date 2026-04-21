package repositories

import (
	"backend-kavling/internal/models"

	"gorm.io/gorm"
)

type KavlingV2Repository struct{ db *gorm.DB }

func NewKavlingV2Repository(db *gorm.DB) *KavlingV2Repository {
	return &KavlingV2Repository{db}
}

func (r *KavlingV2Repository) ListByLokasi(idLokasi int) ([]models.Kavling, error) {
	var list []models.Kavling
	err := r.db.Where("id_lokasi = ?", idLokasi).
		Preload("Customer").
		Order("kode_kavling ASC").
		Find(&list).Error
	return list, err
}

func (r *KavlingV2Repository) FindByKode(idLokasi int, kode string) (*models.Kavling, error) {
	var m models.Kavling
	err := r.db.Where("id_lokasi = ? AND kode_kavling = ?", idLokasi, kode).
		Preload("Lokasi").
		Preload("Customer").
		First(&m).Error
	return &m, err
}

func (r *KavlingV2Repository) FindByID(id int) (*models.Kavling, error) {
	var m models.Kavling
	err := r.db.Preload("Lokasi").Preload("Customer").First(&m, id).Error
	return &m, err
}

func (r *KavlingV2Repository) List(idLokasi int, status int) ([]models.Kavling, error) {
	q := r.db.Model(&models.Kavling{})
	if idLokasi > 0 {
		q = q.Where("id_lokasi = ?", idLokasi)
	}
	if status >= 0 {
		q = q.Where("status = ?", status)
	}
	var list []models.Kavling
	err := q.Preload("Lokasi").Preload("Customer").Order("id_lokasi, kode_kavling").Find(&list).Error
	return list, err
}

func (r *KavlingV2Repository) Create(m *models.Kavling) error {
	return r.db.Create(m).Error
}

func (r *KavlingV2Repository) BulkCreate(items []models.Kavling) error {
	return r.db.CreateInBatches(items, 100).Error
}

func (r *KavlingV2Repository) Update(m *models.Kavling) error {
	return r.db.Save(m).Error
}

func (r *KavlingV2Repository) UpdateStatus(id, status int, idCustomer *int) error {
	updates := map[string]interface{}{"status": status}
	if idCustomer != nil {
		updates["id_customer"] = *idCustomer
	}
	return r.db.Model(&models.Kavling{}).Where("id = ?", id).Updates(updates).Error
}

func (r *KavlingV2Repository) CountByStatus(idLokasi int) (map[int]int64, error) {
	type row struct {
		Status int
		Count  int64
	}
	var rows []row
	q := r.db.Model(&models.Kavling{}).Select("status, count(*) as count").Group("status")
	if idLokasi > 0 {
		q = q.Where("id_lokasi = ?", idLokasi)
	}
	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}
	m := make(map[int]int64)
	for _, r := range rows {
		m[r.Status] = r.Count
	}
	return m, nil
}
