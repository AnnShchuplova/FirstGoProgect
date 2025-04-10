DROP INDEX IF EXISTS idx_users_deleted_at;
ALTER TABLE users DROP COLUMN IF EXISTS deleted_at, DROP COLUMN IF EXISTS username;

DROP INDEX IF EXISTS idx_pets_deleted_at;
ALTER TABLE pets DROP COLUMN IF EXISTS deleted_at, DROP COLUMN IF EXISTS updated_at;

DROP INDEX IF EXISTS idx_vaccines_deleted_at;
ALTER TABLE vaccines DROP COLUMN IF EXISTS deleted_at, DROP COLUMN IF EXISTS updated_at;

DROP INDEX IF EXISTS idx_vaccine_records_deleted_at;
ALTER TABLE vaccine_records DROP COLUMN IF EXISTS deleted_at, DROP COLUMN IF EXISTS updated_at;

DROP INDEX IF EXISTS idx_posts_deleted_at;
ALTER TABLE posts DROP COLUMN IF EXISTS deleted_at, DROP COLUMN IF EXISTS updated_at;