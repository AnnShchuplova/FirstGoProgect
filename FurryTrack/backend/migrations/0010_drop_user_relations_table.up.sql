ALTER TABLE IF EXISTS some_other_table 
DROP CONSTRAINT IF EXISTS fk_user_relation;

DROP TABLE IF EXISTS user_relations CASCADE;