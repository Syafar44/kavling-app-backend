package services

import (
	"errors"
	"time"

	"backend-kavling/internal/helpers"
	"backend-kavling/internal/models"
	"backend-kavling/internal/repositories"
)

type TransaksiService struct {
	kavlingRepo   *repositories.KavlingRepository
	bookingRepo   *repositories.BookingRepository
	trxKavRepo    *repositories.TransaksiKavlingRepository
	pembayaranRepo *repositories.PembayaranRepository
	arusKasRepo   *repositories.ArusKasRepository
}

func NewTransaksiService(
	kavlingRepo *repositories.KavlingRepository,
	bookingRepo *repositories.BookingRepository,
	trxKavRepo *repositories.TransaksiKavlingRepository,
	pembayaranRepo *repositories.PembayaranRepository,
	arusKasRepo *repositories.ArusKasRepository,
) *TransaksiService {
	return &TransaksiService{
		kavlingRepo:    kavlingRepo,
		bookingRepo:    bookingRepo,
		trxKavRepo:     trxKavRepo,
		pembayaranRepo: pembayaranRepo,
		arusKasRepo:    arusKasRepo,
	}
}

// ─── Booking ──────────────────────────────────────────────────────────────────

type BookingInput struct {
	IDKavling      int     `json:"id_kavling" binding:"required"`
	IDCustomer     int     `json:"id_customer" binding:"required"`
	IDMarketing    *int    `json:"id_marketing"`
	NominalBooking float64 `json:"nominal_booking"`
	TglExpired     string  `json:"tgl_expired"`
	Keterangan     string  `json:"keterangan"`
	IDUser         *int    `json:"id_user"`
}

func (s *TransaksiService) CreateBooking(input BookingInput) (*models.TransaksiBooking, error) {
	// Validasi kavling status 0
	kavling, err := s.kavlingRepo.FindByID(input.IDKavling)
	if err != nil {
		return nil, errors.New("kavling tidak ditemukan")
	}
	if kavling.Status != 0 {
		return nil, errors.New("kavling tidak tersedia (status bukan Kosong)")
	}

	var tglExpired *time.Time
	if input.TglExpired != "" {
		t, err := time.Parse("2006-01-02", input.TglExpired)
		if err != nil {
			return nil, errors.New("format tgl_expired tidak valid (gunakan YYYY-MM-DD)")
		}
		tglExpired = &t
	}

	booking := &models.TransaksiBooking{
		NoBooking:      helpers.GenerateNoBooking(),
		IDKavling:      input.IDKavling,
		IDCustomer:     input.IDCustomer,
		IDMarketing:    input.IDMarketing,
		NominalBooking: input.NominalBooking,
		TglBooking:     time.Now(),
		TglExpired:     tglExpired,
		Keterangan:     input.Keterangan,
		Status:         1,
		IDUser:         input.IDUser,
	}

	if err := s.bookingRepo.Create(booking); err != nil {
		return nil, err
	}

	// Update status kavling → 1 (Booking)
	if err := s.kavlingRepo.UpdateStatus(input.IDKavling, 1); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *TransaksiService) CancelBooking(id int) error {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return errors.New("booking tidak ditemukan")
	}
	if booking.Status != 1 {
		return errors.New("booking sudah tidak aktif")
	}

	booking.Status = 0
	if err := s.bookingRepo.Update(booking); err != nil {
		return err
	}

	// Kembalikan status kavling → 0 (Kosong)
	return s.kavlingRepo.UpdateStatus(booking.IDKavling, 0)
}

// ─── Convert Booking ──────────────────────────────────────────────────────────

type ConvertInput struct {
	Jenis          string  `json:"jenis" binding:"required"` // cash | kredit
	HargaJual      float64 `json:"harga_jual" binding:"required"`
	UangMuka       float64 `json:"uang_muka"`
	LamaCicilan    int     `json:"lama_cicilan"`
	TglMulai       string  `json:"tgl_mulai_cicilan"`
	Keterangan     string  `json:"keterangan"`
	IDUser         *int    `json:"id_user"`
}

func (s *TransaksiService) ConvertBooking(bookingID int, input ConvertInput) (*models.TransaksiKavling, error) {
	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return nil, errors.New("booking tidak ditemukan")
	}
	if booking.Status != 1 {
		return nil, errors.New("booking sudah tidak aktif")
	}

	return s.createTransaksi(booking.IDKavling, booking.IDCustomer, booking.IDMarketing,
		&bookingID, input)
}

// ─── Cash Langsung ────────────────────────────────────────────────────────────

type CashInput struct {
	IDKavling   int     `json:"id_kavling" binding:"required"`
	IDCustomer  int     `json:"id_customer" binding:"required"`
	IDMarketing *int    `json:"id_marketing"`
	HargaJual   float64 `json:"harga_jual" binding:"required"`
	Keterangan  string  `json:"keterangan"`
	IDUser      *int    `json:"id_user"`
}

func (s *TransaksiService) CreateCash(input CashInput) (*models.TransaksiKavling, error) {
	kavling, err := s.kavlingRepo.FindByID(input.IDKavling)
	if err != nil {
		return nil, errors.New("kavling tidak ditemukan")
	}
	if kavling.Status != 0 {
		return nil, errors.New("kavling tidak tersedia")
	}

	return s.createTransaksi(input.IDKavling, input.IDCustomer, input.IDMarketing, nil,
		ConvertInput{
			Jenis:      "cash",
			HargaJual:  input.HargaJual,
			Keterangan: input.Keterangan,
			IDUser:     input.IDUser,
		})
}

// ─── Kredit Langsung ──────────────────────────────────────────────────────────

type KreditInput struct {
	IDKavling   int     `json:"id_kavling" binding:"required"`
	IDCustomer  int     `json:"id_customer" binding:"required"`
	IDMarketing *int    `json:"id_marketing"`
	HargaJual   float64 `json:"harga_jual" binding:"required"`
	UangMuka    float64 `json:"uang_muka"`
	LamaCicilan int     `json:"lama_cicilan" binding:"required,min=1"`
	TglMulai    string  `json:"tgl_mulai_cicilan"`
	Keterangan  string  `json:"keterangan"`
	IDUser      *int    `json:"id_user"`
}

func (s *TransaksiService) CreateKredit(input KreditInput) (*models.TransaksiKavling, error) {
	kavling, err := s.kavlingRepo.FindByID(input.IDKavling)
	if err != nil {
		return nil, errors.New("kavling tidak ditemukan")
	}
	if kavling.Status != 0 {
		return nil, errors.New("kavling tidak tersedia")
	}

	return s.createTransaksi(input.IDKavling, input.IDCustomer, input.IDMarketing, nil,
		ConvertInput{
			Jenis:       "kredit",
			HargaJual:   input.HargaJual,
			UangMuka:    input.UangMuka,
			LamaCicilan: input.LamaCicilan,
			TglMulai:    input.TglMulai,
			Keterangan:  input.Keterangan,
			IDUser:      input.IDUser,
		})
}

// ─── Core createTransaksi ─────────────────────────────────────────────────────

func (s *TransaksiService) createTransaksi(
	idKavling, idCustomer int,
	idMarketing *int,
	idBooking *int,
	input ConvertInput,
) (*models.TransaksiKavling, error) {
	jenis := input.Jenis
	jenisPembelian := 2 // CASH
	kavlingStatus := 2
	var cicilan float64
	var tglMulai *time.Time

	if jenis == "kredit" {
		jenisPembelian = 3
		kavlingStatus = 3
		if input.LamaCicilan <= 0 {
			return nil, errors.New("lama cicilan harus lebih dari 0 untuk pembelian kredit")
		}
		sisa := input.HargaJual - input.UangMuka
		cicilan = sisa / float64(input.LamaCicilan)

		if input.TglMulai != "" {
			t, err := time.Parse("2006-01-02", input.TglMulai)
			if err != nil {
				return nil, errors.New("format tgl_mulai_cicilan tidak valid")
			}
			tglMulai = &t
		} else {
			// Default: bulan depan tanggal 10
			next := time.Now().AddDate(0, 1, 0)
			t := time.Date(next.Year(), next.Month(), 10, 0, 0, 0, 0, time.Local)
			tglMulai = &t
		}
	}

	trx := &models.TransaksiKavling{
		NoTransaksi:     helpers.GenerateNoTransaksi("TRX-KAV"),
		IDKavling:       idKavling,
		IDCustomer:      idCustomer,
		IDMarketing:     idMarketing,
		IDBooking:       idBooking,
		JenisPembelian:  jenisPembelian,
		HargaJual:       input.HargaJual,
		UangMuka:        input.UangMuka,
		LamaCicilan:     input.LamaCicilan,
		CicilanPerBulan: cicilan,
		TglTransaksi:    time.Now(),
		TglMulaiCicilan: tglMulai,
		Keterangan:      input.Keterangan,
		IDUser:          input.IDUser,
	}

	if err := s.trxKavRepo.Create(trx); err != nil {
		return nil, err
	}

	// Update status kavling
	if err := s.kavlingRepo.UpdateStatus(idKavling, kavlingStatus); err != nil {
		return nil, err
	}

	// Set jatuh tempo jika kredit
	if jenis == "kredit" && tglMulai != nil {
		s.kavlingRepo.UpdateTglJatuhTempo(idKavling, tglMulai)
	}

	// Jika dari booking, update status booking
	if idBooking != nil {
		booking, _ := s.bookingRepo.FindByID(*idBooking)
		if booking != nil {
			booking.Status = 2 // dikonversi
			s.bookingRepo.Update(booking)
		}
	}

	return trx, nil
}

// ─── Pembayaran Cicilan ───────────────────────────────────────────────────────

type BayarInput struct {
	IDKavling       int     `json:"id_kavling" binding:"required"`
	IDBank          *int    `json:"id_bank"`
	Tanggal         string  `json:"tanggal"`
	JumlahBayar     float64 `json:"jumlah_bayar" binding:"required"`
	BuktiPembayaran string  `json:"bukti_pembayaran"`
	Keterangan      string  `json:"keterangan"`
	IDUser          *int    `json:"id_user"`
}

func (s *TransaksiService) BayarCicilan(input BayarInput) (*models.Pembayaran, error) {
	// Cari transaksi kavling
	trx, err := s.trxKavRepo.FindByKavling(input.IDKavling)
	if err != nil {
		return nil, errors.New("transaksi kavling tidak ditemukan")
	}
	if trx.JenisPembelian != 3 {
		return nil, errors.New("kavling ini bukan pembelian kredit")
	}

	// Hitung pembayaran ke berapa
	count := s.pembayaranRepo.CountByKavling(input.IDKavling)
	pembayaranKe := count + 1

	tanggal := time.Now()
	if input.Tanggal != "" {
		t, err := time.Parse("2006-01-02", input.Tanggal)
		if err != nil {
			return nil, errors.New("format tanggal tidak valid")
		}
		tanggal = t
	}

	pem := &models.Pembayaran{
		NoPembayaran:    helpers.GenerateNoPembayaran(),
		IDTransaksi:     trx.ID,
		IDCustomer:      trx.IDCustomer,
		IDKavling:       input.IDKavling,
		IDBank:          input.IDBank,
		Tanggal:         tanggal,
		PembayaranKe:    pembayaranKe,
		JumlahBayar:     input.JumlahBayar,
		BuktiPembayaran: input.BuktiPembayaran,
		Keterangan:      input.Keterangan,
		IDUser:          input.IDUser,
	}

	if err := s.pembayaranRepo.Create(pem); err != nil {
		return nil, err
	}

	// Otomatis catat ke arus kas
	pemID := pem.ID
	arusKas := &models.Transaksi{
		NoTransaksi:   helpers.GenerateNoArusKas(),
		Jenis:         "Pemasukan",
		Kategori:      "Bayar Angsuran",
		Keterangan:    "Angsuran ke-" + string(rune('0'+pembayaranKe)) + " kavling " + trx.Kavling.KodeKavling,
		Nominal:       input.JumlahBayar,
		IDBank:        input.IDBank,
		Tanggal:       tanggal,
		ReferensiID:   &pemID,
		ReferensiTipe: "pembayaran",
		IDUser:        input.IDUser,
	}
	if err := s.arusKasRepo.Create(arusKas); err != nil {
		return nil, err
	}

	// Update saldo bank
	if input.IDBank != nil {
		_ = s.arusKasRepo.UpdateBankSaldo(*input.IDBank, input.JumlahBayar)
	}

	// Cek apakah sudah lunas
	newCount := s.pembayaranRepo.CountByKavling(input.IDKavling)
	if newCount >= trx.LamaCicilan {
		// Update kavling status → 2 (CASH/Lunas)
		_ = s.kavlingRepo.UpdateStatus(input.IDKavling, 2)
	}

	return pem, nil
}

// ─── Rekap Kredit ─────────────────────────────────────────────────────────────

func (s *TransaksiService) GetRekapKredit() ([]models.RekapKredit, error) {
	trxList, err := s.trxKavRepo.FindAllKredit()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	result := make([]models.RekapKredit, 0, len(trxList))

	for _, trx := range trxList {
		if trx.TglMulaiCicilan == nil {
			continue
		}

		bulanBerjalan := helpers.MonthDiff(*trx.TglMulaiCicilan, now)
		if bulanBerjalan < 0 {
			bulanBerjalan = 0
		}

		jumlahBayar := s.pembayaranRepo.CountByKavling(trx.IDKavling)
		tunggakan := bulanBerjalan - jumlahBayar - 1
		if tunggakan < 0 {
			tunggakan = 0
		}

		statusBulanIni := "Belum Bayar"
		if s.pembayaranRepo.HasPaidThisMonth(trx.IDKavling, now.Year(), int(now.Month())) {
			statusBulanIni = "Sudah Bayar"
		}

		namaCustomer := ""
		if trx.Customer != nil {
			namaCustomer = trx.Customer.Nama
		}
		kodeKavling := ""
		if trx.Kavling != nil {
			kodeKavling = trx.Kavling.KodeKavling
		}

		sisa := trx.LamaCicilan - jumlahBayar
		if sisa < 0 {
			sisa = 0
		}

		result = append(result, models.RekapKredit{
			IDKavling:        trx.IDKavling,
			KodeKavling:      kodeKavling,
			IDCustomer:       trx.IDCustomer,
			NamaCustomer:     namaCustomer,
			HargaJual:        trx.HargaJual,
			UangMuka:         trx.UangMuka,
			LamaCicilan:      trx.LamaCicilan,
			CicilanPerBulan:  trx.CicilanPerBulan,
			TglMulaiCicilan:  trx.TglMulaiCicilan.Format("2006-01-02"),
			BulanBerjalan:    bulanBerjalan,
			JumlahPembayaran: jumlahBayar,
			Tunggakan:        tunggakan,
			NominalTunggakan: float64(tunggakan) * trx.CicilanPerBulan,
			StatusBulanIni:   statusBulanIni,
			SisaAngsuran:     sisa,
			IDTransaksi:      trx.ID,
		})
	}

	return result, nil
}
