-- Migration: 002_denah_kavling.sql
-- Mengubah sistem kavling dari flat list menjadi hierarki Denah → Kavling

-- 1. Tabel denah_kavling (parent)
CREATE TABLE IF NOT EXISTS denah_kavling (
    id          SERIAL PRIMARY KEY,
    nama        VARCHAR(255) NOT NULL,
    svg_content TEXT NOT NULL,        -- kode SVG lengkap yang diinput user
    viewbox     VARCHAR(100),         -- diambil dari atribut viewBox SVG
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Alter kavling_peta: tambah kolom baru, hapus kolom lama
ALTER TABLE kavling_peta
    ADD COLUMN IF NOT EXISTS denah_kavling_id INT REFERENCES denah_kavling(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS kode_map         TEXT,         -- atribut d dari <path>
    ADD COLUMN IF NOT EXISTS panjang_kanan    NUMERIC(10,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS panjang_kiri     NUMERIC(10,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS lebar_depan      NUMERIC(10,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS lebar_belakang   NUMERIC(10,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS luas_tanah       NUMERIC(10,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS harga_jual_cash  NUMERIC(15,2) NOT NULL DEFAULT 0;

-- Hapus kolom lama yang tidak dipakai lagi
-- (kolom luas digantikan luas_tanah, shape_type/shape_data/transform_data digantikan kode_map)
ALTER TABLE kavling_peta
    DROP COLUMN IF EXISTS luas,
    DROP COLUMN IF EXISTS shape_type,
    DROP COLUMN IF EXISTS shape_data,
    DROP COLUMN IF EXISTS transform_data;

-- 3. Index tambahan
CREATE INDEX IF NOT EXISTS idx_kavling_denah ON kavling_peta(denah_kavling_id);
