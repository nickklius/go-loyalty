CREATE TABLE IF NOT EXISTS gophermart.users (
                                                id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
                                                first_name VARCHAR(100),
                                                last_name VARCHAR(100),
                                                login VARCHAR(50) NOT NULL UNIQUE,
                                                password text NOT NULL,
                                                created_at timestamp
)