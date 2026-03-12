CREATE INDEX IF NOT EXISTS idx_users_first_name ON users (first_name);
CREATE INDEX IF NOT EXISTS idx_users_last_name ON users (last_name);
CREATE INDEX IF NOT EXISTS idx_friends_receiver_id ON friends (receiver_id);

