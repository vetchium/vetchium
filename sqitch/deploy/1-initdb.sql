-- Deploy vetchi:1-initdb to pg

BEGIN;

CREATE TABLE IF NOT EXISTS hub_users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS employers (
    client_id TEXT PRIMARY KEY,
    onboard_status TEXT NOT NULL
);

COMMIT;
