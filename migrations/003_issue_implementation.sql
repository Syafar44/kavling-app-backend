-- Migration: 003_issue_implementation.sql
-- Implementasi penuh sesuai issue.md
-- Menambah tabel baru + extend tabel existing

-- ==============================================================
-- 1. EXTEND: users (tambah email, role, nama_lengkap)
-- ==============================================================
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS email        VARCHAR(100),
    ADD COLUMN IF NOT EXISTS role         VARCHAR(50) NOT NULL DEFAULT 'Admin',
    ADD COLUMN IF NOT EXISTS nama_lengkap VARCHAR(255);

UPDATE users SET nama_lengkap = nama WHERE nama_lengkap IS NULL;

-- ==============================================================
-- 2. EXTEND: marketing (tambah kode_marketing, jenis_kelamin, foto)
-- ==============================================================
ALTER TABLE marketing
    ADD COLUMN IF NOT EXISTS kode_marketing VARCHAR(20) UNIQUE,
    ADD COLUMN IF NOT EXISTS jenis_kelamin  VARCHAR(5),
    ADD COLUMN IF NOT EXISTS foto           VARCHAR(255);

-- Auto-generate kode_marketing untuk data existing
DO $$
DECLARE r RECORD;
BEGIN
    FOR r IN SELECT id FROM marketing WHERE kode_marketing IS NULL ORDER BY id LOOP
        UPDATE marketing SET kode_marketing = 'M-' || LPAD(r.id::TEXT, 3, '0') WHERE id = r.id;
    END LOOP;
END $$;

-- ==============================================================
-- 3. EXTEND: konfigurasi (tambah npwp, deskripsi)
-- ==============================================================
ALTER TABLE konfigurasi
    ADD COLUMN IF NOT EXISTS npwp      VARCHAR(30),
    ADD COLUMN IF NOT EXISTS deskripsi TEXT;

-- ==============================================================
-- 4. EXTEND: hak_akses (tambah permission columns)
-- ==============================================================
ALTER TABLE hak_akses
    ADD COLUMN IF NOT EXISTS lihat   BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS beranda BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS tambah  BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS edit    BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS hapus   BOOLEAN NOT NULL DEFAULT FALSE;

-- ==============================================================
-- 5. TABEL BARU: lokasi_kavling
-- ==============================================================
CREATE TABLE IF NOT EXISTS lokasi_kavling (
    id                    SERIAL PRIMARY KEY,
    nama                  VARCHAR(255) NOT NULL,
    nama_singkat          VARCHAR(20),
    header                VARCHAR(255),
    alamat                TEXT,
    nama_perusahaan       VARCHAR(255),
    nama_admin            VARCHAR(255),
    nama_mengetahui       VARCHAR(255),
    alamat_perusahaan     TEXT,
    telp_perusahaan       VARCHAR(50),
    kota_penandatangan    VARCHAR(100),
    nama_penandatangan    VARCHAR(255),
    jabatan_penandatangan VARCHAR(100),
    jenis_pembelian       VARCHAR(50) NOT NULL DEFAULT 'cash_kredit',
    urutan_lokasi         INT NOT NULL DEFAULT 1,
    jumlah_kavling        INT NOT NULL DEFAULT 0,
    kop_surat             VARCHAR(255),
    kwitansi              VARCHAR(255),
    foto_kavling          VARCHAR(255),
    svg_content           TEXT,
    svg_layout_json       JSONB,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==============================================================
-- 6. TABEL BARU: kavling (menggantikan kavling_peta dengan design baru)
-- ==============================================================
CREATE TABLE IF NOT EXISTS kavling (
    id              SERIAL PRIMARY KEY,
    id_lokasi       INT NOT NULL REFERENCES lokasi_kavling(id) ON DELETE CASCADE,
    kode_kavling    VARCHAR(50) NOT NULL,
    panjang_kanan   NUMERIC(10,2) NOT NULL DEFAULT 0,
    panjang_kiri    NUMERIC(10,2) NOT NULL DEFAULT 0,
    lebar_depan     NUMERIC(10,2) NOT NULL DEFAULT 0,
    lebar_belakang  NUMERIC(10,2) NOT NULL DEFAULT 0,
    luas_tanah      NUMERIC(10,2) NOT NULL DEFAULT 0,
    harga_per_meter NUMERIC(15,2) NOT NULL DEFAULT 0,
    harga_jual_cash NUMERIC(15,2) NOT NULL DEFAULT 0,
    no_sertipikat   VARCHAR(100),
    keterangan      TEXT,
    status          SMALLINT NOT NULL DEFAULT 0,
    -- 0=Ready, 1=HOLD, 2=BF, 3=AKAD, 4=UserCancel, 5=LUNAS
    id_customer     INT REFERENCES customer(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(id_lokasi, kode_kavling)
);

CREATE INDEX IF NOT EXISTS idx_kavling_lokasi ON kavling(id_lokasi);
CREATE INDEX IF NOT EXISTS idx_kavling_status ON kavling(status);
CREATE INDEX IF NOT EXISTS idx_kavling_customer ON kavling(id_customer);

-- ==============================================================
-- 7. TABEL BARU: customer_kavling (junction 1 customer → banyak kavling)
-- ==============================================================
CREATE TABLE IF NOT EXISTS customer_kavling (
    id          SERIAL PRIMARY KEY,
    id_customer INT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    id_kavling  INT NOT NULL REFERENCES kavling(id) ON DELETE CASCADE,
    UNIQUE(id_kavling)
);

-- ==============================================================
-- 8. EXTEND: customer (banyak field tambahan)
-- ==============================================================
ALTER TABLE customer
    ADD COLUMN IF NOT EXISTS kode_kontrak        VARCHAR(50) UNIQUE,
    ADD COLUMN IF NOT EXISTS no_ktp_pasangan     VARCHAR(30),
    ADD COLUMN IF NOT EXISTS tempat_lahir        VARCHAR(100),
    ADD COLUMN IF NOT EXISTS tanggal_lahir       DATE,
    ADD COLUMN IF NOT EXISTS jenis_kelamin       VARCHAR(5),
    ADD COLUMN IF NOT EXISTS email               VARCHAR(100),
    ADD COLUMN IF NOT EXISTS npwp                VARCHAR(30),
    ADD COLUMN IF NOT EXISTS alamat_domisili     TEXT,
    ADD COLUMN IF NOT EXISTS id_lokasi           INT REFERENCES lokasi_kavling(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS harga_jual          NUMERIC(15,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS id_marketing        INT REFERENCES marketing(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS status_penjualan    VARCHAR(20) NOT NULL DEFAULT 'BF',
    ADD COLUMN IF NOT EXISTS keterangan_cashback TEXT,
    ADD COLUMN IF NOT EXISTS id_bank_pembayaran  INT REFERENCES bank(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS jenis_pembelian     VARCHAR(30) NOT NULL DEFAULT 'KREDIT',
    ADD COLUMN IF NOT EXISTS cicilan_per_bulan   NUMERIC(15,2),
    ADD COLUMN IF NOT EXISTS tenor               INT,
    ADD COLUMN IF NOT EXISTS jatuh_tempo         INT,
    ADD COLUMN IF NOT EXISTS foto_pemohon        VARCHAR(255),
    ADD COLUMN IF NOT EXISTS foto_ktp_pasangan   VARCHAR(255),
    ADD COLUMN IF NOT EXISTS foto_npwp           VARCHAR(255),
    ADD COLUMN IF NOT EXISTS foto_bpjs           VARCHAR(255);

-- ==============================================================
-- 9. TABEL BARU: kategori_transaksi
-- ==============================================================
CREATE TABLE IF NOT EXISTS kategori_transaksi (
    id         SERIAL PRIMARY KEY,
    kode       VARCHAR(5) NOT NULL UNIQUE,
    kategori   VARCHAR(100) NOT NULL,
    jenis      VARCHAR(20) NOT NULL, -- PEMASUKAN | PENGELUARAN
    is_system  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==============================================================
-- 10. TABEL BARU: list_penjualan (master status progres)
-- ==============================================================
CREATE TABLE IF NOT EXISTS list_penjualan (
    id         SERIAL PRIMARY KEY,
    kode       VARCHAR(20),
    progres    VARCHAR(50) NOT NULL,
    keterangan TEXT,
    status     VARCHAR(20) NOT NULL DEFAULT 'Aktif',
    warna      VARCHAR(20) NOT NULL DEFAULT '#22c55e',
    urutan     INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==============================================================
-- 11. TABEL BARU: tagihan (invoice items per transaksi_kavling)
-- ==============================================================
CREATE TABLE IF NOT EXISTS tagihan (
    id           SERIAL PRIMARY KEY,
    id_transaksi INT NOT NULL REFERENCES transaksi_kavling(id) ON DELETE CASCADE,
    id_kategori  INT NOT NULL REFERENCES kategori_transaksi(id),
    deskripsi    TEXT,
    nominal      NUMERIC(15,2) NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tagihan_transaksi ON tagihan(id_transaksi);

-- ==============================================================
-- 12. EXTEND: pembayaran (tambah id_kategori, cara_bayar)
-- ==============================================================
ALTER TABLE pembayaran
    ADD COLUMN IF NOT EXISTS id_kategori INT REFERENCES kategori_transaksi(id),
    ADD COLUMN IF NOT EXISTS cara_bayar  VARCHAR(30);

-- ==============================================================
-- 13. TABEL BARU: hutang
-- ==============================================================
CREATE TABLE IF NOT EXISTS hutang (
    id                 SERIAL PRIMARY KEY,
    tanggal            DATE NOT NULL DEFAULT CURRENT_DATE,
    deskripsi          TEXT,
    nominal            NUMERIC(15,2) NOT NULL DEFAULT 0,
    status             VARCHAR(20) NOT NULL DEFAULT 'Belum Lunas',
    tanggal_pelunasan  DATE,
    id_bank            INT REFERENCES bank(id),
    lampiran           VARCHAR(255),
    sisa_bayar         NUMERIC(15,2) NOT NULL DEFAULT 0,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==============================================================
-- 14. TABEL BARU: piutang
-- ==============================================================
CREATE TABLE IF NOT EXISTS piutang (
    id                 SERIAL PRIMARY KEY,
    tanggal            DATE NOT NULL DEFAULT CURRENT_DATE,
    deskripsi          TEXT,
    nominal            NUMERIC(15,2) NOT NULL DEFAULT 0,
    status             VARCHAR(20) NOT NULL DEFAULT 'Belum Lunas',
    tanggal_pelunasan  DATE,
    id_bank            INT REFERENCES bank(id),
    lampiran           VARCHAR(255),
    sisa_bayar         NUMERIC(15,2) NOT NULL DEFAULT 0,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==============================================================
-- 15. TABEL BARU: mutasi_saldo
-- ==============================================================
CREATE TABLE IF NOT EXISTS mutasi_saldo (
    id             SERIAL PRIMARY KEY,
    tanggal        DATE NOT NULL DEFAULT CURRENT_DATE,
    id_bank_asal   INT NOT NULL REFERENCES bank(id),
    id_bank_tujuan INT NOT NULL REFERENCES bank(id),
    nominal        NUMERIC(15,2) NOT NULL DEFAULT 0,
    keterangan     TEXT,
    lampiran       VARCHAR(255),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==============================================================
-- 16. TABEL BARU: prospek
-- ==============================================================
CREATE TABLE IF NOT EXISTS prospek (
    id                SERIAL PRIMARY KEY,
    tanggal           DATE NOT NULL DEFAULT CURRENT_DATE,
    nama              VARCHAR(255) NOT NULL,
    no_telp           VARCHAR(20),
    usia              INT,
    pekerjaan         VARCHAR(100),
    penghasilan       NUMERIC(15,2),
    sumber_informasi  VARCHAR(100),
    rangking          VARCHAR(5) NOT NULL DEFAULT 'C',
    id_marketing      INT REFERENCES marketing(id),
    keterangan        TEXT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==============================================================
-- 17. TABEL BARU: notaris
-- ==============================================================
CREATE TABLE IF NOT EXISTS notaris (
    id         SERIAL PRIMARY KEY,
    nama       VARCHAR(255) NOT NULL,
    alamat     TEXT,
    no_telp    VARCHAR(20),
    keterangan VARCHAR(255) NOT NULL DEFAULT 'NOTARIS',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==============================================================
-- 18. TABEL BARU: legalitas
-- ==============================================================
CREATE TABLE IF NOT EXISTS legalitas (
    id              SERIAL PRIMARY KEY,
    id_kavling      INT NOT NULL REFERENCES kavling(id) ON DELETE CASCADE UNIQUE,
    atas_nama_surat VARCHAR(255),
    no_surat        VARCHAR(100),
    id_notaris      INT REFERENCES notaris(id) ON DELETE SET NULL,
    progres         VARCHAR(20) NOT NULL DEFAULT 'BF',
    bukti_foto      VARCHAR(255),
    keterangan      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==============================================================
-- 19. TABEL BARU: customer_file (upload dokumen customer)
-- ==============================================================
CREATE TABLE IF NOT EXISTS customer_file (
    id          SERIAL PRIMARY KEY,
    id_customer INT NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    tanggal     DATE NOT NULL DEFAULT CURRENT_DATE,
    nama_file   VARCHAR(100) NOT NULL,
    path_file   VARCHAR(255) NOT NULL,
    keterangan  TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_customer_file_cust ON customer_file(id_customer);

-- ==============================================================
-- 20. TABEL BARU: customer_arsip
-- ==============================================================
CREATE TABLE IF NOT EXISTS customer_arsip (
    id             SERIAL PRIMARY KEY,
    id_customer    INT REFERENCES customer(id) ON DELETE SET NULL,
    tanggal        DATE NOT NULL DEFAULT CURRENT_DATE,
    nama           VARCHAR(255) NOT NULL,
    id_marketing   INT REFERENCES marketing(id) ON DELETE SET NULL,
    lokasi         VARCHAR(255),
    status_progres VARCHAR(50) NOT NULL DEFAULT 'Aktif',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==============================================================
-- 21. TABEL BARU: media (pengaturan file media)
-- ==============================================================
CREATE TABLE IF NOT EXISTS media (
    id         SERIAL PRIMARY KEY,
    jenis_data VARCHAR(50) NOT NULL UNIQUE,
    keterangan TEXT,
    nama_file  VARCHAR(255),
    path_file  VARCHAR(255),
    is_logo    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==============================================================
-- 22. TABEL BARU: landing_konten
-- ==============================================================
CREATE TABLE IF NOT EXISTS landing_konten (
    id       SERIAL PRIMARY KEY,
    item     VARCHAR(50) NOT NULL,
    judul    VARCHAR(255) NOT NULL,
    artikel  TEXT,
    gambar   VARCHAR(255),
    icon     VARCHAR(255),
    urutan   INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==============================================================
-- SEED: kategori_transaksi system
-- ==============================================================
INSERT INTO kategori_transaksi (kode, kategori, jenis, is_system) VALUES
    ('001', 'Booking Fee',          'PEMASUKAN',    TRUE),
    ('002', 'Cicilan',              'PEMASUKAN',    TRUE),
    ('003', 'Pencairan Kredit',     'PEMASUKAN',    TRUE),
    ('004', 'Pembelian Cash',       'PEMASUKAN',    TRUE),
    ('005', 'Terima Saldo',         'PEMASUKAN',    TRUE),
    ('006', 'Pencairan Hutang',     'PEMASUKAN',    TRUE),
    ('007', 'Harga Tanah Kavling',  'PEMASUKAN',    TRUE),
    ('008', 'Mutasi Saldo',         'PENGELUARAN',  TRUE),
    ('009', 'Pembayaran Hutang',    'PENGELUARAN',  TRUE),
    ('010', 'Bayar Piutang',        'PENGELUARAN',  TRUE),
    ('011', 'Biaya Operasional',    'PENGELUARAN',  TRUE),
    ('012', 'Biaya Pembangunan',    'PENGELUARAN',  TRUE),
    ('013', 'Gaji',                 'PENGELUARAN',  TRUE),
    ('014', 'Komisi Marketing',     'PENGELUARAN',  TRUE),
    ('015', 'Biaya Notaris',        'PENGELUARAN',  TRUE),
    ('016', 'Biaya Admin',          'PENGELUARAN',  TRUE),
    ('017', 'Biaya Marketing',      'PENGELUARAN',  TRUE),
    ('018', 'Pembayaran OP Jalan',  'PENGELUARAN',  TRUE),
    ('019', 'Pembelian Barang',     'PENGELUARAN',  FALSE)
ON CONFLICT (kode) DO NOTHING;

-- ==============================================================
-- SEED: list_penjualan (6 status progres)
-- ==============================================================
INSERT INTO list_penjualan (kode, progres, keterangan, status, warna, urutan) VALUES
    ('Rdy',   'Ready',       'Kavling Siap Dijual',      'Aktif',       '#22c55e', 1),
    ('HLD',   'HOLD',        'Kavling Sedang Ditahan',   'Aktif',       '#eab308', 2),
    ('BF',    'BF',          'Booking Fee',              'Aktif',       '#3b82f6', 3),
    ('AKAD',  'AKAD',        'Proses Akad',              'Aktif',       '#8b5cf6', 4),
    ('',      'User Cancel', 'Dibatalkan oleh Customer', 'Tidak Aktif', '#ef4444', 5),
    ('LUNAS', 'LUNAS',       'Kavling Sudah Lunas',      'Aktif',       '#06b6d4', 6)
ON CONFLICT DO NOTHING;

-- ==============================================================
-- SEED: media default (3 rows)
-- ==============================================================
INSERT INTO media (jenis_data, keterangan, is_logo) VALUES
    ('logo website', 'Logo yang ditampilkan pada halaman login', TRUE),
    ('fav icon',     'Logo yang ditampilkan pada judul tab aplikasi', FALSE),
    ('kop surat',    'Background Cetak Rekap pada menu pembayaran', FALSE)
ON CONFLICT (jenis_data) DO NOTHING;

-- ==============================================================
-- SEED: menu sesuai hak_akses MENU_ROWS (34 row)
-- ==============================================================
TRUNCATE TABLE menu RESTART IDENTITY CASCADE;
INSERT INTO menu (id, nama, icon, url, urutan, parent_id) VALUES
    (1,  'Beranda',             'home',         '/beranda',                         1,  NULL),
    (2,  'Dashboard',           'bar-chart',    '/dashboard',                       2,  NULL),
    (3,  'Siteplan',            'map',          '/siteplan',                        3,  NULL),
    (4,  'Pembayaran',          'credit-card',  '/pembayaran',                      4,  NULL),
    (5,  'Keuangan',            'dollar-sign',  '#',                                5,  NULL),
    (6,  'Pemasukan',           'trending-up',  '/keuangan/pemasukan',              6,  5),
    (7,  'Pengeluaran',         'trending-down','/keuangan/pengeluaran',            7,  5),
    (8,  'Hutang',              'minus-circle', '/keuangan/hutang',                 8,  5),
    (9,  'Piutang',             'plus-circle',  '/keuangan/piutang',                9,  5),
    (10, 'Kategori Transaksi',  'tag',          '/keuangan/kategori',               10, 5),
    (11, 'Mutasi Saldo',        'refresh-cw',   '/keuangan/mutasi-saldo',           11, 5),
    (12, 'Laporan Arus Kas',    'file-text',    '/keuangan/laporan-arus-kas',       12, 5),
    (13, 'Customer',            'users',        '#',                                13, NULL),
    (14, 'Customer',            'user',         '/customer',                        14, 13),
    (15, 'Prospek',             'user-plus',    '/customer/prospek',                15, 13),
    (16, 'Upload File',         'upload',       '/customer/upload-file',            16, 13),
    (17, 'Arsip Customer',      'archive',      '/customer/arsip',                  17, 13),
    (18, 'Legalitas',           'file-badge',   '/legalitas',                       18, NULL),
    (19, 'Master Data',         'database',     '#',                                19, NULL),
    (20, 'Marketing',           'megaphone',    '/master/marketing',                20, 19),
    (21, 'Lokasi Kavling',      'map-pin',      '/master/lokasi-kavling',           21, 19),
    (22, 'Kavling',             'grid',         '/master/kavling',                  22, 19),
    (23, 'Notaris',             'shield',       '/master/notaris',                  23, 19),
    (24, 'Pengaturan',          'settings',     '#',                                24, NULL),
    (25, 'Pengaturan Profile',  'building',     '/pengaturan/profile',              25, 24),
    (26, 'Pengaturan Media',    'image',        '/pengaturan/media',                26, 24),
    (27, 'Hak Akses',           'lock',         '/pengaturan/hak-akses',            27, 24),
    (28, 'Pengguna',            'user-cog',     '/pengaturan/pengguna',             28, 24),
    (29, 'List Penjualan',      'list',         '/pengaturan/list-penjualan',       29, 24),
    (30, 'Pengaturan Landing',  'globe',        '/pengaturan/landing',              30, 24),
    (31, 'Bank',                'credit-card',  '/pengaturan/bank',                 31, 24),
    (32, 'Pengaturan WA',       'message-circle','#',                              32, NULL),
    (33, 'Template Pesan',      'message-square','/pengaturan/wa/template',         33, 32),
    (34, 'Konfigurasi Koneksi', 'wifi',         '/pengaturan/wa/koneksi',           34, 32)
ON CONFLICT DO NOTHING;

-- ==============================================================
-- SEED: lokasi_kavling (3 lokasi dummy)
-- ==============================================================
INSERT INTO lokasi_kavling
    (id, nama, nama_singkat, header, alamat, nama_perusahaan, nama_admin,
     nama_mengetahui, alamat_perusahaan, telp_perusahaan, kota_penandatangan,
     nama_penandatangan, jabatan_penandatangan, jenis_pembelian, urutan_lokasi, jumlah_kavling)
VALUES
    (1, 'Punten Regency', 'PR',  'KAVLING PUNTEN REGENCY',     'Jl. Mulia Kota Batu',
     'PT. BERKAH KAVLING NUSANTARA', 'Budi Santoso', 'Direktur Utama',
     'Jl. Mulia No.1 Kota Batu', '(0341) 555-0001', 'Kota Batu',
     'Ir. Hendra Wijaya', 'Direktur', 'cash_kredit', 1, 206),
    (2, 'SPG',            'SPG', 'KAVLING SPG',                'Jl. Diponegoro Kota Batu',
     'PT. MULIA ASRI SENTOSA', 'Rina Mulyani', 'Manager Operasional',
     'Jl. Diponegoro No.45 Kota Batu', '(0341) 555-0002', 'Kota Batu',
     'Drs. Agung Prasetyo', 'Direktur Utama', 'kredit', 2, 224),
    (3, 'Ngaglik Residen','NR',  'KAVLING NGAGLIK RESIDEN',    'Jl. Warureja - Pemalang',
     'PT. MULIA ASRI SENTOSA', 'Suharto', 'Manager Area',
     'Jl. Pemalang No.10', '(0284) 555-0003', 'Pemalang',
     'Teguh Santoso', 'Kepala Cabang', 'cash_kredit', 3, 227)
ON CONFLICT DO NOTHING;

-- Update sequence
SELECT setval('lokasi_kavling_id_seq', 3);

-- ==============================================================
-- SEED: notaris (3 data dummy)
-- ==============================================================
INSERT INTO notaris (nama, alamat, no_telp, keterangan) VALUES
    ('Notaris Hendra Setiawan', 'Jl. Diponegoro No. 12, Malang',     '0812-3456-7890', 'NOTARIS'),
    ('Notaris Sari Dewi, SH',   'Jl. Ahmad Yani No. 45, Kota Batu',  '0813-2345-6789', 'NOTARIS'),
    ('Notaris Ahmad Fauzi',     'Jl. Perintis Kemerdekaan No. 8',    '0811-9876-5432', 'NOTARIS')
ON CONFLICT DO NOTHING;

-- Index untuk performa
CREATE INDEX IF NOT EXISTS idx_kavling_kode ON kavling(kode_kavling);
CREATE INDEX IF NOT EXISTS idx_customer_kode ON customer(kode_kontrak);
CREATE INDEX IF NOT EXISTS idx_legalitas_kavling ON legalitas(id_kavling);
CREATE INDEX IF NOT EXISTS idx_tagihan_transaksi ON tagihan(id_transaksi);
CREATE INDEX IF NOT EXISTS idx_prospek_marketing ON prospek(id_marketing);
