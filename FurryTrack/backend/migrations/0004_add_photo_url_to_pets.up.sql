SET client_encoding = 'UTF8';

ALTER TABLE pets 
ADD COLUMN photo_url VARCHAR(255);

COMMENT ON COLUMN users.role IS 'Photo URL';