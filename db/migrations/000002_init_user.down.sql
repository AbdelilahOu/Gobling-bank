ALTER TABLE IF EXISTS accounts DROP CONSTRAINT IF EXISTS "owner_currency_key";

ALTER TABLE IF EXISTS accounts DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";

drop table users;

