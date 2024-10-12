-- Deploy vetchi:1-initdb to pg

BEGIN;

CREATE TABLE IF NOT EXISTS hub_users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

--- Must match libvetchi.OnboardStatus
CREATE TYPE onboard_status AS ENUM ('DOMAIN_NOT_VERIFIED', 'DOMAIN_VERIFIED_ONBOARDING_PENDING', 'DOMAIN_ONBOARDED');

CREATE TABLE IF NOT EXISTS employers (
    client_id TEXT PRIMARY KEY,
    onboard_status onboard_status NOT NULL,
    onboarding_admin TEXT NOT NULL,
    onboarding_email_sent_at TIMESTAMP WITH TIME ZONE,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

COMMIT;
