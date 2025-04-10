CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE posts 
ALTER COLUMN id TYPE UUID USING uuid_generate_v4(),
ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE user_relations 
ALTER COLUMN id TYPE UUID USING uuid_generate_v4(),
ALTER COLUMN id SET DEFAULT uuid_generate_v4();


ALTER TABLE comments
ALTER COLUMN post_id TYPE UUID USING uuid_generate_v4();

ALTER TABLE post_likes
ALTER COLUMN post_id TYPE UUID USING uuid_generate_v4();