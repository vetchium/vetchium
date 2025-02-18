BEGIN;

-- Create employer for testing
INSERT INTO emails (email_key, email_from, email_to, email_cc, email_bcc, email_subject, email_html_body, email_text_body, email_state, created_at, processed_at)
VALUES ('12345678-0011-0011-0011-000000000011'::uuid, 'no-reply@vetchi.org', ARRAY['admin@candidacy-comments.example'], NULL, NULL, 'Welcome to Vetchi Subject', 'Welcome HTML', 'Welcome Text', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now()));

INSERT INTO employers (id, client_id_type, employer_state, company_name, onboard_admin_email, onboard_secret_token, token_valid_till, onboard_email_id, created_at)
VALUES ('12345678-0011-0011-0011-000000000201'::uuid, 'DOMAIN', 'ONBOARDED', 'candidacy-comments.example', 'admin@candidacy-comments.example', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0011-0011-0011-000000000011'::uuid, timezone('UTC'::text, now()));

INSERT INTO domains (id, domain_name, domain_state, employer_id, created_at)
VALUES ('12345678-0011-0011-0011-000000003001'::uuid, 'candidacy-comments.example', 'VERIFIED', '12345678-0011-0011-0011-000000000201'::uuid, timezone('UTC'::text, now()));

-- Set primary domain
INSERT INTO employer_primary_domains (employer_id, domain_id)
VALUES ('12345678-0011-0011-0011-000000000201'::uuid, '12345678-0011-0011-0011-000000003001'::uuid);

-- Create org users with different roles
INSERT INTO org_users (id, email, name, password_hash, org_user_roles, org_user_state, employer_id, created_at)
VALUES
    ('12345678-0011-0011-0011-000000040001'::uuid, 'admin@candidacy-comments.example', 'Admin User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0011-0011-0011-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0011-0011-0011-000000040002'::uuid, 'hiringmanager@candidacy-comments.example', 'Hiring Manager', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['APPLICATIONS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0011-0011-0011-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0011-0011-0011-000000040003'::uuid, 'recruiter@candidacy-comments.example', 'Recruiter', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['APPLICATIONS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0011-0011-0011-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0011-0011-0011-000000040004'::uuid, 'watcher@candidacy-comments.example', 'Watcher', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['APPLICATIONS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0011-0011-0011-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0011-0011-0011-000000040005'::uuid, 'regular@candidacy-comments.example', 'Regular User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['APPLICATIONS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0011-0011-0011-000000000201'::uuid, timezone('UTC'::text, now()));

-- Create hub users for testing
INSERT INTO hub_users (id, full_name, handle, email, password_hash, state, resident_country_code, resident_city, preferred_language, short_bio, long_bio, created_at, updated_at)
VALUES
    ('12345678-0011-0011-0011-000000050001'::uuid, 'Active Hub User', 'active_hub_user', '0011-active@hub.example', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'USA', 'New York', 'en', 'Active Hub User is analytical', 'Active Hub User was born in USA and finished education at Columbia University and has 4 years as experience.', timezone('UTC'::text, now()), timezone('UTC'::text, now()));

-- Create cost center for openings
INSERT INTO org_cost_centers (id, cost_center_name, cost_center_state, notes, employer_id, created_at)
VALUES ('12345678-0011-0011-0011-000000060001'::uuid, 'Engineering', 'ACTIVE_CC', 'Engineering Department', '12345678-0011-0011-0011-000000000201'::uuid, timezone('UTC'::text, now()));

-- Create test openings
INSERT INTO openings (
    id,
    employer_id,
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
    state,
    created_at
)
VALUES (
    '2024-Mar-11-001',
    '12345678-0011-0011-0011-000000000201'::uuid,
    'Software Engineer',
    1,
    'Test Opening',
    '12345678-0011-0011-0011-000000040003'::uuid,
    '12345678-0011-0011-0011-000000040002'::uuid,
    '12345678-0011-0011-0011-000000060001'::uuid,
    'FULL_TIME_OPENING',
    2,
    5,
    'BACHELOR_EDUCATION',
    'ACTIVE_OPENING_STATE',
    timezone('UTC'::text, now())
);

-- Add watcher to the test opening
INSERT INTO opening_watchers (
    employer_id,
    opening_id,
    watcher_id
)
VALUES (
    '12345678-0011-0011-0011-000000000201'::uuid,
    '2024-Mar-11-001',
    '12345678-0011-0011-0011-000000040004'::uuid  -- watcher@candidacy-comments.example
);

-- Create test applications
INSERT INTO applications (id, employer_id, opening_id, cover_letter, resume_sha, application_state, hub_user_id, created_at)
VALUES
    ('2024-Dec-01-1', '12345678-0011-0011-0011-000000000201'::uuid, '2024-Mar-11-001', 'Test Cover Letter', 'sha-sha-sha', 'APPLIED', '12345678-0011-0011-0011-000000050001'::uuid, timezone('UTC'::text, now()));

-- Create test candidacies
INSERT INTO candidacies (id, application_id, employer_id, opening_id, candidacy_state, created_by, created_at)
VALUES
    ('12345678-0011-0011-0011-000000060001', '2024-Dec-01-1', '12345678-0011-0011-0011-000000000201'::uuid, '2024-Mar-11-001', 'INTERVIEWING', '12345678-0011-0011-0011-000000040001'::uuid, timezone('UTC'::text, now()));

COMMIT;
