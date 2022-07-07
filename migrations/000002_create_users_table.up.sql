CREATE TABLE IF NOT EXISTS gophermart.users (
                                                id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
                                                login VARCHAR(50) NOT NULL UNIQUE,
                                                password text NOT NULL,
                                                created_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)