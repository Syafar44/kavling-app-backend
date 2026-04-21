package repositories

import (
	"fmt"
	"log"
	"strings"
	"time"

	"backend-kavling/internal/models"

	"gorm.io/gorm"
)

// ─── Booking Repository ───────────────────────────────────────────────────────

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) FindAll() ([]models.TransaksiBooking, error) {
	var list []models.TransaksiBooking
	err := r.db.
		Preload("Kavling").
		Preload("Customer").
		Preload("Marketing").
		Where("status = 1").
		Order("created_at DESC").
		Find(&list).Error
	return list, err
}

func (r *BookingRepository) FindByID(id int) (*models.TransaksiBooking, error) {
	var b models.TransaksiBooking
	err := r.db.
		Preload("Kavling").
		Preload("Customer").
		Preload("Marketing").
		First(&b, id).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BookingRepository) Create(b *models.TransaksiBooking) error {
	return r.db.Create(b).Error
}

func (r *BookingRepository) Update(b *models.TransaksiBooking) error {
	return r.db.Save(b).Error
}

// ─── Transaksi Kavling Repository ────────────────────────────────────────────

type TransaksiKavlingRepository struct {
	db *gorm.DB
}

func NewTransaksiKavlingRepository(db *gorm.DB) *TransaksiKavlingRepository {
	return &TransaksiKavlingRepository{db: db}
}

func (r *TransaksiKavlingRepository) FindAll() ([]models.TransaksiKavling, error) {
	var list []models.TransaksiKavling
	err := r.db.
		Preload("Kavling").
		Preload("Customer").
		Preload("Marketing").
		Order("created_at DESC").
		Find(&list).Error
	return list, err
}

func (r *TransaksiKavlingRepository) FindByID(id int) (*models.TransaksiKavling, error) {
	var t models.TransaksiKavling
	err := r.db.
		Preload("Kavling").
		Preload("Customer").
		Preload("Marketing").
		First(&t, id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TransaksiKavlingRepository) FindByKavling(idKavling int) (*models.TransaksiKavling, error) {
	var t models.TransaksiKavling
	err := r.db.
		Preload("Kavling").
		Preload("Customer").
		Preload("Marketing").
		Where("id_kavling = ?", idKavling).
		Last(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TransaksiKavlingRepository) Create(t *models.TransaksiKavling) error {
	return r.db.Create(t).Error
}

func (r *TransaksiKavlingRepository) FindAllKredit() ([]models.TransaksiKavling, error) {
	var list []models.TransaksiKavling
	err := r.db.
		Preload("Kavling").
		Preload("Customer").
		Where("jenis_pembelian = 3").
		Order("created_at DESC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	// Isi computed fields: jumlah cicilan terbayar & total terbayar
	for i := range list {
		type result struct {
			Count int
			Total float64
		}
		var res result
		r.db.Model(&models.Pembayaran{}).
			Select("COUNT(*) as count, COALESCE(SUM(jumlah_bayar), 0) as total").
			Where("id_kavling = ?", list[i].IDKavling).
			Scan(&res)
		list[i].JumlahCicilanTerbayar = res.Count
		list[i].TotalTerbayar = res.Total
	}
	return list, nil
}

// ─── Pembayaran Repository ────────────────────────────────────────────────────

type PembayaranRepository struct {
	db *gorm.DB
}

func NewPembayaranRepository(db *gorm.DB) *PembayaranRepository {
	return &PembayaranRepository{db: db}
}

func (r *PembayaranRepository) FindByKavling(idKavling int) ([]models.Pembayaran, error) {
	var list []models.Pembayaran
	err := r.db.
		Preload("Bank").
		Preload("Customer").
		Where("id_kavling = ?", idKavling).
		Order("pembayaran_ke ASC").
		Find(&list).Error
	return list, err
}

func (r *PembayaranRepository) CountByKavling(idKavling int) int {
	var count int64
	r.db.Model(&models.Pembayaran{}).Where("id_kavling = ?", idKavling).Count(&count)
	return int(count)
}

func (r *PembayaranRepository) HasPaidThisMonth(idKavling int, year, month int) bool {
	var count int64
	r.db.Model(&models.Pembayaran{}).
		Where("id_kavling = ? AND EXTRACT(YEAR FROM tanggal) = ? AND EXTRACT(MONTH FROM tanggal) = ?",
			idKavling, year, month).
		Count(&count)
	return count > 0
}

func (r *PembayaranRepository) Create(p *models.Pembayaran) error {
	return r.db.Create(p).Error
}

func (r *PembayaranRepository) FindAll() ([]models.Pembayaran, error) {
	var list []models.Pembayaran
	err := r.db.
		Preload("Kavling").
		Preload("Customer").
		Preload("Bank").
		Order("created_at DESC").
		Find(&list).Error
	return list, err
}

// ─── Arus Kas (Transaksi) Repository ─────────────────────────────────────────

type ArusKasRepository struct {
	db *gorm.DB
}

func NewArusKasRepository(db *gorm.DB) *ArusKasRepository {
	return &ArusKasRepository{db: db}
}

func (r *ArusKasRepository) FindAll(dari, sampai string) ([]models.Transaksi, error) {
	var list []models.Transaksi
	q := r.db.Preload("Bank").Order("tanggal DESC, created_at DESC")

	if dari != "" {
		q = q.Where("tanggal >= ?", dari)
	}
	if sampai != "" {
		q = q.Where("tanggal <= ?", sampai)
	}

	err := q.Find(&list).Error
	return list, err
}

func (r *ArusKasRepository) FindByID(id int) (*models.Transaksi, error) {
	var t models.Transaksi
	err := r.db.Preload("Bank").First(&t, id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *ArusKasRepository) Create(t *models.Transaksi) error {
	return r.db.Create(t).Error
}

func (r *ArusKasRepository) Delete(id int) error {
	return r.db.Delete(&models.Transaksi{}, id).Error
}

func (r *ArusKasRepository) SumByJenis(dari, sampai string) (pemasukan, pengeluaran float64) {
	type Result struct {
		Jenis   string
		Total   float64
	}
	var results []Result
	q := r.db.Model(&models.Transaksi{}).
		Select("jenis, SUM(nominal) as total").
		Group("jenis")
	if dari != "" {
		q = q.Where("tanggal >= ?", dari)
	}
	if sampai != "" {
		q = q.Where("tanggal <= ?", sampai)
	}
	q.Scan(&results)

	for _, r := range results {
		if r.Jenis == "Pemasukan" {
			pemasukan = r.Total
		} else {
			pengeluaran = r.Total
		}
	}
	return
}

func (r *ArusKasRepository) FindByKategori(kategori, dari, sampai string) ([]models.Transaksi, error) {
	var list []models.Transaksi
	q := r.db.Preload("Bank").Where("kategori = ?", kategori).Order("tanggal DESC")
	if dari != "" {
		q = q.Where("tanggal >= ?", dari)
	}
	if sampai != "" {
		q = q.Where("tanggal <= ?", sampai)
	}
	err := q.Find(&list).Error
	return list, err
}

func (r *ArusKasRepository) FindByCustomer(idCustomer int, dari, sampai string) ([]models.Transaksi, error) {
	// Fetch ALL pembayaran for this customer (no date filter — we need all IDs)
	var pembayarans []models.Pembayaran
	if err := r.db.Where("id_customer = ?", idCustomer).Find(&pembayarans).Error; err != nil {
		return nil, err
	}

	if len(pembayarans) == 0 {
		return []models.Transaksi{}, nil
	}

	ids := make([]int, 0, len(pembayarans))
	for _, p := range pembayarans {
		ids = append(ids, p.ID)
	}

	// Apply date filter on the arus kas query, not on pembayaran
	q := r.db.Preload("Bank").
		Where("referensi_id IN ? AND referensi_tipe = 'pembayaran'", ids).
		Order("tanggal DESC")
	if dari != "" {
		q = q.Where("tanggal >= ?", dari)
	}
	if sampai != "" {
		q = q.Where("tanggal <= ?", sampai)
	}

	var list []models.Transaksi
	err := q.Find(&list).Error
	return list, err
}

// GetLaporanUmum returns transactions with date filter
func (r *ArusKasRepository) GetLaporanUmum(dari, sampai string) ([]models.Transaksi, error) {
	return r.FindAll(dari, sampai)
}

// GetStatistik returns a summary of transactions by year-month
func (r *ArusKasRepository) GetStatistik(year int) (interface{}, error) {
	type MonthStat struct {
		Bulan      int     `json:"bulan"`
		Pemasukan  float64 `json:"pemasukan"`
		Pengeluaran float64 `json:"pengeluaran"`
	}
	var results []MonthStat
	err := r.db.Raw(`
		SELECT
			EXTRACT(MONTH FROM tanggal)::int AS bulan,
			SUM(CASE WHEN jenis = 'Pemasukan' THEN nominal ELSE 0 END) AS pemasukan,
			SUM(CASE WHEN jenis = 'Pengeluaran' THEN nominal ELSE 0 END) AS pengeluaran
		FROM transaksi
		WHERE EXTRACT(YEAR FROM tanggal) = ?
		GROUP BY EXTRACT(MONTH FROM tanggal)
		ORDER BY bulan ASC
	`, year).Scan(&results).Error
	return results, err
}

// SaldoBank updates bank saldo after a transaction
func (r *ArusKasRepository) UpdateBankSaldo(idBank int, delta float64) error {
	return r.db.Model(&models.Bank{}).Where("id = ?", idBank).
		Update("saldo", gorm.Expr("saldo + ?", delta)).Error
}

// FindPembayaranForExport returns pembayaran with full relations
func (r *ArusKasRepository) FindPembayaranForExport(dari, sampai string) ([]models.Pembayaran, error) {
	var list []models.Pembayaran
	q := r.db.
		Preload("Kavling").
		Preload("Customer").
		Preload("Bank").
		Order("tanggal DESC")
	if dari != "" {
		q = q.Where("tanggal >= ?", dari)
	}
	if sampai != "" {
		q = q.Where("tanggal <= ?", sampai)
	}
	err := q.Find(&list).Error
	return list, err
}

// ArusKasEntry is a helper type for manual arus kas creation from handler
type ArusKasEntry struct {
	ID            int
	NoTransaksi   string
	Jenis         string
	Kategori      string
	Keterangan    string
	Nominal       float64
	IDBank        *int
	Tanggal       time.Time
	ReferensiID   *int
	ReferensiTipe string
	IDUser        *int
}

func (r *ArusKasRepository) CreateEntry(entry *ArusKasEntry) error {
	t := &models.Transaksi{
		NoTransaksi:   entry.NoTransaksi,
		Jenis:         entry.Jenis,
		Kategori:      entry.Kategori,
		Keterangan:    entry.Keterangan,
		Nominal:       entry.Nominal,
		IDBank:        entry.IDBank,
		Tanggal:       entry.Tanggal,
		ReferensiID:   entry.ReferensiID,
		ReferensiTipe: entry.ReferensiTipe,
		IDUser:        entry.IDUser,
	}
	if err := r.db.Create(t).Error; err != nil {
		return err
	}
	entry.ID = t.ID
	return nil
}

// Placeholder for time usage in repo
var _ = time.Now

// DebugCounts returns counts for diagnostic purposes.
func (r *TransaksiKavlingRepository) DebugCounts() map[string]interface{} {
	out := map[string]interface{}{}
	var n int64

	r.db.Table("customer").Count(&n)
	out["customer_total"] = n

	r.db.Table("customer_kavling").Count(&n)
	out["customer_kavling_total"] = n

	r.db.Table("transaksi_kavling").Count(&n)
	out["transaksi_kavling_total"] = n

	r.db.Table("tagihan").Count(&n)
	out["tagihan_total"] = n

	r.db.Table("kategori_transaksi").Count(&n)
	out["kategori_transaksi_total"] = n

	// customer_kavling tanpa transaksi_kavling
	r.db.Raw(`
		SELECT COUNT(*) FROM customer_kavling ck
		LEFT JOIN transaksi_kavling tk
			ON tk.id_kavling = ck.id_kavling AND tk.id_customer = ck.id_customer
		WHERE tk.id IS NULL
	`).Scan(&n)
	out["customer_kavling_tanpa_transaksi"] = n

	// Sample customer_kavling rows
	var ck []map[string]interface{}
	r.db.Raw(`SELECT * FROM customer_kavling LIMIT 10`).Scan(&ck)
	out["sample_customer_kavling"] = ck

	// Sample transaksi_kavling rows
	var tk []map[string]interface{}
	r.db.Raw(`SELECT id, no_transaksi, id_kavling, id_customer, jenis_pembelian FROM transaksi_kavling LIMIT 10`).Scan(&tk)
	out["sample_transaksi_kavling"] = tk

	// Semua customer — lihat jenis_pembelian & field pembayaran apa adanya
	var allCust []map[string]interface{}
	r.db.Raw(`SELECT id, kode_customer, nama, jenis_pembelian, harga_jual, jumlah_pembayaran
	          FROM customer LIMIT 10`).Scan(&allCust)
	out["sample_customer_all"] = allCust

	// Pembayaran total
	r.db.Table("pembayaran").Count(&n)
	out["pembayaran_total"] = n

	// Diagnostik: BF candidates yang HARUSNYA dibuatkan pembayaran
	var bfCand []map[string]interface{}
	r.db.Raw(`
		SELECT tk.id AS id_transaksi, c.id AS id_customer,
		       c.jenis_pembelian, c.jumlah_pembayaran,
		       (SELECT COUNT(*) FROM pembayaran p WHERE p.id_transaksi = tk.id) AS n_pembayaran
		FROM transaksi_kavling tk
		JOIN customer c ON c.id = tk.id_customer
	`).Scan(&bfCand)
	out["bf_candidates"] = bfCand

	// Cek kategori 001 (Booking Fee)
	var kat map[string]interface{}
	r.db.Raw(`SELECT id, kode, kategori, jenis FROM kategori_transaksi WHERE kode = '001' LIMIT 1`).Scan(&kat)
	out["kategori_bf"] = kat

	return out
}

// NOTE: BackfillFromCustomers, FixTagihanBF, dan BackfillBFPembayaran telah DIHAPUS.
// Semua logika pembuatan TransaksiKavling + Tagihan + Pembayaran BF
// sekarang dilakukan langsung di customer_handler_v2.go → generatePembayaran()
// saat customer disimpan. Tidak ada lagi "self-heal" di endpoint GET.

// resolveHargaKavling loads the kavling's harga_jual_cash via transaksi_kavling.
// Ini SUMBER KEBENARAN untuk harga unit — jangan pakai customer.HargaJual.
func resolveHargaKavling(db *gorm.DB, idTransaksi int) float64 {
	var trx models.TransaksiKavling
	if err := db.Select("id_kavling, harga_jual").First(&trx, idTransaksi).Error; err != nil {
		return 0
	}
	var kav models.Kavling
	if err := db.Select("harga_jual_cash").First(&kav, trx.IDKavling).Error; err != nil {
		return trx.HargaJual
	}
	if kav.HargaJualCash > 0 {
		return kav.HargaJualCash
	}
	return trx.HargaJual
}

func resolveKavlingLabel(db *gorm.DB, idTransaksi int) string {
	var trx models.TransaksiKavling
	if err := db.Select("id_kavling").First(&trx, idTransaksi).Error; err != nil {
		return ""
	}
	var kav models.Kavling
	if err := db.Preload("Lokasi").First(&kav, trx.IDKavling).Error; err != nil {
		return ""
	}
	lokasiNama := ""
	if kav.Lokasi != nil {
		lokasiNama = kav.Lokasi.Nama
	}
	return fmt.Sprintf("Harga Tanah Kavling di %s Blok %s", lokasiNama, kav.KodeKavling)
}

func CreateTagihanForCustomer(idTransaksi int, customer *models.Customer, tagihanRepo *TagihanRepository) {
	db := tagihanRepo.db

	// Guard: jika sudah ada tagihan untuk transaksi ini, skip — mencegah duplikasi
	var existingCount int64
	db.Model(&models.Tagihan{}).Where("id_transaksi = ?", idTransaksi).Count(&existingCount)
	if existingCount > 0 {
		log.Printf("[CreateTagihan] SKIP trx %d — sudah punya %d tagihan", idTransaksi, existingCount)
		return
	}

	ensureKategori := func(kode, kategori string) int {
		id := tagihanRepo.FindKategoriByKode(kode)
		if id > 0 {
			return id
		}
		k := &models.KategoriTransaksi{Kode: kode, Kategori: kategori, Jenis: "PEMASUKAN", IsSystem: true}
		if err := db.Create(k).Error; err != nil {
			log.Printf("[CreateTagihan] gagal auto-seed kategori %s: %v", kode, err)
			return 0
		}
		log.Printf("[CreateTagihan] auto-seeded kategori %s (%s) id=%d", kode, kategori, k.ID)
		return k.ID
	}

	// Harga kavling (sumber kebenaran untuk tagihan)
	hargaKavling := resolveHargaKavling(db, idTransaksi)
	kavlingLabel := resolveKavlingLabel(db, idTransaksi)
	if kavlingLabel == "" {
		kavlingLabel = "Harga Unit Kavling"
	}

	switch strings.ToUpper(strings.TrimSpace(customer.JenisPembelian)) {
	case "BOOKING FEE":
		// Tagihan = harga unit kavling (bukan BF!)
		idKat := ensureKategori("001", "Booking Fee")
		if hargaKavling <= 0 {
			log.Printf("[CreateTagihan] WARNING harga kavling 0 untuk trx %d — skip tagihan BF", idTransaksi)
			return
		}
		tagihan := &models.Tagihan{IDTransaksi: idTransaksi, IDKategori: idKat, Deskripsi: kavlingLabel, Nominal: hargaKavling}
		if err := db.Create(tagihan).Error; err != nil {
			log.Printf("[CreateTagihan] ERROR booking fee tagihan: %v", err)
			return
		}
		// Catat BF yang sudah dibayar sebagai pemasukan pertama
		if customer.JumlahPembayaran > 0 {
			var trxRec models.TransaksiKavling
			db.Select("id_kavling").First(&trxRec, idTransaksi)

			var nextKe int64
			db.Model(&models.Pembayaran{}).Where("id_transaksi = ?", idTransaksi).Count(&nextKe)
			noPay := fmt.Sprintf("PAY-%s-%d-%d", time.Now().Format("20060102"), idTransaksi, nextKe+1)
			pay := &models.Pembayaran{
				NoPembayaran: noPay,
				IDTransaksi:  idTransaksi,
				IDCustomer:   customer.ID,
				IDKavling:    trxRec.IDKavling,
				IDKategori:   &idKat,
				CaraBayar:    "TUNAI",
				Tanggal:      time.Now(),
				PembayaranKe: int(nextKe) + 1,
				JumlahBayar:  customer.JumlahPembayaran,
				Keterangan:   fmt.Sprintf("Booking Fee %s", kavlingLabel),
			}
			if err := db.Create(pay).Error; err != nil {
				log.Printf("[CreateTagihan] ERROR booking fee pemasukan: %v", err)
			}
		}
	case "CASH KERAS":
		idKat := ensureKategori("004", "Pembelian Cash")
		harga := hargaKavling
		if harga <= 0 {
			harga = customer.HargaJual
		}
		if err := db.Create(&models.Tagihan{IDTransaksi: idTransaksi, IDKategori: idKat, Deskripsi: kavlingLabel, Nominal: harga}).Error; err != nil {
			log.Printf("[CreateTagihan] ERROR cash: %v", err)
		}
	case "KREDIT":
		idKat := ensureKategori("002", "Cicilan")
		harga := hargaKavling
		if harga <= 0 {
			harga = customer.HargaJual
		}
		
		// Buat 1 Tagihan utama seharga kavling (bukan per bulan)
		tagihan := &models.Tagihan{
			IDTransaksi: idTransaksi,
			IDKategori:  idKat,
			Deskripsi:   kavlingLabel,
			Nominal:     harga,
		}
		if err := db.Create(tagihan).Error; err != nil {
			log.Printf("[CreateTagihan] ERROR kredit tagihan: %v", err)
		}

		// Jika ada DP (JumlahPembayaran > 0), catat sebagai Pembayaran pertama
		if customer.JumlahPembayaran > 0 {
			var trxRec models.TransaksiKavling
			db.Select("id_kavling").First(&trxRec, idTransaksi)

			var nextKe int64
			db.Model(&models.Pembayaran{}).Where("id_transaksi = ?", idTransaksi).Count(&nextKe)
			noPay := fmt.Sprintf("PAY-%s-%d-%d", time.Now().Format("20060102"), idTransaksi, nextKe+1)
			
			pay := &models.Pembayaran{
				NoPembayaran: noPay,
				IDTransaksi:  idTransaksi,
				IDCustomer:   customer.ID,
				IDKavling:    trxRec.IDKavling,
				IDKategori:   &idKat,
				CaraBayar:    "TUNAI",
				Tanggal:      time.Now(),
				PembayaranKe: int(nextKe) + 1,
				JumlahBayar:  customer.JumlahPembayaran,
				Keterangan:   fmt.Sprintf("DP Kredit %s", kavlingLabel),
			}
			if err := db.Create(pay).Error; err != nil {
				log.Printf("[CreateTagihan] ERROR kredit pemasukan DP: %v", err)
			}
		}
	default:
		log.Printf("[CreateTagihan] WARNING jenis_pembelian %q tidak dikenal", customer.JenisPembelian)
	}
}
