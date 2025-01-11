BEGIN;
--- email table primary key uuids should end in 2 digits, 11, 12, 13, etc
--- employer table primary key uuids should end in 3 digits, 201, 202, 203, etc
--- domain table primary key uuids should end in 4 digits, 3001, 3002, 3003, etc
--- org_users table primary key uuids should end in 5 digits, 40001, 40002, 40003, etc
--- cost_centers table primary key uuids should end in 6 digits, 50001, 50002, 50003, etc
--- locations table primary key uuids should end in 7 digits, 60001, 60002, 60003, etc

-- Create test hub users
INSERT INTO hub_users (
    id,
    full_name,
    handle,
    email,
    password_hash,
    state,
    resident_country_code,
    resident_city,
    preferred_language,
    created_at,
    updated_at
) VALUES 
    ('12345678-0009-0009-0009-000000050001'::uuid, 'Active Hub User', 'active_hub_user', 'active@applyopening.example', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'IND', 'Bangalore', 'en', timezone('UTC'::text, now()), timezone('UTC'::text, now())),
    ('12345678-0009-0009-0009-000000050002'::uuid, 'Disabled Hub User', 'disabled_hub_user', 'disabled@applyopening.example', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'DISABLED_HUB_USER', 'IND', 'Chennai', 'en', timezone('UTC'::text, now()), timezone('UTC'::text, now())),
    ('12345678-0009-0009-0009-000000050003'::uuid, 'Deleted Hub User', 'deleted_hub_user', 'deleted@applyopening.example', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'DELETED_HUB_USER', 'IND', 'Mumbai', 'en', timezone('UTC'::text, now()), timezone('UTC'::text, now()));

-- Create test employer
INSERT INTO emails (email_key, email_from, email_to, email_cc, email_bcc, email_subject, email_html_body, email_text_body, email_state, created_at, processed_at)
VALUES ('12345678-0009-0009-0009-000000000011'::uuid, 'no-reply@vetchi.org', ARRAY['admin@applyopening.example'], NULL, NULL, 'Welcome to Vetchi', 'Welcome HTML', 'Welcome Text', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now()));

INSERT INTO employers (id, client_id_type, employer_state, company_name, onboard_admin_email, onboard_secret_token, token_valid_till, onboard_email_id, created_at)
VALUES ('12345678-0009-0009-0009-000000000201'::uuid, 'DOMAIN', 'ONBOARDED', 'applyopening.example', 'admin@applyopening.example', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0009-0009-0009-000000000011'::uuid, timezone('UTC'::text, now()));

INSERT INTO domains (id, domain_name, domain_state, employer_id, created_at)
VALUES ('12345678-0009-0009-0009-000000003001'::uuid, 'applyopening.example', 'VERIFIED', '12345678-0009-0009-0009-000000000201'::uuid, timezone('UTC'::text, now()));

-- Set primary domain
INSERT INTO employer_primary_domains (employer_id, domain_id)
VALUES ('12345678-0009-0009-0009-000000000201'::uuid, '12345678-0009-0009-0009-000000003001'::uuid);

-- Create org users
INSERT INTO org_users (id, email, name, password_hash, org_user_roles, org_user_state, employer_id, created_at)
VALUES ('12345678-0009-0009-0009-000000040001'::uuid, 'admin@applyopening.example', 'Admin User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0009-0009-0009-000000000201'::uuid, timezone('UTC'::text, now()));

-- Create cost center
INSERT INTO org_cost_centers (id, cost_center_name, cost_center_state, notes, employer_id, created_at)
VALUES ('12345678-0009-0009-0009-000000050001'::uuid, 'Engineering', 'ACTIVE_CC', 'Engineering department', '12345678-0009-0009-0009-000000000201'::uuid, timezone('UTC'::text, now()));

-- Create location
INSERT INTO locations (id, title, country_code, postal_address, postal_code, city_aka, location_state, employer_id, created_at)
VALUES ('12345678-0009-0009-0009-000000060001'::uuid, 'Bangalore Office', 'IND', '123 MG Road', '560001', ARRAY['Bengaluru'], 'ACTIVE_LOCATION', '12345678-0009-0009-0009-000000000201'::uuid, timezone('UTC'::text, now()));

-- Create test opening
INSERT INTO openings (
    employer_id,
    id,
    title,
    positions,
    jd,
    recruiter,
    hiring_manager,
    cost_center_id,
    opening_type,
    yoe_min,
    yoe_max,
    min_education_level,
    salary_min,
    salary_max,
    salary_currency,
    state,
    created_at,
    last_updated_at
)
VALUES (
    '12345678-0009-0009-0009-000000000201'::uuid,
    '2024-Mar-09-001',
    'Software Engineer',
    2,
    'Looking for talented engineers...',
    '12345678-0009-0009-0009-000000040001'::uuid,
    '12345678-0009-0009-0009-000000040001'::uuid,
    '12345678-0009-0009-0009-000000050001'::uuid,
    'FULL_TIME_OPENING',
    2,
    5,
    'BACHELOR_EDUCATION',
    50000,
    100000,
    'USD',
    'ACTIVE_OPENING_STATE',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
);

-- Link opening to location
INSERT INTO opening_locations (employer_id, opening_id, location_id)
VALUES ('12345678-0009-0009-0009-000000000201'::uuid, '2024-Mar-09-001', '12345678-0009-0009-0009-000000060001'::uuid);

COMMIT;
