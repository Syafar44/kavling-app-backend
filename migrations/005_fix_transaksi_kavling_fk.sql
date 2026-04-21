-- Fix FK: transaksi_kavling.id_kavling originally referenced legacy kavling_peta(id).
-- App now stores kavlings in the v2 `kavling` table, so INSERTs fail silently on FK violation.
-- Repoint the FK to kavling(id).

DO $$
DECLARE
    fk_name text;
BEGIN
    SELECT conname INTO fk_name
    FROM pg_constraint
    WHERE conrelid = 'transaksi_kavling'::regclass
      AND contype = 'f'
      AND conkey = ARRAY[(
          SELECT attnum FROM pg_attribute
          WHERE attrelid = 'transaksi_kavling'::regclass AND attname = 'id_kavling'
      )::smallint];

    IF fk_name IS NOT NULL THEN
        EXECUTE format('ALTER TABLE transaksi_kavling DROP CONSTRAINT %I', fk_name);
    END IF;
END $$;

ALTER TABLE transaksi_kavling
    ADD CONSTRAINT transaksi_kavling_id_kavling_fkey
    FOREIGN KEY (id_kavling) REFERENCES kavling(id) ON DELETE RESTRICT;
