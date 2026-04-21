-- Migration 004: Add booking fee and cash payment fields to customer table
ALTER TABLE customer
    ADD COLUMN IF NOT EXISTS tanggal              DATE NOT NULL DEFAULT CURRENT_DATE,
    ADD COLUMN IF NOT EXISTS tanggal_batas_booking DATE,
    ADD COLUMN IF NOT EXISTS jumlah_pembayaran     NUMERIC(15,2) NOT NULL DEFAULT 0;

-- Convert jatuh_tempo from INT to DATE
ALTER TABLE customer
    ALTER COLUMN jatuh_tempo TYPE DATE USING NULL;

-- Reset status kavling yang id_customer-nya sudah NULL (customer dihapus).
UPDATE kavling SET status = 0 WHERE id_customer IS NULL AND status <> 0;
