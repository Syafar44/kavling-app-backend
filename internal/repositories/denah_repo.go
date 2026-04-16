package repositories

import (
	"backend-kavling/internal/models"

	"gorm.io/gorm"
)

type DenahRepository struct {
	db *gorm.DB
}

func NewDenahRepository(db *gorm.DB) *DenahRepository {
	return &DenahRepository{db: db}
}

// FindAll returns all denah kavling without kavling list (lightweight)
func (r *DenahRepository) FindAll() ([]models.DenahKavling, error) {
	var list []models.DenahKavling
	err := r.db.Order("created_at DESC").Find(&list).Error
	return list, err
}

// FindAllWithKavling returns denah kavling with their kavlings (count only summary)
func (r *DenahRepository) FindAllWithSummary() ([]DenahSummary, error) {
	type row struct {
		ID           int    `json:"id"`
		Nama         string `json:"nama"`
		Viewbox      string `json:"viewbox"`
		JumlahKavling int64 `json:"jumlah_kavling"`
		JumlahKosong  int64 `json:"jumlah_kosong"`
		JumlahTerjual int64 `json:"jumlah_terjual"`
	}
	var rows []row
	err := r.db.Raw(`
		SELECT
			d.id, d.nama, d.viewbox,
			COUNT(k.id)                                              AS jumlah_kavling,
			COUNT(k.id) FILTER (WHERE k.status = 0)                 AS jumlah_kosong,
			COUNT(k.id) FILTER (WHERE k.status IN (2,3))            AS jumlah_terjual
		FROM denah_kavling d
		LEFT JOIN kavling_peta k ON k.denah_kavling_id = d.id
		GROUP BY d.id, d.nama, d.viewbox
		ORDER BY d.created_at DESC
	`).Scan(&rows).Error

	result := make([]DenahSummary, 0, len(rows))
	for _, r := range rows {
		result = append(result, DenahSummary{
			ID:            r.ID,
			Nama:          r.Nama,
			Viewbox:       r.Viewbox,
			JumlahKavling: int(r.JumlahKavling),
			JumlahKosong:  int(r.JumlahKosong),
			JumlahTerjual: int(r.JumlahTerjual),
		})
	}
	return result, err
}

type DenahSummary struct {
	ID            int    `json:"id"`
	Nama          string `json:"nama"`
	Viewbox       string `json:"viewbox"`
	JumlahKavling int    `json:"jumlah_kavling"`
	JumlahKosong  int    `json:"jumlah_kosong"`
	JumlahTerjual int    `json:"jumlah_terjual"`
}

// FindByID returns a denah with all its kavlings
func (r *DenahRepository) FindByID(id int) (*models.DenahKavling, error) {
	var d models.DenahKavling
	err := r.db.
		Preload("Kavlings", func(db *gorm.DB) *gorm.DB {
			return db.Order("kode_kavling ASC")
		}).
		First(&d, id).Error
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DenahRepository) Create(d *models.DenahKavling) error {
	return r.db.Create(d).Error
}

func (r *DenahRepository) Update(d *models.DenahKavling) error {
	return r.db.Save(d).Error
}

// Delete removes a denah and all its kavlings (CASCADE in DB)
func (r *DenahRepository) Delete(id int) error {
	return r.db.Delete(&models.DenahKavling{}, id).Error
}

// HasTransaksi checks if any kavling in the denah has transactions
func (r *DenahRepository) HasTransaksi(id int) bool {
	var count int64
	r.db.Model(&TransaksiKavlingModel{}).
		Joins("JOIN kavling_peta ON kavling_peta.id = transaksi_kavling.id_kavling").
		Where("kavling_peta.denah_kavling_id = ?", id).
		Count(&count)
	return count > 0
}

// Tiny helper model just for the join query above
type TransaksiKavlingModel struct{}

func (TransaksiKavlingModel) TableName() string { return "transaksi_kavling" }
