ALTER TABLE users
    ADD COLUMN IF NOT EXISTS balance FLOAT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS spent FLOAT DEFAULT 0;
