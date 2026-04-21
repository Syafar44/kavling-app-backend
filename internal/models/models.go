package models

import "time"

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

// ─── Users ────────────────────────────────────────────────────────────────────

type User struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Username    string    `json:"username" gorm:"uniqueIndex;not null"`
	Password    string    `json:"-" gorm:"not null"`
	Nama        string    `json:"nama" gorm:"not null"`
	NamaLengkap string    `json:"nama_lengkap"`
	Email       string    `json:"email"`
	Role        string    `json:"role" gorm:"default:'Admin'"`
	IsAdmin     int       `json:"is_admin" gorm:"default:0"`
	Status      string    `json:"status" gorm:"default:'AKTIF'"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (User) TableName() string { return "users" }

// ─── Menu & Hak Akses ─────────────────────────────────────────────────────────

type Menu struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Nama     string `json:"nama"`
	Icon     string `json:"icon"`
	URL      string `json:"url"`
	Urutan   int    `json:"urutan"`
	ParentID *int   `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (Menu) TableName() string { return "menu" }

type HakAkses struct {
	ID        int     `json:"id" gorm:"primaryKey"`
	IDUser    int     `json:"id_user"`
	IDMenu    int     `json:"id_menu"`
	StatusHak int     `json:"status_hak" gorm:"default:1"`
	Lihat     bool    `json:"lihat" gorm:"default:true"`
	Beranda   bool    `json:"beranda" gorm:"default:false"`
	Tambah    bool    `json:"tambah" gorm:"default:false"`
	Edit      bool    `json:"edit" gorm:"default:false"`
	Hapus     bool    `json:"hapus" gorm:"default:false"`
	Menu      Menu    `json:"menu,omitempty" gorm:"foreignKey:IDMenu"`
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
	ID         int       `json:"id" gorm:"primaryKey"`
	IDUser     *int      `json:"id_user"`
	Aksi       string    `json:"aksi"`
	Keterangan string    `json:"keterangan"`
	IPAddress  string    `json:"ip_address"`
	CreatedAt  time.Time `json:"created_at"`
	User       *User     `json:"user,omitempty" gorm:"foreignKey:IDUser"`
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
	NPWP           string    `json:"npwp"`
	Deskripsi      string    `json:"deskripsi"`
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
	ID               int       `json:"id" gorm:"primaryKey"`
	KodeMarketing    string    `json:"kode_marketing"`
	Nama             string    `json:"nama"`
	NoTelp           string    `json:"no_telp"`
	Alamat           string    `json:"alamat"`
	JenisKelamin     string    `json:"jenis_kelamin"`
	Email            string    `json:"email"`
	PersentaseKomisi float64   `json:"persentase_komisi" gorm:"default:0"`
	Foto             string    `json:"foto"`
	Status           int       `json:"status" gorm:"default:1"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	// Computed
	JumlahOpen   int `json:"jumlah_open,omitempty" gorm:"-"`
	JumlahClosed int `json:"jumlah_closed,omitempty" gorm:"-"`
}

func (Marketing) TableName() string { return "marketing" }

// ─── Notaris ──────────────────────────────────────────────────────────────────

type Notaris struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	Nama       string    `json:"nama"`
	Alamat     string    `json:"alamat"`
	NoTelp     string    `json:"no_telp"`
	Keterangan string    `json:"keterangan" gorm:"default:'NOTARIS'"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Notaris) TableName() string { return "notaris" }

// ─── Lokasi Kavling ───────────────────────────────────────────────────────────

type LokasiKavling struct {
	ID                   int       `json:"id" gorm:"primaryKey"`
	Nama                 string    `json:"nama"`
	NamaSingkat          string    `json:"nama_singkat"`
	Header               string    `json:"header"`
	Alamat               string    `json:"alamat"`
	NamaPerusahaan       string    `json:"nama_perusahaan"`
	NamaAdmin            string    `json:"nama_admin"`
	NamaMengetahui       string    `json:"nama_mengetahui"`
	AlamatPerusahaan     string    `json:"alamat_perusahaan"`
	TelpPerusahaan       string    `json:"telp_perusahaan"`
	KotaPenandatangan    string    `json:"kota_penandatangan"`
	NamaPenandatangan    string    `json:"nama_penandatangan"`
	JabatanPenandatangan string    `json:"jabatan_penandatangan"`
	JenisPembelian       string    `json:"jenis_pembelian" gorm:"default:'cash_kredit'"`
	UrutanLokasi         int       `json:"urutan_lokasi" gorm:"default:1"`
	JumlahKavling        int       `json:"jumlah_kavling" gorm:"default:0"`
	KopSurat             string    `json:"kop_surat"`
	Kwitansi             string    `json:"kwitansi"`
	FotoKavling          string    `json:"foto_kavling"`
	SvgContent           string    `json:"svg_content"`
	SvgLayoutJson        *string   `json:"svg_layout_json" gorm:"type:jsonb"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func (LokasiKavling) TableName() string { return "lokasi_kavling" }

// ─── Kavling ──────────────────────────────────────────────────────────────────

type Kavling struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	IDLokasi       int       `json:"id_lokasi"`
	KodeKavling    string    `json:"kode_kavling"`
	PanjangKanan   float64   `json:"panjang_kanan" gorm:"default:0"`
	PanjangKiri    float64   `json:"panjang_kiri" gorm:"default:0"`
	LebarDepan     float64   `json:"lebar_depan" gorm:"default:0"`
	LebarBelakang  float64   `json:"lebar_belakang" gorm:"default:0"`
	LuasTanah      float64   `json:"luas_tanah" gorm:"default:0"`
	HargaPerMeter  float64   `json:"harga_per_meter" gorm:"default:0"`
	HargaJualCash  float64   `json:"harga_jual_cash" gorm:"default:0"`
	NoSertipikat   string    `json:"no_sertipikat"`
	Keterangan     string    `json:"keterangan"`
	Status         int       `json:"status" gorm:"default:0"`
	IDCustomer     *int      `json:"id_customer"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	// Relations
	Lokasi   *LokasiKavling `json:"lokasi,omitempty" gorm:"foreignKey:IDLokasi"`
	Customer *Customer      `json:"customer,omitempty" gorm:"foreignKey:IDCustomer"`
}

func (Kavling) TableName() string { return "kavling" }

func (k *Kavling) StatusLabel() string {
	labels := []string{"Ready", "HOLD", "BF", "AKAD", "User Cancel", "LUNAS"}
	if k.Status >= 0 && k.Status < len(labels) {
		return labels[k.Status]
	}
	return "Unknown"
}

// ─── Customer ─────────────────────────────────────────────────────────────────

type Customer struct {
	ID                 int        `json:"id" gorm:"primaryKey"`
	KodeKontrak        string     `json:"kode_kontrak"`
	Tanggal            *string    `json:"tanggal" gorm:"type:date"`
	Nama               string     `json:"nama"`
	NoTelp             string     `json:"no_telp"`
	NoKTP              string     `json:"no_ktp"`
	NoKTPPasangan      string     `json:"no_ktp_pasangan"`
	TempatLahir        string     `json:"tempat_lahir"`
	TanggalLahir       *string    `json:"tanggal_lahir" gorm:"type:date"`
	JenisKelamin       string     `json:"jenis_kelamin"`
	Email              string     `json:"email"`
	NPWP               string     `json:"npwp"`
	Alamat             string     `json:"alamat"`
	AlamatDomisili     string     `json:"alamat_domisili"`
	Pekerjaan          string     `json:"pekerjaan"`
	IDLokasi           *int       `json:"id_lokasi"`
	HargaJual          float64    `json:"harga_jual" gorm:"default:0"`
	IDMarketing        *int       `json:"id_marketing"`
	StatusPenjualan    string     `json:"status_penjualan" gorm:"default:'BF'"`
	KeteranganCashback string     `json:"keterangan_cashback"`
	IDBankPembayaran   *int       `json:"id_bank_pembayaran"`
	JenisPembelian      string     `json:"jenis_pembelian" gorm:"default:'KREDIT'"`
	JumlahPembayaran    float64    `json:"jumlah_pembayaran" gorm:"default:0"`
	TanggalBatasBooking *string    `json:"tanggal_batas_booking" gorm:"type:date"`
	CicilanPerBulan     *float64   `json:"cicilan_per_bulan"`
	Tenor               *int       `json:"tenor"`
	JatuhTempo          *string    `json:"jatuh_tempo" gorm:"type:date"`
	FotoKTP            string     `json:"foto_ktp"`
	FotoKK             string     `json:"foto_kk"`
	FotoPemohon        string     `json:"foto_pemohon"`
	FotoKTPPasangan    string     `json:"foto_ktp_pasangan"`
	FotoNPWP           string     `json:"foto_npwp"`
	FotoBPJS           string     `json:"foto_bpjs"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	// Relations
	Lokasi    *LokasiKavling `json:"lokasi,omitempty" gorm:"foreignKey:IDLokasi"`
	Marketing *Marketing     `json:"marketing,omitempty" gorm:"foreignKey:IDMarketing"`
}

func (Customer) TableName() string { return "customer" }

// ─── Customer Kavling (junction) ──────────────────────────────────────────────

type CustomerKavling struct {
	ID         int     `json:"id" gorm:"primaryKey"`
	IDCustomer int     `json:"id_customer"`
	IDKavling  int     `json:"id_kavling"`
	Kavling    Kavling `json:"kavling,omitempty" gorm:"foreignKey:IDKavling"`
}

func (CustomerKavling) TableName() string { return "customer_kavling" }

// ─── Customer File ────────────────────────────────────────────────────────────

type CustomerFile struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	IDCustomer int       `json:"id_customer"`
	Tanggal    time.Time `json:"tanggal"`
	NamaFile   string    `json:"nama_file"`
	PathFile   string    `json:"path_file"`
	Keterangan string    `json:"keterangan"`
	CreatedAt  time.Time `json:"created_at"`
}

func (CustomerFile) TableName() string { return "customer_file" }

// ─── Customer Arsip ───────────────────────────────────────────────────────────

type CustomerArsip struct {
	ID            int       `json:"id" gorm:"primaryKey"`
	IDCustomer    *int      `json:"id_customer"`
	Tanggal       time.Time `json:"tanggal"`
	Nama          string    `json:"nama"`
	IDMarketing   *int      `json:"id_marketing"`
	Lokasi        string    `json:"lokasi"`
	StatusProgres string    `json:"status_progres" gorm:"default:'Aktif'"`
	CreatedAt     time.Time `json:"created_at"`
	// Relations
	Marketing *Marketing `json:"marketing,omitempty" gorm:"foreignKey:IDMarketing"`
}

func (CustomerArsip) TableName() string { return "customer_arsip" }

// ─── Prospek ──────────────────────────────────────────────────────────────────

type Prospek struct {
	ID               int       `json:"id" gorm:"primaryKey"`
	Tanggal          time.Time `json:"tanggal"`
	Nama             string    `json:"nama"`
	NoTelp           string    `json:"no_telp"`
	Usia             *int      `json:"usia"`
	Pekerjaan        string    `json:"pekerjaan"`
	Penghasilan      *float64  `json:"penghasilan"`
	SumberInformasi  string    `json:"sumber_informasi"`
	Rangking         string    `json:"rangking" gorm:"default:'C'"`
	IDMarketing      *int      `json:"id_marketing"`
	Keterangan       string    `json:"keterangan"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	// Relations
	Marketing *Marketing `json:"marketing,omitempty" gorm:"foreignKey:IDMarketing"`
}

func (Prospek) TableName() string { return "prospek" }

// ─── Kategori Transaksi ───────────────────────────────────────────────────────

type KategoriTransaksi struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Kode      string    `json:"kode" gorm:"uniqueIndex"`
	Kategori  string    `json:"kategori"`
	Jenis     string    `json:"jenis"` // PEMASUKAN | PENGELUARAN
	IsSystem  bool      `json:"is_system" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (KategoriTransaksi) TableName() string { return "kategori_transaksi" }

// ─── List Penjualan ───────────────────────────────────────────────────────────

type ListPenjualan struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	Kode       string    `json:"kode"`
	Progres    string    `json:"progres"`
	Keterangan string    `json:"keterangan"`
	Status     string    `json:"status" gorm:"default:'Aktif'"`
	Warna      string    `json:"warna" gorm:"default:'#22c55e'"`
	Urutan     int       `json:"urutan" gorm:"default:1"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (ListPenjualan) TableName() string { return "list_penjualan" }

// ─── Tagihan (invoice items per transaksi) ────────────────────────────────────

type Tagihan struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	IDTransaksi int       `json:"id_transaksi"`
	IDKategori  int       `json:"id_kategori"`
	Deskripsi   string    `json:"deskripsi"`
	Nominal     float64   `json:"nominal" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// Relations
	Kategori *KategoriTransaksi `json:"kategori,omitempty" gorm:"foreignKey:IDKategori"`
}

func (Tagihan) TableName() string { return "tagihan" }

// ─── Legalitas ────────────────────────────────────────────────────────────────

type Legalitas struct {
	ID            int       `json:"id" gorm:"primaryKey"`
	IDKavling     int       `json:"id_kavling" gorm:"uniqueIndex"`
	AtasNamaSurat string    `json:"atas_nama_surat"`
	NoSurat       string    `json:"no_surat"`
	IDNotaris     *int      `json:"id_notaris"`
	Progres       string    `json:"progres" gorm:"default:'BF'"`
	BuktiFoto     string    `json:"bukti_foto"`
	Keterangan    string    `json:"keterangan"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	// Relations
	Kavling  *Kavling  `json:"kavling,omitempty" gorm:"foreignKey:IDKavling"`
	Notaris  *Notaris  `json:"notaris,omitempty" gorm:"foreignKey:IDNotaris"`
}

func (Legalitas) TableName() string { return "legalitas" }

// ─── Hutang ───────────────────────────────────────────────────────────────────

type Hutang struct {
	ID               int        `json:"id" gorm:"primaryKey"`
	Tanggal          time.Time  `json:"tanggal"`
	Deskripsi        string     `json:"deskripsi"`
	Nominal          float64    `json:"nominal" gorm:"default:0"`
	Status           string     `json:"status" gorm:"default:'Belum Lunas'"`
	TanggalPelunasan *time.Time `json:"tanggal_pelunasan"`
	IDBank           *int       `json:"id_bank"`
	Lampiran         string     `json:"lampiran"`
	SisaBayar        float64    `json:"sisa_bayar" gorm:"default:0"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	// Relations
	Bank *Bank `json:"bank,omitempty" gorm:"foreignKey:IDBank"`
}

func (Hutang) TableName() string { return "hutang" }

// ─── Piutang ──────────────────────────────────────────────────────────────────

type Piutang struct {
	ID               int        `json:"id" gorm:"primaryKey"`
	Tanggal          time.Time  `json:"tanggal"`
	Deskripsi        string     `json:"deskripsi"`
	Nominal          float64    `json:"nominal" gorm:"default:0"`
	Status           string     `json:"status" gorm:"default:'Belum Lunas'"`
	TanggalPelunasan *time.Time `json:"tanggal_pelunasan"`
	IDBank           *int       `json:"id_bank"`
	Lampiran         string     `json:"lampiran"`
	SisaBayar        float64    `json:"sisa_bayar" gorm:"default:0"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	// Relations
	Bank *Bank `json:"bank,omitempty" gorm:"foreignKey:IDBank"`
}

func (Piutang) TableName() string { return "piutang" }

// ─── Mutasi Saldo ─────────────────────────────────────────────────────────────

type MutasiSaldo struct {
	ID            int       `json:"id" gorm:"primaryKey"`
	Tanggal       time.Time `json:"tanggal"`
	IDBankAsal    int       `json:"id_bank_asal"`
	IDBankTujuan  int       `json:"id_bank_tujuan"`
	Nominal       float64   `json:"nominal" gorm:"default:0"`
	Keterangan    string    `json:"keterangan"`
	Lampiran      string    `json:"lampiran"`
	CreatedAt     time.Time `json:"created_at"`
	// Relations
	BankAsal   *Bank `json:"bank_asal,omitempty" gorm:"foreignKey:IDBankAsal"`
	BankTujuan *Bank `json:"bank_tujuan,omitempty" gorm:"foreignKey:IDBankTujuan"`
}

func (MutasiSaldo) TableName() string { return "mutasi_saldo" }

// ─── Media ────────────────────────────────────────────────────────────────────

type Media struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	JenisData  string    `json:"jenis_data" gorm:"uniqueIndex"`
	Keterangan string    `json:"keterangan"`
	NamaFile   string    `json:"nama_file"`
	PathFile   string    `json:"path_file"`
	IsLogo     bool      `json:"is_logo" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Media) TableName() string { return "media" }

// ─── Landing Konten ───────────────────────────────────────────────────────────

type LandingKonten struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Item      string    `json:"item"`
	Judul     string    `json:"judul"`
	Artikel   string    `json:"artikel"`
	Gambar    string    `json:"gambar"`
	Icon      string    `json:"icon"`
	Urutan    int       `json:"urutan" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (LandingKonten) TableName() string { return "landing_konten" }

// ─── KavlingPeta (legacy — maps to kavling_peta table, kept for backward compat) ──

type KavlingPeta struct {
	ID             int        `json:"id" gorm:"primaryKey"`
	DenahKavlingID *int       `json:"denah_kavling_id"`
	KodeKavling    string     `json:"kode_kavling" gorm:"uniqueIndex;not null"`
	KodeMap        string     `json:"kode_map" gorm:"type:text"`
	PanjangKanan   float64    `json:"panjang_kanan" gorm:"default:0"`
	PanjangKiri    float64    `json:"panjang_kiri" gorm:"default:0"`
	LebarDepan     float64    `json:"lebar_depan" gorm:"default:0"`
	LebarBelakang  float64    `json:"lebar_belakang" gorm:"default:0"`
	LuasTanah      float64    `json:"luas_tanah" gorm:"default:0"`
	HargaPerMeter  float64    `json:"harga_per_meter" gorm:"default:0"`
	HargaJualCash  float64    `json:"harga_jual_cash" gorm:"default:0"`
	Status         int        `json:"status" gorm:"default:0"`
	TglJatuhTempo  *time.Time `json:"tgl_jatuh_tempo"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	Denah          *DenahKavling `json:"denah,omitempty" gorm:"foreignKey:DenahKavlingID"`
}

func (KavlingPeta) TableName() string { return "kavling_peta" }

func (k *KavlingPeta) StatusLabel() string {
	switch k.Status {
	case 0: return "Kosong"
	case 1: return "Booking"
	case 2: return "CASH"
	case 3: return "KREDIT"
	}
	return "Unknown"
}

// ─── Denah Kavling (legacy — kept for backward compat) ────────────────────────

type DenahKavling struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Nama      string    `json:"nama" gorm:"not null"`
	SvgContent string   `json:"svg_content" gorm:"type:text;not null"`
	Viewbox   string    `json:"viewbox"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (DenahKavling) TableName() string { return "denah_kavling" }

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
	Customer  *Customer  `json:"customer,omitempty" gorm:"foreignKey:IDCustomer"`
	Marketing *Marketing `json:"marketing,omitempty" gorm:"foreignKey:IDMarketing"`
}

func (TransaksiBooking) TableName() string { return "transaksi_booking" }

// ─── Transaksi Kavling ────────────────────────────────────────────────────────

type TransaksiKavling struct {
	ID              int        `json:"id" gorm:"primaryKey"`
	NoTransaksi     string     `json:"no_transaksi" gorm:"uniqueIndex;not null"`
	IDKavling       int        `json:"id_kavling"`
	IDCustomer      int        `json:"id_customer"`
	IDMarketing     *int       `json:"id_marketing"`
	IDBooking       *int       `json:"id_booking"`
	JenisPembelian  int        `json:"jenis_pembelian"` // 2=CASH, 3=KREDIT
	HargaJual       float64    `json:"harga_jual"`
	UangMuka        float64    `json:"uang_muka" gorm:"default:0"`
	LamaCicilan     int        `json:"lama_cicilan" gorm:"default:0"`
	CicilanPerBulan float64    `json:"cicilan_per_bulan" gorm:"default:0"`
	TglTransaksi    time.Time  `json:"tgl_transaksi"`
	TglMulaiCicilan *time.Time `json:"tgl_mulai_cicilan"`
	Keterangan      string     `json:"keterangan"`
	IDUser          *int       `json:"id_user"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	// Relations
	Kavling   *Kavling   `json:"kavling,omitempty" gorm:"foreignKey:IDKavling"`
	Customer  *Customer  `json:"customer,omitempty" gorm:"foreignKey:IDCustomer"`
	Marketing *Marketing `json:"marketing,omitempty" gorm:"foreignKey:IDMarketing"`
	// Computed
	JumlahCicilanTerbayar int     `json:"jumlah_cicilan_terbayar" gorm:"-"`
	TotalTerbayar         float64 `json:"total_terbayar" gorm:"-"`
	TotalTagihan          float64 `json:"total_tagihan" gorm:"-"`
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
	IDKategori      *int      `json:"id_kategori"`
	CaraBayar       string    `json:"cara_bayar"`
	Tanggal         time.Time `json:"tanggal"`
	PembayaranKe    int       `json:"pembayaran_ke" gorm:"default:1"`
	JumlahBayar     float64   `json:"jumlah_bayar"`
	BuktiPembayaran string    `json:"bukti_pembayaran"`
	Keterangan      string    `json:"keterangan"`
	IDUser          *int      `json:"id_user"`
	CreatedAt       time.Time `json:"created_at"`
	// Relations
	Kategori *KategoriTransaksi `json:"kategori,omitempty" gorm:"foreignKey:IDKategori"`
	Bank     *Bank              `json:"bank,omitempty" gorm:"foreignKey:IDBank"`
	Customer *Customer          `json:"customer,omitempty" gorm:"foreignKey:IDCustomer"`
	Kavling  *Kavling           `json:"kavling,omitempty" gorm:"foreignKey:IDKavling"`
}

func (Pembayaran) TableName() string { return "pembayaran" }

// ─── Transaksi (Arus Kas) ─────────────────────────────────────────────────────

type Transaksi struct {
	ID            int       `json:"id" gorm:"primaryKey"`
	NoTransaksi   string    `json:"no_transaksi" gorm:"uniqueIndex;not null"`
	Jenis         string    `json:"jenis"` // Pemasukan | Pengeluaran
	Kategori      string    `json:"kategori"`
	Keterangan    string    `json:"keterangan"`
	Nominal       float64   `json:"nominal"`
	IDBank        *int      `json:"id_bank"`
	Tanggal       time.Time `json:"tanggal"`
	ReferensiID   *int      `json:"referensi_id"`
	ReferensiTipe string    `json:"referensi_tipe"`
	IDUser        *int      `json:"id_user"`
	CreatedAt     time.Time `json:"created_at"`
	// Relations
	Bank *Bank `json:"bank,omitempty" gorm:"foreignKey:IDBank"`
}

func (Transaksi) TableName() string { return "transaksi" }

// ─── Rekap Kredit (computed struct) ──────────────────────────────────────────

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
	NamaLokasi       string  `json:"nama_lokasi"`
}

// ─── Dashboard Ringkasan ──────────────────────────────────────────────────────

type DashboardRingkasan struct {
	TotalKavling     int     `json:"total_kavling"`
	Kosong           int     `json:"kosong"`
	Hold             int     `json:"hold"`
	BF               int     `json:"bf"`
	Akad             int     `json:"akad"`
	Lunas            int     `json:"lunas"`
	UserCancel       int     `json:"user_cancel"`
	JumlahCustomer   int     `json:"jumlah_customer"`
	JumlahProspek    int     `json:"jumlah_prospek"`
	JumlahTunggakan  int     `json:"jumlah_tunggakan"`
	NominalTunggakan float64 `json:"nominal_tunggakan"`
}

type ArusKasBulan struct {
	Bulan       int     `json:"bulan"`
	Pemasukan   float64 `json:"pemasukan"`
	Pengeluaran float64 `json:"pengeluaran"`
}

// ─── Pembayaran Detail Response ───────────────────────────────────────────────

type PembayaranDetailResponse struct {
	NamaCustomer    string     `json:"nama_customer"`
	LokasiPerumahan string     `json:"lokasi_perumahan"`
	KodeKavling     string     `json:"kode_kavling"`
	JenisPembayaran string     `json:"jenis_pembayaran"`
	TotalTagihan    float64    `json:"total_tagihan"`
	TotalBayar      float64    `json:"total_bayar"`
	SisaBayar       float64    `json:"sisa_bayar"`
	TagihanItems    []Tagihan  `json:"tagihan_items"`
	PemasukanItems  []Pembayaran `json:"pemasukan_items"`
}

// ─── Jatuh Tempo Row ──────────────────────────────────────────────────────────

type JatuhTempoRow struct {
	IDTransaksi    int     `json:"id_transaksi"`
	Nama           string  `json:"nama"`
	Lokasi         string  `json:"lokasi"`
	KodeKavling    string  `json:"kode_kavling"`
	HargaTanah     float64 `json:"harga_tanah"`
	Pembayaran     float64 `json:"pembayaran"`
	Pencairan      float64 `json:"pencairan"`
	Sisa           float64 `json:"sisa"`
	Keterlambatan  string  `json:"keterlambatan"`
	Cicilan        float64 `json:"cicilan"`
	Tenor          int     `json:"tenor"`
	JatuhTempo     string  `json:"jatuh_tempo"`
	JenisPembelian string  `json:"jenis_pembelian"`
}
