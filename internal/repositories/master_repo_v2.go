package repositories

import (
	"backend-kavling/internal/models"

	"gorm.io/gorm"
)

// ─── Legalitas ────────────────────────────────────────────────────────────────

type LegalitasRepository struct{ db *gorm.DB }

func NewLegalitasRepository(db *gorm.DB) *LegalitasRepository {
	return &LegalitasRepository{db}
}

func (r *LegalitasRepository) List(idLokasi int, progres, q string) ([]models.Legalitas, error) {
	query := r.db.Model(&models.Legalitas{}).
		Preload("Kavling").Preload("Kavling.Lokasi").Preload("Kavling.Customer").Preload("Notaris")
	if idLokasi > 0 {
		query = query.Joins("JOIN kavling ON kavling.id = legalitas.id_kavling").
			Where("kavling.id_lokasi = ?", idLokasi)
	}
	if progres != "" {
		query = query.Where("legalitas.progres = ?", progres)
	}
	if q != "" {
		query = query.Joins("JOIN kavling ON kavling.id = legalitas.id_kavling").
			Where("kavling.kode_kavling ILIKE ? OR legalitas.atas_nama_surat ILIKE ?", "%"+q+"%", "%"+q+"%")
	}
	var list []models.Legalitas
	err := query.Order("legalitas.id ASC").Find(&list).Error
	return list, err
}

func (r *LegalitasRepository) FindByID(id int) (*models.Legalitas, error) {
	var m models.Legalitas
	err := r.db.Preload("Kavling").Preload("Kavling.Lokasi").Preload("Notaris").First(&m, id).Error
	return &m, err
}

func (r *LegalitasRepository) FindByKavling(idKavling int) (*models.Legalitas, error) {
	var m models.Legalitas
	err := r.db.Where("id_kavling = ?", idKavling).First(&m).Error
	return &m, err
}

func (r *LegalitasRepository) Create(m *models.Legalitas) error {
	return r.db.Create(m).Error
}

func (r *LegalitasRepository) Update(m *models.Legalitas) error {
	return r.db.Save(m).Error
}

// ─── Notaris ─────────────────────────────────────────────────────────────────

type NotarisRepository struct{ db *gorm.DB }

func NewNotarisRepository(db *gorm.DB) *NotarisRepository {
	return &NotarisRepository{db}
}

func (r *NotarisRepository) List() ([]models.Notaris, error) {
	var list []models.Notaris
	err := r.db.Order("nama ASC").Find(&list).Error
	return list, err
}

func (r *NotarisRepository) FindByID(id int) (*models.Notaris, error) {
	var m models.Notaris
	err := r.db.First(&m, id).Error
	return &m, err
}

func (r *NotarisRepository) Create(m *models.Notaris) error {
	return r.db.Create(m).Error
}

func (r *NotarisRepository) Update(m *models.Notaris) error {
	return r.db.Save(m).Error
}

func (r *NotarisRepository) Delete(id int) error {
	return r.db.Delete(&models.Notaris{}, id).Error
}

// ─── List Penjualan ───────────────────────────────────────────────────────────

type ListPenjualanRepository struct{ db *gorm.DB }

func NewListPenjualanRepository(db *gorm.DB) *ListPenjualanRepository {
	return &ListPenjualanRepository{db}
}

func (r *ListPenjualanRepository) List() ([]models.ListPenjualan, error) {
	var list []models.ListPenjualan
	err := r.db.Order("urutan ASC").Find(&list).Error
	return list, err
}

func (r *ListPenjualanRepository) FindByID(id int) (*models.ListPenjualan, error) {
	var m models.ListPenjualan
	err := r.db.First(&m, id).Error
	return &m, err
}

func (r *ListPenjualanRepository) Create(m *models.ListPenjualan) error {
	return r.db.Create(m).Error
}

func (r *ListPenjualanRepository) Update(m *models.ListPenjualan) error {
	return r.db.Save(m).Error
}

func (r *ListPenjualanRepository) Delete(id int) error {
	return r.db.Delete(&models.ListPenjualan{}, id).Error
}

// ─── Media ────────────────────────────────────────────────────────────────────

type MediaRepository struct{ db *gorm.DB }

func NewMediaRepository(db *gorm.DB) *MediaRepository {
	return &MediaRepository{db}
}

func (r *MediaRepository) List() ([]models.Media, error) {
	var list []models.Media
	err := r.db.Order("id ASC").Find(&list).Error
	return list, err
}

func (r *MediaRepository) FindByID(id int) (*models.Media, error) {
	var m models.Media
	err := r.db.First(&m, id).Error
	return &m, err
}

func (r *MediaRepository) Update(m *models.Media) error {
	return r.db.Save(m).Error
}

// ─── Landing Konten ───────────────────────────────────────────────────────────

type LandingRepository struct{ db *gorm.DB }

func NewLandingRepository(db *gorm.DB) *LandingRepository {
	return &LandingRepository{db}
}

func (r *LandingRepository) List(item string) ([]models.LandingKonten, error) {
	q := r.db.Model(&models.LandingKonten{})
	if item != "" {
		q = q.Where("item = ?", item)
	}
	var list []models.LandingKonten
	err := q.Order("item, urutan").Find(&list).Error
	return list, err
}

func (r *LandingRepository) FindByID(id int) (*models.LandingKonten, error) {
	var m models.LandingKonten
	err := r.db.First(&m, id).Error
	return &m, err
}

func (r *LandingRepository) Create(m *models.LandingKonten) error {
	return r.db.Create(m).Error
}

func (r *LandingRepository) Update(m *models.LandingKonten) error {
	return r.db.Save(m).Error
}

func (r *LandingRepository) Delete(id int) error {
	return r.db.Delete(&models.LandingKonten{}, id).Error
}

// ─── Beranda / Dashboard ─────────────────────────────────────────────────────

type BerandaRepository struct{ db *gorm.DB }

func NewBerandaRepository(db *gorm.DB) *BerandaRepository {
	return &BerandaRepository{db}
}

func (r *BerandaRepository) Ringkasan() (*models.DashboardRingkasan, error) {
	var res models.DashboardRingkasan

	type statusCount struct {
		Status int
		Count  int
	}
	var counts []statusCount
	r.db.Raw("SELECT status, COUNT(*) as count FROM kavling GROUP BY status").Scan(&counts)
	for _, c := range counts {
		switch c.Status {
		case 0:
			res.Kosong = c.Count
		case 1:
			res.Hold = c.Count
		case 2:
			res.BF = c.Count
		case 3:
			res.Akad = c.Count
		case 4:
			res.UserCancel = c.Count
		case 5:
			res.Lunas = c.Count
		}
		res.TotalKavling += c.Count
	}

	r.db.Raw("SELECT COUNT(*) FROM customer").Scan(&res.JumlahCustomer)
	r.db.Raw("SELECT COUNT(*) FROM prospek").Scan(&res.JumlahProspek)

	// Tunggakan: kredit yang bulan berjalan belum dibayar
	r.db.Raw(`
		SELECT COUNT(*) FROM transaksi_kavling tk
		WHERE tk.jenis_pembelian = 3
		AND (
			SELECT COUNT(*) FROM pembayaran p
			WHERE p.id_transaksi = tk.id
			AND EXTRACT(YEAR FROM p.tanggal) = EXTRACT(YEAR FROM NOW())
			AND EXTRACT(MONTH FROM p.tanggal) = EXTRACT(MONTH FROM NOW())
		) = 0
		AND tk.tgl_mulai_cicilan IS NOT NULL
		AND tk.tgl_mulai_cicilan <= NOW()
	`).Scan(&res.JumlahTunggakan)

	return &res, nil
}

func (r *BerandaRepository) AktifitasTerbaru(limit int) ([]models.Aktifitas, error) {
	var list []models.Aktifitas
	err := r.db.Preload("User").Order("created_at DESC").Limit(limit).Find(&list).Error
	return list, err
}

// ─── Pembayaran V2 (tagihan + pemasukan per transaksi) ────────────────────────

type PembayaranV2Repository struct{ db *gorm.DB }

func NewPembayaranV2Repository(db *gorm.DB) *PembayaranV2Repository {
	return &PembayaranV2Repository{db}
}

func (r *PembayaranV2Repository) List(q, progres string) ([]map[string]interface{}, error) {
	type row struct {
		ID           int     `gorm:"column:id"`
		NamaCustomer string  `gorm:"column:nama_customer"`
		KodeCustomer string  `gorm:"column:kode_customer"`
		NoTelp       string  `gorm:"column:no_telp"`
		Jenis        string  `gorm:"column:jenis_pembelian"`
		Lokasi       string  `gorm:"column:lokasi_perumahan"`
		KodeKavling  string  `gorm:"column:kode_kavling"`
		Progres      string  `gorm:"column:progres"`
		Marketing    string  `gorm:"column:nama_marketing"`
		Tagihan      float64 `gorm:"column:total_tagihan"`
		Terbayar     float64 `gorm:"column:total_terbayar"`
	}

	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if q != "" {
		whereClause += " AND (c.nama ILIKE ? OR c.kode_kontrak ILIKE ?)"
		args = append(args, "%"+q+"%", "%"+q+"%")
	}
	if progres != "" && progres != "Semua Progres" {
		statusMap := map[string]int{"Ready": 0, "HOLD": 1, "BF": 2, "AKAD": 3, "User Cancel": 4, "LUNAS": 5}
		if val, ok := statusMap[progres]; ok {
			whereClause += " AND k.status = ?"
			args = append(args, val)
		}
	}

	sql := `
		SELECT
			tk.id,
			COALESCE(c.nama,'')                                                       AS nama_customer,
			COALESCE(c.kode_kontrak,'')                                               AS kode_customer,
			COALESCE(c.no_telp,'')                                                    AS no_telp,
			CASE tk.jenis_pembelian WHEN 1 THEN 'Booking Fee' WHEN 2 THEN 'Cash Keras' WHEN 3 THEN 'Kredit' ELSE 'Lainnya' END AS jenis_pembelian,
			COALESCE(lk.nama,'')                                                      AS lokasi_perumahan,
			COALESCE(k.kode_kavling,'')                                               AS kode_kavling,
			CASE k.status WHEN 0 THEN 'Ready' WHEN 1 THEN 'HOLD' WHEN 2 THEN 'BF' WHEN 3 THEN 'AKAD' WHEN 4 THEN 'User Cancel' WHEN 5 THEN 'LUNAS' ELSE '' END AS progres,
			COALESCE(m.nama,'')                                                       AS nama_marketing,
			COALESCE((SELECT SUM(t.nominal)    FROM tagihan   t WHERE t.id_transaksi = tk.id), 0) AS total_tagihan,
			COALESCE((SELECT SUM(p.jumlah_bayar) FROM pembayaran p WHERE p.id_transaksi = tk.id), 0) AS total_terbayar
		FROM transaksi_kavling tk
		LEFT JOIN customer       c  ON c.id  = tk.id_customer
		LEFT JOIN kavling        k  ON k.id  = tk.id_kavling
		LEFT JOIN lokasi_kavling lk ON lk.id = k.id_lokasi
		LEFT JOIN marketing      m  ON m.id  = tk.id_marketing
		` + whereClause + `
		ORDER BY tk.id DESC`

	var rows []row
	if err := r.db.Raw(sql, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}

	list := make([]map[string]interface{}, 0, len(rows))
	for _, rw := range rows {
		list = append(list, map[string]interface{}{
			"id":               rw.ID,
			"id_transaksi":     rw.ID,
			"nama_customer":    rw.NamaCustomer,
			"kode_customer":    rw.KodeCustomer,
			"no_telp":          rw.NoTelp,
			"jenis_pembelian":  rw.Jenis,
			"lokasi_perumahan": rw.Lokasi,
			"kode_kavling":     rw.KodeKavling,
			"progres":          rw.Progres,
			"nama_marketing":   rw.Marketing,
			"total_tagihan":    rw.Tagihan,
			"total_terbayar":   rw.Terbayar,
		})
	}
	return list, nil
}

func (r *PembayaranV2Repository) Detail(idTransaksi int) (*models.PembayaranDetailResponse, error) {
	var tk models.TransaksiKavling
	if err := r.db.Preload("Kavling").Preload("Kavling.Lokasi").Preload("Customer").
		First(&tk, idTransaksi).Error; err != nil {
		return nil, err
	}

	var tagihans []models.Tagihan
	r.db.Where("id_transaksi = ?", idTransaksi).Preload("Kategori").Order("id ASC").Find(&tagihans)

	var pembayarans []models.Pembayaran
	r.db.Where("id_transaksi = ?", idTransaksi).Preload("Kategori").Preload("Bank").
		Order("tanggal ASC").Find(&pembayarans)

	var totalTagihan, totalBayar float64
	for _, t := range tagihans {
		totalTagihan += t.Nominal
	}
	for _, p := range pembayarans {
		totalBayar += p.JumlahBayar
	}

	lokasi := ""
	kode := ""
	if tk.Kavling != nil {
		if tk.Kavling.Lokasi != nil {
			lokasi = tk.Kavling.Lokasi.Nama
		}
		kode = tk.Kavling.KodeKavling
	}

	nama := ""
	if tk.Customer != nil {
		nama = tk.Customer.Nama
	}

	jenis := "Cash"
	switch tk.JenisPembelian {
	case 1:
		jenis = "Booking Fee"
	case 2:
		jenis = "Cash Keras"
	case 3:
		jenis = "Kredit"
	}

	sisa := totalTagihan - totalBayar
	if sisa < 0 {
		sisa = 0
	}

	return &models.PembayaranDetailResponse{
		NamaCustomer:    nama,
		LokasiPerumahan: lokasi,
		KodeKavling:     kode,
		JenisPembayaran: jenis,
		TotalTagihan:    totalTagihan,
		TotalBayar:      totalBayar,
		SisaBayar:       sisa,
		TagihanItems:    tagihans,
		PemasukanItems:  pembayarans,
	}, nil
}

func (r *PembayaranV2Repository) SumPembayaran(idTransaksi int) float64 {
	var total float64
	r.db.Model(&models.Pembayaran{}).Where("id_transaksi = ?", idTransaksi).
		Select("COALESCE(SUM(jumlah_bayar),0)").Scan(&total)
	return total
}

func (r *PembayaranV2Repository) NextPembayaranKe(idTransaksi int) int {
	var max int
	r.db.Model(&models.Pembayaran{}).Where("id_transaksi = ?", idTransaksi).
		Select("COALESCE(MAX(pembayaran_ke),0)").Scan(&max)
	return max + 1
}

func (r *PembayaranV2Repository) CreatePembayaran(m *models.Pembayaran) error {
	return r.db.Create(m).Error
}

func (r *PembayaranV2Repository) UpdatePembayaran(m *models.Pembayaran) error {
	return r.db.Save(m).Error
}

func (r *PembayaranV2Repository) DeletePembayaran(id int) error {
	return r.db.Delete(&models.Pembayaran{}, id).Error
}

func (r *PembayaranV2Repository) FindPembayaranByID(id int) (*models.Pembayaran, error) {
	var m models.Pembayaran
	err := r.db.First(&m, id).Error
	return &m, err
}

// ─── Jatuh Tempo ─────────────────────────────────────────────────────────────

type JatuhTempoRepository struct{ db *gorm.DB }

func NewJatuhTempoRepository(db *gorm.DB) *JatuhTempoRepository {
	return &JatuhTempoRepository{db}
}

func (r *JatuhTempoRepository) List(idLokasi int, jenis, q string) ([]models.JatuhTempoRow, error) {
	query := `
		SELECT
			tk.id as id_transaksi,
			c.nama,
			lk.nama as lokasi,
			k.kode_kavling,
			COALESCE((SELECT SUM(t.nominal) FROM tagihan t WHERE t.id_transaksi = tk.id), 0) as harga_tanah,
			COALESCE((SELECT SUM(p.jumlah_bayar) FROM pembayaran p WHERE p.id_transaksi = tk.id), 0) as pembayaran,
			0 as pencairan,
			COALESCE((SELECT SUM(t.nominal) FROM tagihan t WHERE t.id_transaksi = tk.id), 0) - COALESCE((SELECT SUM(p.jumlah_bayar) FROM pembayaran p WHERE p.id_transaksi = tk.id),0) as sisa,
			CASE
				-- Lunas: total bayar >= total tagihan
				WHEN COALESCE((SELECT SUM(p.jumlah_bayar) FROM pembayaran p WHERE p.id_transaksi = tk.id), 0)
					>= COALESCE((SELECT SUM(t.nominal) FROM tagihan t WHERE t.id_transaksi = tk.id), 1)
					AND COALESCE((SELECT SUM(t.nominal) FROM tagihan t WHERE t.id_transaksi = tk.id), 0) > 0
				THEN 'Lunas'
				-- Belum mulai: tgl_mulai_cicilan belum ada atau di masa depan
				WHEN tk.tgl_mulai_cicilan IS NULL OR tk.tgl_mulai_cicilan > NOW()
				THEN 'Belum Mulai'
				-- Hitung keterlambatan
				ELSE (
					CASE
						WHEN (
							-- Bulan berjalan sejak mulai cicilan (1-based)
							(EXTRACT(YEAR FROM AGE(NOW(), tk.tgl_mulai_cicilan)) * 12 +
							 EXTRACT(MONTH FROM AGE(NOW(), tk.tgl_mulai_cicilan)) + 1)::int
						) > (
							-- Jumlah pembayaran yang sudah masuk
							SELECT COUNT(*) FROM pembayaran p WHERE p.id_transaksi = tk.id
						)
						THEN 'Telat ' || (
							NOW()::date - (
								tk.tgl_mulai_cicilan + (
									(SELECT COUNT(*) FROM pembayaran p WHERE p.id_transaksi = tk.id)
									* INTERVAL '1 month'
								)
							)::date
						)::text || ' Hari'
						ELSE 'Lancar'
					END
				)
			END as keterlambatan,
			tk.cicilan_per_bulan as cicilan,
			tk.lama_cicilan as tenor,
			COALESCE(TO_CHAR(tk.tgl_mulai_cicilan, 'YYYY-MM-DD'), '') as jatuh_tempo,
			CASE tk.jenis_pembelian WHEN 2 THEN 'Cash Keras' WHEN 3 THEN 'Kredit' ELSE 'Cash' END as jenis_pembelian
		FROM transaksi_kavling tk
		JOIN customer c ON c.id = tk.id_customer
		JOIN kavling k ON k.id = tk.id_kavling
		JOIN lokasi_kavling lk ON lk.id = k.id_lokasi
		WHERE tk.jenis_pembelian = 3
	`
	args := []interface{}{}
	if idLokasi > 0 {
		query += " AND k.id_lokasi = ?"
		args = append(args, idLokasi)
	}
	if q != "" {
		query += " AND c.nama ILIKE ?"
		args = append(args, "%"+q+"%")
	}
	query += " ORDER BY c.nama ASC"

	var list []models.JatuhTempoRow
	err := r.db.Raw(query, args...).Scan(&list).Error
	return list, err
}
