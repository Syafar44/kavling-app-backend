package repositories

import (
	"backend-kavling/internal/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type KeuanganRepository struct{ db *gorm.DB }

func NewKeuanganRepository(db *gorm.DB) *KeuanganRepository {
	return &KeuanganRepository{db}
}

// ─── Pemasukan / Pengeluaran (dari tabel transaksi) ──────────────────────────

func (r *KeuanganRepository) ListTransaksi(jenis, q string, tahun, bulan, idBank, idKategori, page, perPage int) ([]models.Transaksi, int64, error) {
	query := r.db.Model(&models.Transaksi{}).Preload("Bank")
	if jenis != "" {
		query = query.Where("jenis = ?", jenis)
	}
	if q != "" {
		query = query.Where("kategori ILIKE ? OR keterangan ILIKE ?", "%"+q+"%", "%"+q+"%")
	}
	if tahun > 0 {
		query = query.Where("EXTRACT(YEAR FROM tanggal) = ?", tahun)
	}
	if bulan > 0 {
		query = query.Where("EXTRACT(MONTH FROM tanggal) = ?", bulan)
	}
	if idBank > 0 {
		query = query.Where("id_bank = ?", idBank)
	}

	var total int64
	query.Count(&total)

	var list []models.Transaksi
	err := query.Order("tanggal DESC, id DESC").
		Offset((page - 1) * perPage).Limit(perPage).
		Find(&list).Error
	return list, total, err
}

func (r *KeuanganRepository) FindTransaksiByID(id int) (*models.Transaksi, error) {
	var m models.Transaksi
	err := r.db.Preload("Bank").First(&m, id).Error
	return &m, err
}

func (r *KeuanganRepository) CreateTransaksi(m *models.Transaksi) error {
	return r.db.Create(m).Error
}

func (r *KeuanganRepository) DeleteTransaksi(id int) error {
	return r.db.Delete(&models.Transaksi{}, id).Error
}

func (r *KeuanganRepository) GenerateNoTransaksi(prefix string) string {
	today := time.Now().Format("20060102")
	var count int64
	r.db.Model(&models.Transaksi{}).
		Where("no_transaksi LIKE ?", prefix+"-"+today+"%").
		Count(&count)
	return fmt.Sprintf("%s-%s-%04d", prefix, today, count+1)
}

// ─── Hutang ──────────────────────────────────────────────────────────────────

func (r *KeuanganRepository) ListHutang(q string) ([]models.Hutang, error) {
	query := r.db.Model(&models.Hutang{}).Preload("Bank")
	if q != "" {
		query = query.Where("deskripsi ILIKE ?", "%"+q+"%")
	}
	var list []models.Hutang
	err := query.Order("tanggal DESC").Find(&list).Error
	return list, err
}

func (r *KeuanganRepository) FindHutangByID(id int) (*models.Hutang, error) {
	var m models.Hutang
	err := r.db.Preload("Bank").First(&m, id).Error
	return &m, err
}

func (r *KeuanganRepository) CreateHutang(m *models.Hutang) error {
	m.SisaBayar = m.Nominal
	return r.db.Create(m).Error
}

func (r *KeuanganRepository) UpdateHutang(m *models.Hutang) error {
	return r.db.Save(m).Error
}

func (r *KeuanganRepository) DeleteHutang(id int) error {
	return r.db.Delete(&models.Hutang{}, id).Error
}

func (r *KeuanganRepository) BayarHutang(id int, nominal float64) error {
	var m models.Hutang
	if err := r.db.First(&m, id).Error; err != nil {
		return err
	}
	m.SisaBayar -= nominal
	if m.SisaBayar <= 0 {
		m.SisaBayar = 0
		m.Status = "Sudah Lunas"
		now := time.Now()
		m.TanggalPelunasan = &now
	}
	return r.db.Save(&m).Error
}

// ─── Piutang ─────────────────────────────────────────────────────────────────

func (r *KeuanganRepository) ListPiutang(q string) ([]models.Piutang, error) {
	query := r.db.Model(&models.Piutang{}).Preload("Bank")
	if q != "" {
		query = query.Where("deskripsi ILIKE ?", "%"+q+"%")
	}
	var list []models.Piutang
	err := query.Order("tanggal DESC").Find(&list).Error
	return list, err
}

func (r *KeuanganRepository) FindPiutangByID(id int) (*models.Piutang, error) {
	var m models.Piutang
	err := r.db.Preload("Bank").First(&m, id).Error
	return &m, err
}

func (r *KeuanganRepository) CreatePiutang(m *models.Piutang) error {
	m.SisaBayar = m.Nominal
	return r.db.Create(m).Error
}

func (r *KeuanganRepository) UpdatePiutang(m *models.Piutang) error {
	return r.db.Save(m).Error
}

func (r *KeuanganRepository) DeletePiutang(id int) error {
	return r.db.Delete(&models.Piutang{}, id).Error
}

func (r *KeuanganRepository) BayarPiutang(id int, nominal float64) error {
	var m models.Piutang
	if err := r.db.First(&m, id).Error; err != nil {
		return err
	}
	m.SisaBayar -= nominal
	if m.SisaBayar <= 0 {
		m.SisaBayar = 0
		m.Status = "Sudah Lunas"
		now := time.Now()
		m.TanggalPelunasan = &now
	}
	return r.db.Save(&m).Error
}

// ─── Mutasi Saldo ─────────────────────────────────────────────────────────────

func (r *KeuanganRepository) ListMutasi() ([]models.MutasiSaldo, error) {
	var list []models.MutasiSaldo
	err := r.db.Preload("BankAsal").Preload("BankTujuan").Order("tanggal DESC").Find(&list).Error
	return list, err
}

func (r *KeuanganRepository) CreateMutasi(m *models.MutasiSaldo) error {
	return r.db.Create(m).Error
}

func (r *KeuanganRepository) DeleteMutasi(id int) error {
	return r.db.Delete(&models.MutasiSaldo{}, id).Error
}

// ─── Bank saldo operations ────────────────────────────────────────────────────

func (r *KeuanganRepository) AddSaldo(idBank int, nominal float64) error {
	return r.db.Exec("UPDATE bank SET saldo = saldo + ? WHERE id = ?", nominal, idBank).Error
}

func (r *KeuanganRepository) SubtractSaldo(idBank int, nominal float64) error {
	return r.db.Exec("UPDATE bank SET saldo = saldo - ? WHERE id = ?", nominal, idBank).Error
}

// ─── Laporan Arus Kas ─────────────────────────────────────────────────────────

func (r *KeuanganRepository) ArusKas(tahun, bulan, idBank, idLokasi int) ([]models.Transaksi, error) {
	query := r.db.Model(&models.Transaksi{}).Preload("Bank").
		Where("EXTRACT(YEAR FROM tanggal) = ?", tahun)
	if bulan > 0 {
		query = query.Where("EXTRACT(MONTH FROM tanggal) = ?", bulan)
	}
	if idBank > 0 {
		query = query.Where("id_bank = ?", idBank)
	}
	if idLokasi > 0 {
		query = query.Where(`
			referensi_tipe = 'pembayaran' AND referensi_id IN (
				SELECT id FROM transaksi_kavling WHERE id_kavling IN (
					SELECT id FROM kavling WHERE id_lokasi = ?
				)
			)
		`, idLokasi)
	}
	var list []models.Transaksi
	err := query.Order("tanggal ASC").Find(&list).Error
	return list, err
}

func (r *KeuanganRepository) ArusKasBulanan(tahun int) ([]models.ArusKasBulan, error) {
	type rawRow struct {
		Bulan       int
		Jenis       string
		TotalNominal float64
	}
	var rows []rawRow
	err := r.db.Raw(`
		SELECT EXTRACT(MONTH FROM tanggal)::INT AS bulan, jenis, COALESCE(SUM(nominal),0) AS total_nominal
		FROM transaksi
		WHERE EXTRACT(YEAR FROM tanggal) = ?
		GROUP BY bulan, jenis
		ORDER BY bulan
	`, tahun).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]models.ArusKasBulan, 12)
	for i := range result {
		result[i].Bulan = i + 1
	}
	for _, r := range rows {
		idx := r.Bulan - 1
		if idx >= 0 && idx < 12 {
			if r.Jenis == "Pemasukan" {
				result[idx].Pemasukan = r.TotalNominal
			} else {
				result[idx].Pengeluaran = r.TotalNominal
			}
		}
	}
	return result, nil
}
