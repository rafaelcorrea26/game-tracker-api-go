DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'games' AND column_name = 'platform'
    ) THEN
        ALTER TABLE games ADD COLUMN platform TEXT DEFAULT '';
    END IF;
END $$;
