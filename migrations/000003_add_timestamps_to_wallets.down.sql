ALTER TABLE paym_processor.wallets 
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS updated_at;