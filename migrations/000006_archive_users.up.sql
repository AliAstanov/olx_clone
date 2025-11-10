
CREATE TABLE IF NOT EXISTS users_archive (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP,
    deleted_at TIMESTAMP
);