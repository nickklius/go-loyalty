CREATE TABLE IF NOT EXISTS orders (
                                      id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
                                      user_id uuid REFERENCES users(id) ON DELETE CASCADE,
                                      number VARCHAR(50) NOT NULL UNIQUE,
                                      status VARCHAR(50),
                                      uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                      accrual FLOAT DEFAULT 0
);