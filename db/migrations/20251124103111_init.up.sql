CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    phone_number TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    birthdate DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);