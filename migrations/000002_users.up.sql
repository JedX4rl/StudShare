CREATE TABLE IF NOT EXISTS users
(
    id            UUID PRIMARY KEY,
    email         TEXT UNIQUE NOT NULL,
    password_hash TEXT        NOT NULL,
    name          TEXT        NOT NULL,
    surname       TEXT        NOT NULL,
    PHONE         TEXT        NOT NULL,
    is_admin      BOOLEAN     NOT NULL DEFAULT FALSE, -- может быть 'user' или 'admin'
    rating DECIMAL(10, 2) NOT NULL DEFAULT 5.0,
    created_at    TIMESTAMP   NOT NULL DEFAULT NOW()
);