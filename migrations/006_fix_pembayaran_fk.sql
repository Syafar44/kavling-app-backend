-- Same root cause as 005: pembayaran.id_kavling FK originally pointed to legacy kavling_peta(id).
-- Repoint to kavling(id).

DO $$
DECLARE
    fk_name text;
BEGIN
    SELECT conname INTO fk_name
    FROM pg_constraint
    WHERE conrelid = 'pembayaran'::regclass
      AND contype = 'f'
      AND conkey = ARRAY[(
          SELECT attnum FROM pg_attribute
          WHERE attrelid = 'pembayaran'::regclass AND attname = 'id_kavling'
      )::smallint];

    IF fk_name IS NOT NULL THEN
        EXECUTE format('ALTER TABLE pembayaran DROP CONSTRAINT %I', fk_name);
    END IF;
END $$;

ALTER TABLE pembayaran
    ADD CONSTRAINT pembayaran_id_kavling_fkey
    FOREIGN KEY (id_kavling) REFERENCES kavling(id) ON DELETE RESTRICT;
