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

INSERT INTO employers (client_id, onboard_status) VALUES ('domain-verified-email-not-sent.example', 'DOMAIN_VERIFIED_EMAIL_NOT_SENT');
INSERT INTO employers (client_id, onboard_status) VALUES ('domain-verified-email-sent.example', 'DOMAIN_VERIFIED_EMAIL_SENT');
INSERT INTO employers (client_id, onboard_status) VALUES ('domain-onboarded.example', 'DOMAIN_ONBOARDED');

COMMIT;
