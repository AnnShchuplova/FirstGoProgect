ALTER TABLE users 
ADD COLUMN deleted_at TIMESTAMP,
ADD COLUMN username TEXT NOT NULL DEFAULT '';

UPDATE users SET username = split_part(email, '@', 1);

ALTER TABLE pets
ADD COLUMN deleted_at TIMESTAMP,
ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE vaccines
ADD COLUMN deleted_at TIMESTAMP,
ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE vaccine_records
ADD COLUMN deleted_at TIMESTAMP,
ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE posts
ADD COLUMN deleted_at TIMESTAMP,
ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_pets_deleted_at ON pets(deleted_at);
CREATE INDEX idx_vaccines_deleted_at ON vaccines(deleted_at);
CREATE INDEX idx_vaccine_records_deleted_at ON vaccine_records(deleted_at);
CREATE INDEX idx_posts_deleted_at ON posts(deleted_at);