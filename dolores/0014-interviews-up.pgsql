BEGIN;
-- Create employer
INSERT INTO public.emails (email_key, email_from, email_to, email_cc, email_bcc, email_subject, email_html_body, email_text_body, email_state, created_at, processed_at)
    VALUES ('12345678-0014-0014-0014-000000000011'::uuid, 'no-reply@vetchi.org', ARRAY['admin@interview.example'], NULL, NULL, 'Welcome to Vetchi Subject', 'Welcome to Vetchi HTML Body', 'Welcome to Vetchi Text Body', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now()));

INSERT INTO public.employers (id, client_id_type, employer_state, company_name, onboard_admin_email, onboard_secret_token, token_valid_till, onboard_email_id, created_at)
    VALUES ('12345678-0014-0014-0014-000000000201'::uuid, 'DOMAIN', 'ONBOARDED', 'interview.example', 'admin@interview.example', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0014-0014-0014-000000000011'::uuid, timezone('UTC'::text, now()));

INSERT INTO public.domains (id, domain_name, domain_state, employer_id, created_at)
    VALUES ('12345678-0014-0014-0014-000000003001'::uuid, 'interview.example', 'VERIFIED', '12345678-0014-0014-0014-000000000201'::uuid, timezone('UTC'::text, now()));

-- Insert users with different roles
INSERT INTO public.org_users (id, email, name, password_hash, org_user_roles, org_user_state, employer_id, created_at)
    VALUES 
    ('12345678-0014-0014-0014-000000040001'::uuid, 'admin@interview.example', 'Admin User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0014-0014-0014-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0014-0014-0014-000000040002'::uuid, 'recruiter@interview.example', 'Recruiter User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['OPENINGS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0014-0014-0014-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0014-0014-0014-000000040003'::uuid, 'hiring-manager@interview.example', 'Hiring Manager User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['OPENINGS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0014-0014-0014-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0014-0014-0014-000000040004'::uuid, 'interviewer1@interview.example', 'Interviewer One', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['OPENINGS_VIEWER']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0014-0014-0014-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0014-0014-0014-000000040005'::uuid, 'interviewer2@interview.example', 'Interviewer Two', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['OPENINGS_VIEWER']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0014-0014-0014-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0014-0014-0014-000000040006'::uuid, 'interviewer3@interview.example', 'Interviewer Three', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['OPENINGS_VIEWER']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0014-0014-0014-000000000201'::uuid, timezone('UTC'::text, now()));

-- Insert cost centers
INSERT INTO public.org_cost_centers (id, cost_center_name, cost_center_state, notes, employer_id, created_at)
    VALUES ('12345678-0014-0014-0014-000000050001'::uuid, 'Engineering', 'ACTIVE_CC', 'Engineering department', '12345678-0014-0014-0014-000000000201'::uuid, timezone('UTC'::text, now()));

-- Insert locations
INSERT INTO public.locations (id, title, country_code, postal_address, postal_code, openstreetmap_url, city_aka, location_state, employer_id, created_at)
    VALUES ('12345678-0014-0014-0014-000000060001'::uuid, 'Main Office', 'IND', '123 Main St', '600001', NULL, ARRAY['Chennai', 'Madras'], 'ACTIVE_LOCATION', '12345678-0014-0014-0014-000000000201'::uuid, timezone('UTC'::text, now()));

-- Create hub user
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
) VALUES (
    '12345678-0014-0014-0014-000000050001'::uuid,
    'Interview Test User',
    'interview',
    'interview@hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'IND',
    'Chennai',
    'en',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
);

-- Create opening
INSERT INTO openings (
    employer_id,
    id,
    title,
    positions,
    description,
    recruiter_id,
    hiring_manager_id,
    cost_center_id,
    opening_type,
    yoe_min,
    yoe_max,
    education_level,
    salary_min,
    salary_max,
    salary_currency,
    remote_country_codes,
    remote_timezones,
    opening_state,
    created_at,
    updated_at
) VALUES (
    '12345678-0014-0014-0014-000000000201'::uuid,
    '2024-Mar-15-001',
    'Software Engineer',
    2,
    'Looking for talented engineers...',
    '12345678-0014-0014-0014-000000040002'::uuid,
    '12345678-0014-0014-0014-000000040003'::uuid,
    '12345678-0014-0014-0014-000000050001'::uuid,
    'FULL_TIME_OPENING',
    2,
    5,
    'BACHELOR_EDUCATION',
    80000,
    120000,
    'USD',
    ARRAY['IND'],
    ARRAY['IST Indian Standard Time GMT+0530'],
    'ACTIVE_OPENING_STATE',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
);

-- Link opening to location
INSERT INTO opening_locations (employer_id, opening_id, location_id)
VALUES ('12345678-0014-0014-0014-000000000201'::uuid, '2024-Mar-15-001', '12345678-0014-0014-0014-000000060001'::uuid);

COMMIT; 