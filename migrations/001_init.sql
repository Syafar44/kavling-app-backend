-- Migration: 001_init.sql
-- Buat semua tabel secara berurutan berdasarkan dependensi

-- 1. Tabel users
CREATE TABLE IF NOT EXISTS users (
    id          SERIAL PRIMARY KEY,
    username    VARCHAR(100) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,
    nama        VARCHAR(255) NOT NULL,
    is_admin    SMALLINT NOT NULL DEFAULT 0, -- 0=non-admin, 1=admin
    status      VARCHAR(20) NOT NULL DEFAULT 'AKTIF', -- AKTIF | BLOKIR
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Tabel menu
CREATE TABLE IF NOT EXISTS menu (
    id          SERIAL PRIMARY KEY,
    nama        VARCHAR(100) NOT NULL,
    icon        VARCHAR(100),
    url         VARCHAR(200),
    urutan      INT NOT NULL DEFAULT 0,
    parent_id   INT REFERENCES menu(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 3. Tabel hak_akses
CREATE TABLE IF NOT EXISTS hak_akses (
    id          SERIAL PRIMARY KEY,
    id_user     INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    id_menu     INT NOT NULL REFERENCES menu(id) ON DELETE CASCADE,
    status_hak  SMALLINT NOT NULL DEFAULT 1, -- 1=boleh, 0=tidak boleh
    UNIQUE(id_user, id_menu)
);

-- 4. Tabel throttle (pencatatan login gagal)
CREATE TABLE IF NOT EXISTS throttle (
    id          SERIAL PRIMARY KEY,
    ip_address  VARCHAR(45) NOT NULL,
    username    VARCHAR(100) NOT NULL,
    attempts    INT NOT NULL DEFAULT 1,
    last_attempt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(ip_address, username)
);

-- 5. Tabel aktifitas (log aktivitas user)
CREATE TABLE IF NOT EXISTS aktifitas (
    id          SERIAL PRIMARY KEY,
    id_user     INT REFERENCES users(id) ON DELETE SET NULL,
    aksi        VARCHAR(255) NOT NULL,
    keterangan  TEXT,
    ip_address  VARCHAR(45),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 6. Tabel bank (rekening bank/kas)
CREATE TABLE IF NOT EXISTS bank (
    id          SERIAL PRIMARY KEY,
    nama_bank   VARCHAR(100) NOT NULL,
    nama_rekening VARCHAR(100) NOT NULL,
    no_rekening VARCHAR(50) NOT NULL,
    saldo       NUMERIC(15,2) NOT NULL DEFAULT 0,
    is_kas      SMALLINT NOT NULL DEFAULT 0, -- 1=kas tunai
    status      SMALLINT NOT NULL DEFAULT 1, -- 1=aktif
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 7. Tabel konfigurasi (pengaturan umum)
CREATE TABLE IF NOT EXISTS konfigurasi (
    id              SERIAL PRIMARY KEY,
    nama_perusahaan VARCHAR(255) NOT NULL DEFAULT 'Kavling Mentaya',
    alamat          TEXT,
    telepon         VARCHAR(50),
    email           VARCHAR(100),
    website         VARCHAR(200),
    logo            VARCHAR(255),
    ttd_digital     VARCHAR(255), -- path file tanda tangan digital
    nama_ttd        VARCHAR(100),
    jabatan_ttd     VARCHAR(100),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 8. Tabel konfigurasi_wa
CREATE TABLE IF NOT EXISTS konfigurasi_wa (
    id          SERIAL PRIMARY KEY,
    id_device   VARCHAR(100) NOT NULL,
    no_telp     VARCHAR(20) NOT NULL,
    api_url     VARCHAR(255) NOT NULL,
    api_key     VARCHAR(255),
    is_aktif    SMALLINT NOT NULL DEFAULT 1,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 9. Tabel template (template pesan WhatsApp)
CREATE TABLE IF NOT EXISTS template (
    id          SERIAL PRIMARY KEY,
    nama        VARCHAR(100) NOT NULL,
    tipe        VARCHAR(50) NOT NULL, -- kwitansi | tagihan | dll
    isi         TEXT NOT NULL,
    is_aktif    SMALLINT NOT NULL DEFAULT 1,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 10. Tabel marketing (data sales/marketing)
CREATE TABLE IF NOT EXISTS marketing (
    id          SERIAL PRIMARY KEY,
    nama        VARCHAR(255) NOT NULL,
    no_telp     VARCHAR(20),
    alamat      TEXT,
    email       VARCHAR(100),
    persentase_komisi NUMERIC(5,2) NOT NULL DEFAULT 0,
    status      SMALLINT NOT NULL DEFAULT 1, -- 1=aktif, 0=nonaktif
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 11. Tabel customer (data pembeli)
CREATE TABLE IF NOT EXISTS customer (
    id          SERIAL PRIMARY KEY,
    nama        VARCHAR(255) NOT NULL,
    no_telp     VARCHAR(20),
    no_ktp      VARCHAR(30),
    alamat      TEXT,
    pekerjaan   VARCHAR(100),
    foto_ktp    VARCHAR(255),
    foto_kk     VARCHAR(255),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 12. Tabel kavling_peta (data kavling + koordinat SVG + status)
CREATE TABLE IF NOT EXISTS kavling_peta (
    id              SERIAL PRIMARY KEY,
    kode_kavling    VARCHAR(50) NOT NULL UNIQUE,
    luas            NUMERIC(10,2),
    harga_per_meter NUMERIC(15,2),
    harga_total     NUMERIC(15,2),
    status          SMALLINT NOT NULL DEFAULT 0,
    -- 0=Kosong, 1=Booking, 2=CASH/Lunas, 3=KREDIT
    shape_type      VARCHAR(20), -- polygon | path
    shape_data      TEXT,        -- atribut points atau d
    transform_data  VARCHAR(255),
    tgl_jatuh_tempo DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 13. Tabel transaksi_booking
CREATE TABLE IF NOT EXISTS transaksi_booking (
    id                  SERIAL PRIMARY KEY,
    no_booking          VARCHAR(50) NOT NULL UNIQUE,
    id_kavling          INT NOT NULL REFERENCES kavling_peta(id),
    id_customer         INT NOT NULL REFERENCES customer(id),
    id_marketing        INT REFERENCES marketing(id),
    nominal_booking     NUMERIC(15,2) NOT NULL DEFAULT 0,
    tgl_booking         DATE NOT NULL DEFAULT CURRENT_DATE,
    tgl_expired         DATE,
    keterangan          TEXT,
    status              SMALLINT NOT NULL DEFAULT 1, -- 1=aktif, 0=dibatalkan, 2=dikonversi
    id_user             INT REFERENCES users(id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 14. Tabel transaksi_kavling (data jual-beli kavling)
CREATE TABLE IF NOT EXISTS transaksi_kavling (
    id                  SERIAL PRIMARY KEY,
    no_transaksi        VARCHAR(50) NOT NULL UNIQUE,
    id_kavling          INT NOT NULL REFERENCES kavling_peta(id),
    id_customer         INT NOT NULL REFERENCES customer(id),
    id_marketing        INT REFERENCES marketing(id),
    id_booking          INT REFERENCES transaksi_booking(id),
    jenis_pembelian     SMALLINT NOT NULL,
    -- 2=CASH, 3=KREDIT
    harga_jual          NUMERIC(15,2) NOT NULL,
    uang_muka           NUMERIC(15,2) NOT NULL DEFAULT 0,
    lama_cicilan        INT NOT NULL DEFAULT 0, -- bulan
    cicilan_per_bulan   NUMERIC(15,2) NOT NULL DEFAULT 0,
    tgl_transaksi       DATE NOT NULL DEFAULT CURRENT_DATE,
    tgl_mulai_cicilan   DATE,
    keterangan          TEXT,
    id_user             INT REFERENCES users(id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 15. Tabel pembayaran (catatan pembayaran cicilan)
CREATE TABLE IF NOT EXISTS pembayaran (
    id                  SERIAL PRIMARY KEY,
    no_pembayaran       VARCHAR(50) NOT NULL UNIQUE,
    id_transaksi        INT NOT NULL REFERENCES transaksi_kavling(id),
    id_customer         INT NOT NULL REFERENCES customer(id),
    id_kavling          INT NOT NULL REFERENCES kavling_peta(id),
    id_bank             INT REFERENCES bank(id),
    tanggal             DATE NOT NULL DEFAULT CURRENT_DATE,
    pembayaran_ke       INT NOT NULL DEFAULT 1,
    jumlah_bayar        NUMERIC(15,2) NOT NULL,
    bukti_pembayaran    VARCHAR(255),
    keterangan          TEXT,
    id_user             INT REFERENCES users(id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 16. Tabel transaksi (arus kas pemasukan/pengeluaran)
CREATE TABLE IF NOT EXISTS transaksi (
    id          SERIAL PRIMARY KEY,
    no_transaksi VARCHAR(50) NOT NULL UNIQUE,
    jenis       VARCHAR(20) NOT NULL, -- Pemasukan | Pengeluaran
    kategori    VARCHAR(100) NOT NULL,
    keterangan  TEXT,
    nominal     NUMERIC(15,2) NOT NULL,
    id_bank     INT REFERENCES bank(id),
    tanggal     DATE NOT NULL DEFAULT CURRENT_DATE,
    referensi_id INT, -- id pembayaran atau lainnya
    referensi_tipe VARCHAR(50), -- pembayaran | manual
    id_user     INT REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index untuk performa query
CREATE INDEX IF NOT EXISTS idx_kavling_status ON kavling_peta(status);
CREATE INDEX IF NOT EXISTS idx_transaksi_kavling_id ON transaksi_kavling(id_kavling);
CREATE INDEX IF NOT EXISTS idx_pembayaran_kavling ON pembayaran(id_kavling);
CREATE INDEX IF NOT EXISTS idx_pembayaran_transaksi ON pembayaran(id_transaksi);
CREATE INDEX IF NOT EXISTS idx_transaksi_tanggal ON transaksi(tanggal);
CREATE INDEX IF NOT EXISTS idx_aktifitas_user ON aktifitas(id_user);
CREATE INDEX IF NOT EXISTS idx_throttle_ip_user ON throttle(ip_address, username);

-- Seed data default
INSERT INTO konfigurasi (nama_perusahaan, alamat, telepon)
VALUES ('Kavling Mentaya', 'Jl. Contoh No. 1, Sampit', '0812-0000-0000')
ON CONFLICT DO NOTHING;

-- User admin default (password: admin123 - harus diganti setelah pertama login)
-- Password di-hash menggunakan bcrypt, ini placeholder; hash dilakukan di aplikasi saat seed
INSERT INTO menu (nama, icon, url, urutan) VALUES
    ('Dashboard', 'dashboard', '/dashboard', 1),
    ('Master Data', 'folder', NULL, 2),
    ('Kavling', 'map', '/kavling', 3),
    ('Marketing', 'people', '/marketing', 4),
    ('Customer', 'person', '/customer', 5),
    ('Transaksi', 'receipt', NULL, 6),
    ('Booking', 'book', '/booking', 7),
    ('Transaksi Kavling', 'sell', '/transaksi', 8),
    ('Pembayaran', 'payment', '/pembayaran', 9),
    ('Keuangan', 'account_balance', NULL, 10),
    ('Arus Kas', 'swap_horiz', '/keuangan/arus-kas', 11),
    ('Rekap Kredit', 'credit_score', '/keuangan/rekap-kredit', 12),
    ('Laporan', 'bar_chart', '/laporan', 13),
    ('Pengaturan', 'settings', '/pengaturan', 14),
    ('Manajemen User', 'manage_accounts', '/users', 15)
ON CONFLICT DO NOTHING;
