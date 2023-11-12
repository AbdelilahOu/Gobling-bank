DROP TABLE IF EXISTS entries;
DROP TABLE IF EXISTS transfers;
DROP TABLE IF EXISTS acounts;

ALTER TABLE entries DROP CONSTRAINT IF EXISTS entries_account_id_fkey;
ALTER TABLE transfers DROP CONSTRAINT IF EXISTS transfers_from_account_id_fkey;
ALTER TABLE transfers DROP CONSTRAINT IF EXISTS transfers_to_account_id_fkey;

DROP INDEX IF EXISTS accounts_owner_idx;
DROP INDEX IF EXISTS entries_account_id_idx;
DROP INDEX IF EXISTS transfers_from_account_id_idx;
DROP INDEX IF EXISTS transfers_to_account_id_idx;
DROP INDEX IF EXISTS transfers_from_to_account_id_idx;

COMMENT ON COLUMN entries.amount IS NULL;
COMMENT ON COLUMN transfers.amount IS NULL;