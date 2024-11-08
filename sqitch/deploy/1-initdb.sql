-- Deploy vetchi:1-initdb to pg

ALTER SYSTEM SET log_statement = 'all';
ALTER SYSTEM SET log_min_duration_statement = 0;
ALTER SYSTEM SET log_duration = 'on';
SELECT pg_reload_conf();

BEGIN;

CREATE TABLE IF NOT EXISTS hub_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name TEXT NOT NULL,
    handle TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TYPE email_states AS ENUM ('PENDING', 'PROCESSED');
CREATE TABLE emails(
	email_key UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	email_from TEXT NOT NULL,
	email_to TEXT ARRAY NOT NULL,
	email_cc TEXT ARRAY,
	email_bcc TEXT ARRAY,
	email_subject TEXT NOT NULL,
	email_html_body TEXT NOT NULL,
	email_text_body TEXT NOT NULL,
	email_state email_states NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),
	processed_at TIMESTAMP WITH TIME ZONE
);

---

CREATE TYPE client_id_types AS ENUM ('DOMAIN');
CREATE TYPE employer_states AS ENUM (
    'ONBOARD_PENDING',
    'ONBOARDED',
    'DEBOARDED'
);
CREATE TABLE employers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id_type client_id_types NOT NULL,
    employer_state employer_states NOT NULL,
    onboard_admin_email TEXT NOT NULL,

    -- TODO: Perhaps we can move this to org_user_tokens ?
    onboard_secret_token TEXT,
    token_valid_till TIMESTAMP WITH TIME ZONE,

    --- Despite its name, it should not be confused with an email address.
    --- This is the rowid in the 'emails' table for the welcome email sent.
    onboard_email_id UUID REFERENCES emails(email_key),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

---

CREATE TYPE domain_states AS ENUM (
    'VERIFIED',
    'DEBOARDED'
);
CREATE TABLE domains (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    domain_name TEXT NOT NULL,
    CONSTRAINT uniq_domain_name UNIQUE (domain_name),

    domain_state domain_states NOT NULL,

    employer_id UUID REFERENCES employers(id) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

---

CREATE TYPE org_user_roles AS ENUM (
    'ADMIN',
    'COST_CENTERS_CRUD',
    'COST_CENTERS_VIEWER',
    'LOCATIONS_CRUD',
    'LOCATIONS_VIEWER',
    'ORG_USERS_CRUD',
    'ORG_USERS_VIEWER',
    'OPENINGS_CRUD',
    'OPENINGS_VIEWER'
);
CREATE TYPE org_user_states AS ENUM (
    'ACTIVE_ORG_USER',
    'INVITED_ORG_USER',
    'ADDED_ORG_USER',
    'DISABLED_ORG_USER',
    'REPLICATED_ORG_USER'
);
CREATE TABLE org_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL,
    name TEXT NOT NULL,
    password_hash TEXT,
    org_user_roles org_user_roles[] NOT NULL,
    org_user_state org_user_states NOT NULL,

--- As of now, we have only one org per employer. This may change in future.
    employer_id UUID REFERENCES employers(id) NOT NULL,
    CONSTRAINT uniq_email_employer_id UNIQUE (email, employer_id),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TYPE token_types AS ENUM (
    'EMPLOYER_SESSION',
    'EMPLOYER_LTS',
    'EMPLOYER_TFA_TOKEN',

    -- Perhaps this alone can be moved to a new table, to minimize collisions !?
    'EMPLOYER_TFA_CODE',

    'EMPLOYER_INVITE'
);
CREATE TABLE org_user_tokens (
    token TEXT,
    token_valid_till TIMESTAMP WITH TIME ZONE NOT NULL,
    token_type token_types NOT NULL,

    CONSTRAINT org_user_tokens_pkey PRIMARY KEY (token),
    org_user_id UUID REFERENCES org_users(id) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

---

CREATE TYPE cost_center_states AS ENUM ('ACTIVE_CC', 'DEFUNCT_CC');
CREATE TABLE org_cost_centers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cost_center_name TEXT NOT NULL,
    cost_center_state cost_center_states NOT NULL,
    notes TEXT NOT NULL,

    employer_id UUID REFERENCES employers(id) NOT NULL,
    CONSTRAINT uniq_cost_center_name_employer_id UNIQUE (cost_center_name, employer_id),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

---

CREATE TYPE location_states AS ENUM ('ACTIVE_LOCATION', 'DEFUNCT_LOCATION');
CREATE TABLE locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    country_code TEXT NOT NULL,
    postal_address TEXT NOT NULL,
    postal_code TEXT NOT NULL,
    openstreetmap_url TEXT,
    city_aka TEXT ARRAY,

    location_state location_states NOT NULL,

    employer_id UUID REFERENCES employers(id) NOT NULL,
    CONSTRAINT uniq_location_title_employer_id UNIQUE (title, employer_id),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

---

CREATE TYPE opening_states AS ENUM ('DRAFT_OPENING', 'ACTIVE_OPENING', 'SUSPENDED_OPENING', 'CLOSED_OPENING');
CREATE TYPE opening_types AS ENUM ('FULL_TIME', 'PART_TIME', 'CONTRACT', 'INTERNSHIP', 'UNSPECIFIED');
CREATE TYPE education_levels AS ENUM ('BACHELOR', 'MASTER', 'DOCTORATE', 'NOT_MATTERS', 'UNSPECIFIED');
CREATE TABLE openings (
    employer_id UUID REFERENCES employers(id) NOT NULL,
    id TEXT NOT NULL,
    CONSTRAINT openings_pk PRIMARY KEY (employer_id, id),

    title TEXT NOT NULL,
    positions INTEGER NOT NULL,
    jd TEXT NOT NULL,
    hiring_manager UUID REFERENCES org_users(id) NOT NULL,
    cost_center_id UUID REFERENCES org_cost_centers(id),
    employer_notes TEXT NOT NULL,
    remote_country_codes TEXT[] NOT NULL,
    remote_timezones TEXT[] NOT NULL,
    opening_type opening_types NOT NULL,
    yoe_min INTEGER NOT NULL,
    yoe_max INTEGER NOT NULL,
    min_education_level education_levels NOT NULL,
    salary_min NUMERIC NOT NULL,
    salary_max NUMERIC NOT NULL,
    salary_currency TEXT NOT NULL,
    current_state opening_states NOT NULL,
    approval_waiting_state opening_states,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),
    last_updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TABLE opening_recruiters (
    employer_id UUID NOT NULL,
    opening_id TEXT NOT NULL,
    CONSTRAINT fk_opening FOREIGN KEY (employer_id, opening_id) REFERENCES openings (employer_id, id),

    recruiter_id UUID REFERENCES org_users(id) NOT NULL,

    PRIMARY KEY (employer_id, opening_id, recruiter_id)
);

CREATE TABLE opening_hiring_team(
    employer_id UUID NOT NULL,
    opening_id TEXT NOT NULL,
    CONSTRAINT fk_opening FOREIGN KEY (employer_id, opening_id) REFERENCES openings (employer_id, id),

    hiring_team_member_id UUID REFERENCES hub_users(id) NOT NULL,

    PRIMARY KEY (employer_id, opening_id, hiring_team_member_id)
);

CREATE TABLE opening_locations(
    employer_id UUID NOT NULL,
    opening_id TEXT NOT NULL,
    CONSTRAINT fk_opening FOREIGN KEY (employer_id, opening_id) REFERENCES openings (employer_id, id),

    location_id UUID REFERENCES locations(id) NOT NULL,

    PRIMARY KEY (employer_id, opening_id, location_id)
);

COMMIT;
