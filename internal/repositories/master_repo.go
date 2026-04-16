package repositories

import (
	"backend-kavling/internal/models"

	"gorm.io/gorm"
)

// ─── Marketing Repository ─────────────────────────────────────────────────────

type MarketingRepository struct {
	db *gorm.DB
}

func NewMarketingRepository(db *gorm.DB) *MarketingRepository {
	return &MarketingRepository{db: db}
}

func (r *MarketingRepository) FindAll() ([]models.Marketing, error) {
	var list []models.Marketing
	err := r.db.Order("nama ASC").Find(&list).Error
	return list, err
}

func (r *MarketingRepository) FindByID(id int) (*models.Marketing, error) {
	var m models.Marketing
	err := r.db.First(&m, id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MarketingRepository) Create(m *models.Marketing) error {
	return r.db.Create(m).Error
}

func (r *MarketingRepository) Update(m *models.Marketing) error {
	return r.db.Save(m).Error
}

func (r *MarketingRepository) Delete(id int) error {
	return r.db.Delete(&models.Marketing{}, id).Error
}

func (r *MarketingRepository) HasTransaksi(id int) bool {
	var count int64
	r.db.Model(&models.TransaksiKavling{}).Where("id_marketing = ?", id).Count(&count)
	return count > 0
}

func (r *MarketingRepository) CountOpenClosed(id int) (open, closed int) {
	// OPEN: status kavling 1 (booking) atau 3 (kredit)
	r.db.Model(&models.TransaksiKavling{}).
		Joins("JOIN kavling_peta ON kavling_peta.id = transaksi_kavling.id_kavling").
		Where("transaksi_kavling.id_marketing = ? AND kavling_peta.status IN (1,3)", id).
		Count((*int64)(func() *int64 { n := int64(open); return &n }()))

	var openCount, closedCount int64
	r.db.Model(&models.TransaksiKavling{}).
		Joins("JOIN kavling_peta ON kavling_peta.id = transaksi_kavling.id_kavling").
		Where("transaksi_kavling.id_marketing = ? AND kavling_peta.status IN (1,3)", id).
		Count(&openCount)
	r.db.Model(&models.TransaksiKavling{}).
		Joins("JOIN kavling_peta ON kavling_peta.id = transaksi_kavling.id_kavling").
		Where("transaksi_kavling.id_marketing = ? AND kavling_peta.status = 2", id).
		Count(&closedCount)
	return int(openCount), int(closedCount)
}

// ─── Customer Repository ──────────────────────────────────────────────────────

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) FindAll() ([]models.Customer, error) {
	var list []models.Customer
	err := r.db.Order("nama ASC").Find(&list).Error
	return list, err
}

func (r *CustomerRepository) FindByID(id int) (*models.Customer, error) {
	var c models.Customer
	err := r.db.First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CustomerRepository) Create(c *models.Customer) error {
	return r.db.Create(c).Error
}

func (r *CustomerRepository) Update(c *models.Customer) error {
	return r.db.Save(c).Error
}

func (r *CustomerRepository) Delete(id int) error {
	return r.db.Delete(&models.Customer{}, id).Error
}

// ─── Bank Repository ──────────────────────────────────────────────────────────

type BankRepository struct {
	db *gorm.DB
}

func NewBankRepository(db *gorm.DB) *BankRepository {
	return &BankRepository{db: db}
}

func (r *BankRepository) FindAll() ([]models.Bank, error) {
	var list []models.Bank
	err := r.db.Where("status = 1").Order("nama_bank ASC").Find(&list).Error
	return list, err
}

func (r *BankRepository) FindByID(id int) (*models.Bank, error) {
	var b models.Bank
	err := r.db.First(&b, id).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BankRepository) Create(b *models.Bank) error {
	return r.db.Create(b).Error
}

func (r *BankRepository) Update(b *models.Bank) error {
	return r.db.Save(b).Error
}

func (r *BankRepository) Delete(id int) error {
	return r.db.Delete(&models.Bank{}, id).Error
}

// ─── Konfigurasi Repository ───────────────────────────────────────────────────

type KonfigurasiRepository struct {
	db *gorm.DB
}

func NewKonfigurasiRepository(db *gorm.DB) *KonfigurasiRepository {
	return &KonfigurasiRepository{db: db}
}

func (r *KonfigurasiRepository) Get() (*models.Konfigurasi, error) {
	var k models.Konfigurasi
	err := r.db.First(&k).Error
	if err != nil {
		return nil, err
	}
	return &k, nil
}

func (r *KonfigurasiRepository) Update(k *models.Konfigurasi) error {
	return r.db.Save(k).Error
}

func (r *KonfigurasiRepository) GetWA() (*models.KonfigurasiWA, error) {
	var wa models.KonfigurasiWA
	err := r.db.Where("is_aktif = ?", 1).First(&wa).Error
	if err != nil {
		return nil, err
	}
	return &wa, nil
}

func (r *KonfigurasiRepository) UpsertWA(wa *models.KonfigurasiWA) error {
	if wa.ID > 0 {
		return r.db.Save(wa).Error
	}
	return r.db.Create(wa).Error
}

func (r *KonfigurasiRepository) FindAllTemplate() ([]models.Template, error) {
	var list []models.Template
	err := r.db.Find(&list).Error
	return list, err
}

func (r *KonfigurasiRepository) SaveTemplate(t *models.Template) error {
	if t.ID > 0 {
		return r.db.Save(t).Error
	}
	return r.db.Create(t).Error
}

// ResetAllData menghapus semua data transaksi dan master (kecuali users, menu, hak_akses, konfigurasi).
func (r *KonfigurasiRepository) ResetAllData() error {
	return r.ResetSelected([]string{"pembayaran", "transaksi", "keuangan", "kavling", "customer", "marketing"})
}

// ResetSelected menghapus data berdasarkan kategori yang dipilih.
// Kategori: pembayaran, transaksi, keuangan, kavling, customer, marketing
func (r *KonfigurasiRepository) ResetSelected(items []string) error {
	// Mapping kategori → tabel yang perlu dihapus
	groupTables := map[string][]string{
		"pembayaran": {"pembayaran"},
		"transaksi":  {"pembayaran", "transaksi_kavling", "transaksi_booking"},
		"keuangan":   {"transaksi"},
		"kavling":    {"pembayaran", "transaksi_kavling", "transaksi_booking", "kavling_peta", "denah_kavling"},
		"customer":   {"pembayaran", "transaksi_kavling", "transaksi_booking", "customer"},
		"marketing":  {"pembayaran", "transaksi_kavling", "transaksi_booking", "marketing"},
	}

	// Kumpulkan tabel unik lalu urutkan sesuai dependency (children first)
	seen := map[string]bool{}
	for _, item := range items {
		for _, tbl := range groupTables[item] {
			seen[tbl] = true
		}
	}

	order := []string{
		"pembayaran", "transaksi_kavling", "transaksi_booking",
		"transaksi", "kavling_peta", "denah_kavling", "customer", "marketing",
	}
	for _, tbl := range order {
		if seen[tbl] {
			if err := r.db.Exec("TRUNCATE TABLE " + tbl + " RESTART IDENTITY CASCADE").Error; err != nil {
				return err
			}
		}
	}
	return nil
}
