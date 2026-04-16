package models

import "time"

// ─── Users ────────────────────────────────────────────────────────────────────

type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Nama      string    `json:"nama" gorm:"not null"`
	IsAdmin   int       `json:"is_admin" gorm:"default:0"`
	Status    string    `json:"status" gorm:"default:'AKTIF'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string { return "users" }

// ─── Menu & Hak Akses ─────────────────────────────────────────────────────────

type Menu struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Nama      string    `json:"nama"`
	Icon      string    `json:"icon"`
	URL       string    `json:"url"`
	Urutan    int       `json:"urutan"`
	ParentID  *int      `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (Menu) TableName() string { return "menu" }

type HakAkses struct {
	ID        int  `json:"id" gorm:"primaryKey"`
	IDUser    int  `json:"id_user"`
	IDMenu    int  `json:"id_menu"`
	StatusHak int  `json:"status_hak" gorm:"default:1"`
	Menu      Menu `json:"menu,omitempty" gorm:"foreignKey:IDMenu"`
}

func (HakAkses) TableName() string { return "hak_akses" }

// ─── Throttle ─────────────────────────────────────────────────────────────────

type Throttle struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	IPAddress   string    `json:"ip_address"`
	Username    string    `json:"username"`
	Attempts    int       `json:"attempts" gorm:"default:1"`
	LastAttempt time.Time `json:"last_attempt"`
}

func (Throttle) TableName() string { return "throttle" }

// ─── Aktifitas ────────────────────────────────────────────────────────────────

type Aktifitas struct {
	ID         int        `json:"id" gorm:"primaryKey"`
	IDUser     *int       `json:"id_user"`
	Aksi       string     `json:"aksi"`
	Keterangan string     `json:"keterangan"`
	IPAddress  string     `json:"ip_address"`
	CreatedAt  time.Time  `json:"created_at"`
	User       *User      `json:"user,omitempty" gorm:"foreignKey:IDUser"`
}

func (Aktifitas) TableName() string { return "aktifitas" }

// ─── Bank ─────────────────────────────────────────────────────────────────────

type Bank struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	NamaBank     string    `json:"nama_bank"`
	NamaRekening string    `json:"nama_rekening"`
	NoRekening   string    `json:"no_rekening"`
	Saldo        float64   `json:"saldo" gorm:"default:0"`
	IsKas        int       `json:"is_kas" gorm:"default:0"`
	Status       int       `json:"status" gorm:"default:1"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Bank) TableName() string { return "bank" }

// ─── Konfigurasi ──────────────────────────────────────────────────────────────

type Konfigurasi struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	NamaPerusahaan string    `json:"nama_perusahaan"`
	Alamat         string    `json:"alamat"`
	Telepon        string    `json:"telepon"`
	Email          string    `json:"email"`
	Website        string    `json:"website"`
	Logo           string    `json:"logo"`
	TtdDigital     string    `json:"ttd_digital"`
	NamaTtd        string    `json:"nama_ttd"`
	JabatanTtd     string    `json:"jabatan_ttd"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (Konfigurasi) TableName() string { return "konfigurasi" }

// ─── Konfigurasi WA ───────────────────────────────────────────────────────────

type KonfigurasiWA struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	IDDevice  string    `json:"id_device"`
	NoTelp    string    `json:"no_telp"`
	ApiURL    string    `json:"api_url"`
	ApiKey    string    `json:"api_key"`
	IsAktif   int       `json:"is_aktif" gorm:"default:1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (KonfigurasiWA) TableName() string { return "konfigurasi_wa" }

// ─── Template ─────────────────────────────────────────────────────────────────

type Template struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Nama      string    `json:"nama"`
	Tipe      string    `json:"tipe"`
	Isi       string    `json:"isi"`
	IsAktif   int       `json:"is_aktif" gorm:"default:1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Template) TableName() string { return "template" }

// ─── Marketing ────────────────────────────────────────────────────────────────

type Marketing struct {
	ID                 int       `json:"id" gorm:"primaryKey"`
	Nama               string    `json:"nama"`
	NoTelp             string    `json:"no_telp"`
	Alamat             string    `json:"alamat"`
	Email              string    `json:"email"`
	PersentaseKomisi   float64   `json:"persentase_komisi" gorm:"default:0"`
	Status             int       `json:"status" gorm:"default:1"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	// Fields computed (not in DB)
	JumlahOpen   int     `json:"jumlah_open,omitempty" gorm:"-"`
	JumlahClosed int     `json:"jumlah_closed,omitempty" gorm:"-"`
}

func (Marketing) TableName() string { return "marketing" }

// ─── Customer ─────────────────────────────────────────────────────────────────

type Customer struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Nama      string    `json:"nama"`
	NoTelp    string    `json:"no_telp"`
	NoKTP     string    `json:"no_ktp"`
	Alamat    string    `json:"alamat"`
	Pekerjaan string    `json:"pekerjaan"`
	FotoKTP   string    `json:"foto_ktp"`
	FotoKK    string    `json:"foto_kk"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Customer) TableName() string { return "customer" }

// ─── Denah Kavling ────────────────────────────────────────────────────────────

// DenahKavling adalah kumpulan kavling dalam satu peta SVG.
// Satu project bisa memiliki multiple denah (Blok A, Blok B, dll).
type DenahKavling struct {
	ID         int           `json:"id" gorm:"primaryKey"`
	Nama       string        `json:"nama" gorm:"not null"`
	SvgContent string        `json:"svg_content" gorm:"type:text;not null"`
	Viewbox    string        `json:"viewbox"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
	// Relasi
	Kavlings []KavlingPeta `json:"kavlings,omitempty" gorm:"foreignKey:DenahKavlingID"`
}

func (DenahKavling) TableName() string { return "denah_kavling" }

// ─── Kavling Peta ─────────────────────────────────────────────────────────────

// KavlingPeta merepresentasikan satu kavling/tanah dalam sebuah denah.
// kode_kavling dan kode_map bersifat immutable setelah dibuat dari SVG.
type KavlingPeta struct {
	ID              int        `json:"id" gorm:"primaryKey"`
	DenahKavlingID  *int       `json:"denah_kavling_id"`
	KodeKavling     string     `json:"kode_kavling" gorm:"uniqueIndex;not null"`
	KodeMap         string     `json:"kode_map" gorm:"type:text"`    // atribut d dari <path>
	PanjangKanan    float64    `json:"panjang_kanan" gorm:"default:0"`
	PanjangKiri     float64    `json:"panjang_kiri" gorm:"default:0"`
	LebarDepan      float64    `json:"lebar_depan" gorm:"default:0"`
	LebarBelakang   float64    `json:"lebar_belakang" gorm:"default:0"`
	LuasTanah       float64    `json:"luas_tanah" gorm:"default:0"`
	HargaPerMeter   float64    `json:"harga_per_meter" gorm:"default:0"`
	HargaJualCash   float64    `json:"harga_jual_cash" gorm:"default:0"`
	Status          int        `json:"status" gorm:"default:0"`
	// 0=Kosong, 1=Booking, 2=Cash, 3=Kredit
	TglJatuhTempo   *time.Time `json:"tgl_jatuh_tempo"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	// Relasi
	Denah *DenahKavling `json:"denah,omitempty" gorm:"foreignKey:DenahKavlingID"`
}

func (KavlingPeta) TableName() string { return "kavling_peta" }

// StatusLabel returns human-readable status
func (k *KavlingPeta) StatusLabel() string {
	switch k.Status {
	case 0:
		return "Kosong"
	case 1:
		return "Booking"
	case 2:
		return "CASH"
	case 3:
		return "KREDIT"
	}
	return "Unknown"
}

// ─── Transaksi Booking ────────────────────────────────────────────────────────

type TransaksiBooking struct {
	ID             int        `json:"id" gorm:"primaryKey"`
	NoBooking      string     `json:"no_booking" gorm:"uniqueIndex;not null"`
	IDKavling      int        `json:"id_kavling"`
	IDCustomer     int        `json:"id_customer"`
	IDMarketing    *int       `json:"id_marketing"`
	NominalBooking float64    `json:"nominal_booking" gorm:"default:0"`
	TglBooking     time.Time  `json:"tgl_booking"`
	TglExpired     *time.Time `json:"tgl_expired"`
	Keterangan     string     `json:"keterangan"`
	Status         int        `json:"status" gorm:"default:1"`
	IDUser         *int       `json:"id_user"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	// Relations
	Kavling   *KavlingPeta `json:"kavling,omitempty" gorm:"foreignKey:IDKavling"`
	Customer  *Customer    `json:"customer,omitempty" gorm:"foreignKey:IDCustomer"`
	Marketing *Marketing   `json:"marketing,omitempty" gorm:"foreignKey:IDMarketing"`
}

func (TransaksiBooking) TableName() string { return "transaksi_booking" }

// ─── Transaksi Kavling ────────────────────────────────────────────────────────

type TransaksiKavling struct {
	ID               int        `json:"id" gorm:"primaryKey"`
	NoTransaksi      string     `json:"no_transaksi" gorm:"uniqueIndex;not null"`
	IDKavling        int        `json:"id_kavling"`
	IDCustomer       int        `json:"id_customer"`
	IDMarketing      *int       `json:"id_marketing"`
	IDBooking        *int       `json:"id_booking"`
	JenisPembelian   int        `json:"jenis_pembelian"` // 2=CASH, 3=KREDIT
	HargaJual        float64    `json:"harga_jual"`
	UangMuka         float64    `json:"uang_muka" gorm:"default:0"`
	LamaCicilan      int        `json:"lama_cicilan" gorm:"default:0"`
	CicilanPerBulan  float64    `json:"cicilan_per_bulan" gorm:"default:0"`
	TglTransaksi     time.Time  `json:"tgl_transaksi"`
	TglMulaiCicilan  *time.Time `json:"tgl_mulai_cicilan"`
	Keterangan       string     `json:"keterangan"`
	IDUser           *int       `json:"id_user"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	// Relations
	Kavling   *KavlingPeta `json:"kavling,omitempty" gorm:"foreignKey:IDKavling"`
	Customer  *Customer    `json:"customer,omitempty" gorm:"foreignKey:IDCustomer"`
	Marketing *Marketing   `json:"marketing,omitempty" gorm:"foreignKey:IDMarketing"`
	// Computed (not in DB)
	JumlahCicilanTerbayar int     `json:"jumlah_cicilan_terbayar" gorm:"-"`
	TotalTerbayar         float64 `json:"total_terbayar" gorm:"-"`
}

func (TransaksiKavling) TableName() string { return "transaksi_kavling" }

// ─── Pembayaran ───────────────────────────────────────────────────────────────

type Pembayaran struct {
	ID              int       `json:"id" gorm:"primaryKey"`
	NoPembayaran    string    `json:"no_pembayaran" gorm:"uniqueIndex;not null"`
	IDTransaksi     int       `json:"id_transaksi"`
	IDCustomer      int       `json:"id_customer"`
	IDKavling       int       `json:"id_kavling"`
	IDBank          *int      `json:"id_bank"`
	Tanggal         time.Time `json:"tanggal"`
	PembayaranKe    int       `json:"pembayaran_ke" gorm:"default:1"`
	JumlahBayar     float64   `json:"jumlah_bayar"`
	BuktiPembayaran string    `json:"bukti_pembayaran"`
	Keterangan      string    `json:"keterangan"`
	IDUser          *int      `json:"id_user"`
	CreatedAt       time.Time `json:"created_at"`
	// Relations
	Transaksi *TransaksiKavling `json:"transaksi,omitempty" gorm:"foreignKey:IDTransaksi"`
	Customer  *Customer         `json:"customer,omitempty" gorm:"foreignKey:IDCustomer"`
	Kavling   *KavlingPeta      `json:"kavling,omitempty" gorm:"foreignKey:IDKavling"`
	Bank      *Bank             `json:"bank,omitempty" gorm:"foreignKey:IDBank"`
}

func (Pembayaran) TableName() string { return "pembayaran" }

// ─── Transaksi (Arus Kas) ─────────────────────────────────────────────────────

type Transaksi struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	NoTransaksi    string    `json:"no_transaksi" gorm:"uniqueIndex;not null"`
	Jenis          string    `json:"jenis"` // Pemasukan | Pengeluaran
	Kategori       string    `json:"kategori"`
	Keterangan     string    `json:"keterangan"`
	Nominal        float64   `json:"nominal"`
	IDBank         *int      `json:"id_bank"`
	Tanggal        time.Time `json:"tanggal"`
	ReferensiID    *int      `json:"referensi_id"`
	ReferensiTipe  string    `json:"referensi_tipe"`
	IDUser         *int      `json:"id_user"`
	CreatedAt      time.Time `json:"created_at"`
	// Relations
	Bank *Bank `json:"bank,omitempty" gorm:"foreignKey:IDBank"`
}

func (Transaksi) TableName() string { return "transaksi" }

// ─── Response helpers ─────────────────────────────────────────────────────────

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	PerPage int         `json:"per_page"`
}

// ─── RekapKredit (computed struct) ───────────────────────────────────────────

type RekapKredit struct {
	IDKavling        int     `json:"id_kavling"`
	KodeKavling      string  `json:"kode_kavling"`
	IDCustomer       int     `json:"id_customer"`
	NamaCustomer     string  `json:"nama_customer"`
	HargaJual        float64 `json:"harga_jual"`
	UangMuka         float64 `json:"uang_muka"`
	LamaCicilan      int     `json:"lama_cicilan"`
	CicilanPerBulan  float64 `json:"cicilan_per_bulan"`
	TglMulaiCicilan  string  `json:"tgl_mulai_cicilan"`
	BulanBerjalan    int     `json:"bulan_berjalan"`
	JumlahPembayaran int     `json:"jumlah_pembayaran"`
	Tunggakan        int     `json:"tunggakan"`
	NominalTunggakan float64 `json:"nominal_tunggakan"`
	StatusBulanIni   string  `json:"status_bulan_ini"`
	SisaAngsuran     int     `json:"sisa_angsuran"`
	IDTransaksi      int     `json:"id_transaksi"`
}
