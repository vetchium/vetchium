BEGIN;
--- email table primary key uuids should end in 2 digits, 11, 12, 13, etc
--- employer table primary key uuids should end in 3 digits, 201, 202, 203, etc
--- domain table primary key uuids should end in 4 digits, 3001, 3002, 3003, etc
--- org_users table primary key uuids should end in 5 digits, 40001, 40002, 40003, etc
--- hub_users table primary key uuids should end in 6 digits, 60001, 60002, 60003, etc
--- applications table ids should be sequential strings
--- candidacies table ids should be sequential strings

-- Create employer and related records
INSERT INTO emails (email_key, email_from, email_to, email_cc, email_bcc, email_subject, email_html_body, email_text_body, email_state, created_at, processed_at)
VALUES ('12345678-0011-0011-0011-000000000011'::uuid, 'no-reply@vetchi.org', ARRAY['admin@filter-candidacy-infos.example'], NULL, NULL, 'Welcome to Vetchium', 'Welcome HTML', 'Welcome Text', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now()));

INSERT INTO employers (id, client_id_type, employer_state, company_name, onboard_admin_email, onboard_secret_token, token_valid_till, onboard_email_id, created_at)
VALUES ('12345678-0011-0011-0011-000000000201'::uuid, 'DOMAIN', 'ONBOARDED', 'filter-candidacy-infos.example', 'admin@filter-candidacy-infos.example', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0011-0011-0011-000000000011'::uuid, timezone('UTC'::text, now()));

INSERT INTO domains (id, domain_name, domain_state, employer_id, created_at)
VALUES ('12345678-0011-0011-0011-000000003001'::uuid, 'filter-candidacy-infos.example', 'VERIFIED', '12345678-0011-0011-0011-000000000201'::uuid, timezone('UTC'::text, now()));

-- Set primary domain
INSERT INTO employer_primary_domains (employer_id, domain_id)
VALUES ('12345678-0011-0011-0011-000000000201'::uuid, '12345678-0011-0011-0011-000000003001'::uuid);

-- Create org users with different roles
INSERT INTO org_users (id, email, name, password_hash, org_user_roles, org_user_state, employer_id, created_at)
VALUES 
    ('12345678-0011-0011-0011-000000040001'::uuid, 'admin@filter-candidacy-infos.example', 'Admin User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0011-0011-0011-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0011-0011-0011-000000040002'::uuid, 'recruiter1@filter-candidacy-infos.example', 'Recruiter One', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['OPENINGS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0011-0011-0011-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0011-0011-0011-000000040003'::uuid, 'recruiter2@filter-candidacy-infos.example', 'Recruiter Two', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['OPENINGS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0011-0011-0011-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0011-0011-0011-000000040004'::uuid, 'viewer@filter-candidacy-infos.example', 'Viewer User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['OPENINGS_VIEWER']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0011-0011-0011-000000000201'::uuid, timezone('UTC'::text, now()));

-- Create hub users (applicants)
INSERT INTO hub_users (id, full_name, handle, email, password_hash, state, tier, resident_country_code, resident_city, preferred_language, short_bio, long_bio, created_at, updated_at)
VALUES
    ('12345678-0011-0011-0011-000000060001'::uuid, 'Applicant One', 'applicant1', 'applicant1@hub.example', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'USA', 'New York', 'en', 'Applicant One is innovative', 'Applicant One was born in USA and finished education at NYU and has 4 years as experience.', timezone('UTC'::text, now()), timezone('UTC'::text, now())),
    ('12345678-0011-0011-0011-000000060002'::uuid, 'Applicant Two', 'applicant2', 'applicant2@hub.example', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'IND', 'Bangalore', 'en', 'Applicant Two is dedicated', 'Applicant Two was born in India and finished education at IISc Bangalore and has 6 years as experience.', timezone('UTC'::text, now()), timezone('UTC'::text, now())),
    ('12345678-0011-0011-0011-000000060003'::uuid, 'Applicant Three', 'applicant3', 'applicant3@hub.example', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'GBR', 'London', 'en', 'Applicant Three is resourceful', 'Applicant Three was born in UK and finished education at Imperial College London and has 5 years as experience.', timezone('UTC'::text, now()), timezone('UTC'::text, now()));

-- Create cost centers
INSERT INTO org_cost_centers (id, cost_center_name, cost_center_state, notes, employer_id, created_at)
VALUES ('12345678-0011-0011-0011-000000050001'::uuid, 'Engineering', 'ACTIVE_CC', 'Engineering department', '12345678-0011-0011-0011-000000000201'::uuid, timezone('UTC'::text, now()));

-- Create openings
INSERT INTO openings (employer_id, id, title, positions, jd, recruiter, hiring_manager, cost_center_id, opening_type, yoe_min, yoe_max, min_education_level, state, created_at, last_updated_at)
VALUES
    ('12345678-0011-0011-0011-000000000201'::uuid, '2024-Mar-01-001', 'Software Engineer', 2, 'Looking for talented engineers', '12345678-0011-0011-0011-000000040002'::uuid, '12345678-0011-0011-0011-000000040001'::uuid, '12345678-0011-0011-0011-000000050001'::uuid, 'FULL_TIME_OPENING', 2, 5, 'BACHELOR_EDUCATION', 'ACTIVE_OPENING_STATE', timezone('UTC'::text, now()), timezone('UTC'::text, now())),
    ('12345678-0011-0011-0011-000000000201'::uuid, '2024-Mar-01-002', 'Senior Engineer', 1, 'Looking for senior engineers', '12345678-0011-0011-0011-000000040003'::uuid, '12345678-0011-0011-0011-000000040001'::uuid, '12345678-0011-0011-0011-000000050001'::uuid, 'FULL_TIME_OPENING', 5, 8, 'BACHELOR_EDUCATION', 'ACTIVE_OPENING_STATE', timezone('UTC'::text, now()), timezone('UTC'::text, now()));

-- Create applications
INSERT INTO applications (id, employer_id, opening_id, cover_letter, resume_sha, application_state, color_tag, hub_user_id, created_at)
VALUES
    ('APP-001', '12345678-0011-0011-0011-000000000201'::uuid, '2024-Mar-01-001', 'Cover letter 1', 'sha-sha-sha', 'SHORTLISTED', 'GREEN', '12345678-0011-0011-0011-000000060001'::uuid, timezone('UTC'::text, now())),
    ('APP-002', '12345678-0011-0011-0011-000000000201'::uuid, '2024-Mar-01-001', 'Cover letter 2', 'sha-sha-sha', 'SHORTLISTED', 'YELLOW', '12345678-0011-0011-0011-000000060002'::uuid, timezone('UTC'::text, now())),
    ('APP-003', '12345678-0011-0011-0011-000000000201'::uuid, '2024-Mar-01-002', 'Cover letter 3', 'sha-sha-sha', 'SHORTLISTED', 'GREEN', '12345678-0011-0011-0011-000000060003'::uuid, timezone('UTC'::text, now()));

-- Create candidacies with different states
INSERT INTO candidacies (id, application_id, employer_id, opening_id, candidacy_state, created_by, created_at)
VALUES
    ('CAND-001', 'APP-001', '12345678-0011-0011-0011-000000000201'::uuid, '2024-Mar-01-001', 'INTERVIEWING', '12345678-0011-0011-0011-000000040002'::uuid, timezone('UTC'::text, now())),
    ('CAND-002', 'APP-002', '12345678-0011-0011-0011-000000000201'::uuid, '2024-Mar-01-001', 'OFFERED', '12345678-0011-0011-0011-000000040002'::uuid, timezone('UTC'::text, now())),
    ('CAND-003', 'APP-003', '12345678-0011-0011-0011-000000000201'::uuid, '2024-Mar-01-002', 'INTERVIEWING', '12345678-0011-0011-0011-000000040003'::uuid, timezone('UTC'::text, now()));

COMMIT;
