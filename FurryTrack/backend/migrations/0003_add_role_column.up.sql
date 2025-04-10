ALTER TABLE users 
ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'USER';

COMMENT ON COLUMN users.role IS 'USER, ADMIN, VET, etc.';

UPDATE users SET role = 'ADMIN' WHERE is_admin = true;