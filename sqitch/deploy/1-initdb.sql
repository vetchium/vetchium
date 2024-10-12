-- Deploy vetchi:1-initdb to pg

BEGIN;

CREATE TABLE IF NOT EXISTS hub_users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

--- Must match dao.EmailState
CREATE TYPE email_states AS ENUM ('PENDING', 'PROCESSED');

CREATE TABLE emails(
	id BIGSERIAL PRIMARY KEY,
	email_from TEXT NOT NULL,
	email_to TEXT ARRAY NOT NULL,
	email_cc TEXT ARRAY,
	email_bcc TEXT ARRAY,
	email_subject TEXT NOT NULL,
	email_html_body TEXT NOT NULL,
	email_text_body TEXT NOT NULL,
	email_state email_states NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('UTC', now()),
	processed_at TIMESTAMP WITH TIME ZONE
);

--- Must match libvetchi.OnboardStatus
CREATE TYPE onboard_status AS ENUM ('DOMAIN_NOT_VERIFIED', 'DOMAIN_VERIFIED_ONBOARD_PENDING', 'DOMAIN_ONBOARDED');

CREATE TABLE IF NOT EXISTS employers (
    client_id TEXT PRIMARY KEY,
    onboard_status onboard_status NOT NULL,
    onboard_admin TEXT,
    onboard_secret_token TEXT,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),

    onboard_email_id BIGINT REFERENCES emails(id)
);

COMMIT;
