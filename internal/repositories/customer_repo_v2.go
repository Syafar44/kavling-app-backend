package repositories

import (
	"backend-kavling/internal/models"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type CustomerV2Repository struct{ db *gorm.DB }

func NewCustomerV2Repository(db *gorm.DB) *CustomerV2Repository {
	return &CustomerV2Repository{db}
}

func (r *CustomerV2Repository) List(q, idLokasi, status, jenisPembelian string, page, perPage int) ([]models.Customer, int64, error) {
	query := r.db.Model(&models.Customer{}).Preload("Lokasi").Preload("Marketing")
	if q != "" {
		query = query.Where("nama ILIKE ? OR kode_kontrak ILIKE ? OR no_ktp ILIKE ?", "%"+q+"%", "%"+q+"%", "%"+q+"%")
	}
	if idLokasi != "" {
		query = query.Where("id_lokasi = ?", idLokasi)
	}
	if status != "" {
		query = query.Where("status_penjualan = ?", status)
	}
	if jenisPembelian != "" {
		query = query.Where("jenis_pembelian = ?", jenisPembelian)
	}
	var total int64
	query.Count(&total)
	var list []models.Customer
	err := query.Order("id DESC").Offset((page-1)*perPage).Limit(perPage).Find(&list).Error
	return list, total, err
}

func (r *CustomerV2Repository) FindByID(id int) (*models.Customer, error) {
	var m models.Customer
	err := r.db.Preload("Lokasi").Preload("Marketing").First(&m, id).Error
	return &m, err
}

func (r *CustomerV2Repository) Create(m *models.Customer) error {
	return r.db.Create(m).Error
}

func (r *CustomerV2Repository) Update(m *models.Customer) error {
	return r.db.Save(m).Error
}

func (r *CustomerV2Repository) Delete(id int) error {
	return r.db.Delete(&models.Customer{}, id).Error
}

func (r *CustomerV2Repository) GenerateKodeKontrak(namaLokasi, jenisPembelian string) (string, error) {
	// Format: <SINGKAT>-<JENIS>-<NNNN>
	singkat := "XX"
	var lok models.LokasiKavling
	if err := r.db.Where("nama = ?", namaLokasi).First(&lok).Error; err == nil {
		singkat = lok.NamaSingkat
	}

	jenis := "KRD"
	switch strings.ToUpper(jenisPembelian) {
	case "CASH KERAS":
		jenis = "CSH"
	case "CASH BERTAHAP":
		jenis = "CST"
	}

	prefix := fmt.Sprintf("%s-%s", singkat, jenis)
	var count int64
	r.db.Model(&models.Customer{}).Where("kode_kontrak LIKE ?", prefix+"%").Count(&count)
	return fmt.Sprintf("%s-%04d", prefix, count+1), nil
}

// ─── Customer Kavling ─────────────────────────────────────────────────────────

func (r *CustomerV2Repository) AddKavling(idCustomer, idKavling int) error {
	ck := models.CustomerKavling{IDCustomer: idCustomer, IDKavling: idKavling}
	return r.db.Create(&ck).Error
}

func (r *CustomerV2Repository) RemoveKavling(idCustomer, idKavling int) error {
	return r.db.Where("id_customer = ? AND id_kavling = ?", idCustomer, idKavling).
		Delete(&models.CustomerKavling{}).Error
}

func (r *CustomerV2Repository) GetKavlings(idCustomer int) ([]models.CustomerKavling, error) {
	var list []models.CustomerKavling
	err := r.db.Where("id_customer = ?", idCustomer).Preload("Kavling").Find(&list).Error
	return list, err
}

// ─── Customer File ────────────────────────────────────────────────────────────

func (r *CustomerV2Repository) ListFiles(idCustomer int) ([]models.CustomerFile, error) {
	var list []models.CustomerFile
	err := r.db.Where("id_customer = ?", idCustomer).Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *CustomerV2Repository) CreateFile(m *models.CustomerFile) error {
	return r.db.Create(m).Error
}

func (r *CustomerV2Repository) DeleteFile(id int) error {
	return r.db.Delete(&models.CustomerFile{}, id).Error
}

// ─── Customer Arsip ───────────────────────────────────────────────────────────

func (r *CustomerV2Repository) ListArsip(q string) ([]models.CustomerArsip, error) {
	query := r.db.Model(&models.CustomerArsip{}).Preload("Marketing")
	if q != "" {
		query = query.Where("nama ILIKE ? OR lokasi ILIKE ?", "%"+q+"%", "%"+q+"%")
	}
	var list []models.CustomerArsip
	err := query.Order("tanggal DESC").Find(&list).Error
	return list, err
}

func (r *CustomerV2Repository) CreateArsip(m *models.CustomerArsip) error {
	return r.db.Create(m).Error
}

func (r *CustomerV2Repository) DeleteArsip(id int) error {
	return r.db.Delete(&models.CustomerArsip{}, id).Error
}

// ─── Prospek ─────────────────────────────────────────────────────────────────

func (r *CustomerV2Repository) ListProspek(q string) ([]models.Prospek, error) {
	query := r.db.Model(&models.Prospek{}).Preload("Marketing")
	if q != "" {
		query = query.Where("nama ILIKE ? OR no_telp ILIKE ?", "%"+q+"%", "%"+q+"%")
	}
	var list []models.Prospek
	err := query.Order("tanggal DESC").Find(&list).Error
	return list, err
}

func (r *CustomerV2Repository) FindProspekByID(id int) (*models.Prospek, error) {
	var m models.Prospek
	err := r.db.Preload("Marketing").First(&m, id).Error
	return &m, err
}

func (r *CustomerV2Repository) CreateProspek(m *models.Prospek) error {
	return r.db.Create(m).Error
}

func (r *CustomerV2Repository) UpdateProspek(m *models.Prospek) error {
	return r.db.Save(m).Error
}

func (r *CustomerV2Repository) DeleteProspek(id int) error {
	return r.db.Delete(&models.Prospek{}, id).Error
}
