CREATE TABLE IF NOT EXISTS withdrawals (
                                           id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
                                           user_id uuid REFERENCES users(id) ON DELETE CASCADE,
                                           order_num VARCHAR (50) NOT NULL UNIQUE,
                                           status VARCHAR(50) DEFAULT 'NEW',
                                           processed_at TIMESTAMP,
                                           sum FLOAT DEFAULT 0
);
