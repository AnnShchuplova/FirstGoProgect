CREATE TABLE user_relations (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    follower_id UUID NOT NULL,
    following_id UUID NOT NULL,
    FOREIGN KEY (follower_id) REFERENCES users(id),
    FOREIGN KEY (following_id) REFERENCES users(id),
    UNIQUE(follower_id, following_id)
);

CREATE INDEX idx_user_relations_follower ON user_relations(follower_id);
CREATE INDEX idx_user_relations_following ON user_relations(following_id);