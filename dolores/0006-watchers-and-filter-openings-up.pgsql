BEGIN;
--- email table primary key uuids should end in 2 digits, 11, 12, 13, etc
--- employer table primary key uuids should end in 3 digits, 201, 202, 203, etc
--- domain table primary key uuids should end in 4 digits, 3001, 3002, 3003, etc
--- org_users table primary key uuids should end in 5 digits, 40001, 40002, 40003, etc
--- cost_centers table primary key uuids should end in 6 digits, 50001, 50002, 50003, etc
--- locations table primary key uuids should end in 7 digits, 60001, 60002, 60003, etc
--- openings table ids should be sequential strings

INSERT INTO public.emails (email_key, email_from, email_to, email_cc, email_bcc, email_subject, email_html_body, email_text_body, email_state, created_at, processed_at)
    VALUES ('12345678-0006-0006-0006-000000000011'::uuid, 'no-reply@vetchi.org', ARRAY['admin@openings0006.example'], NULL, NULL, 'Welcome to Vetchi Subject', 'Welcome to Vetchi HTML Body', 'Welcome to Vetchi Text Body', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now()));

INSERT INTO public.employers (id, client_id_type, employer_state, onboard_admin_email, onboard_secret_token, token_valid_till, onboard_email_id, created_at)
    VALUES ('12345678-0006-0006-0006-000000000201'::uuid, 'DOMAIN', 'ONBOARDED', 'admin@openings0006.example', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0006-0006-0006-000000000011'::uuid, timezone('UTC'::text, now()));

INSERT INTO public.domains (id, domain_name, domain_state, employer_id, created_at)
    VALUES ('12345678-0006-0006-0006-000000003001'::uuid, 'openings0006.example', 'VERIFIED', '12345678-0006-0006-0006-000000000201'::uuid, timezone('UTC'::text, now()));

-- Insert users with different roles
INSERT INTO public.org_users (id, email, name, password_hash, org_user_roles, org_user_state, employer_id, created_at)
    VALUES 
    ('12345678-0006-0006-0006-000000040001'::uuid, 'admin@openings0006.example', 'Admin User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0006-0006-0006-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000040002'::uuid, 'crud@openings0006.example', 'CRUD User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['OPENINGS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0006-0006-0006-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000040003'::uuid, 'viewer@openings0006.example', 'Viewer User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['OPENINGS_VIEWER']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0006-0006-0006-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000040004'::uuid, 'recruiter@openings0006.example', 'Recruiter User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['OPENINGS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0006-0006-0006-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000040005'::uuid, 'hiring-manager@openings0006.example', 'Hiring Manager User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['OPENINGS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0006-0006-0006-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000040006'::uuid, 'non-openings@openings0006.example', 'Non Openings User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['COST_CENTERS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0006-0006-0006-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000040007'::uuid, 'watcher1@openings0006.example', 'Watcher One', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['OPENINGS_VIEWER']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0006-0006-0006-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000040008'::uuid, 'watcher2@openings0006.example', 'Watcher Two', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['OPENINGS_VIEWER']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0006-0006-0006-000000000201'::uuid, timezone('UTC'::text, now()));

-- Insert cost centers
INSERT INTO public.org_cost_centers (id, cost_center_name, cost_center_state, notes, employer_id, created_at)
    VALUES 
    ('12345678-0006-0006-0006-000000050001'::uuid, 'Engineering', 'ACTIVE_CC', 'Engineering department', '12345678-0006-0006-0006-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000050002'::uuid, 'Sales', 'ACTIVE_CC', 'Sales department', '12345678-0006-0006-0006-000000000201'::uuid, timezone('UTC'::text, now()));

-- Insert locations
INSERT INTO public.locations (id, title, country_code, postal_address, postal_code, openstreetmap_url, city_aka, location_state, employer_id, created_at)
    VALUES 
    ('12345678-0006-0006-0006-000000060001'::uuid, 'Bangalore Office', 'IND', '123 MG Road, Bangalore', '560001', NULL, ARRAY['Bengaluru', 'Silicon Valley of India'], 'ACTIVE_LOCATION', '12345678-0006-0006-0006-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000060002'::uuid, 'Chennai Office', 'IND', '456 Anna Salai, Chennai', '600002', NULL, ARRAY['Madras'], 'ACTIVE_LOCATION', '12345678-0006-0006-0006-000000000201'::uuid, timezone('UTC'::text, now()));

-- Insert openings with different states and dates
INSERT INTO public.openings (employer_id, id, title, positions, jd, recruiter, hiring_manager, cost_center_id, opening_type, yoe_min, yoe_max, min_education_level, current_state, created_at, last_updated_at)
VALUES
    -- February 15 (1 opening)
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Feb-15-001', 'Software Engineer', 2, 'JD for Software Engineer', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050001'::uuid, 'FULL_TIME_OPENING', 2, 5, 'BACHELOR_EDUCATION', 'DRAFT_OPENING_STATE', timezone('UTC'::text, now()) - interval '30 days', timezone('UTC'::text, now())),

    -- February 25 (2 openings)
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Feb-25-001', 'Senior Engineer', 1, 'JD for Senior Engineer', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050001'::uuid, 'FULL_TIME_OPENING', 5, 8, 'BACHELOR_EDUCATION', 'ACTIVE_OPENING_STATE', timezone('UTC'::text, now()) - interval '20 days', timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Feb-25-002', 'Junior Engineer', 2, 'JD for Junior Engineer', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050001'::uuid, 'FULL_TIME_OPENING', 0, 2, 'BACHELOR_EDUCATION', 'CLOSED_OPENING_STATE', timezone('UTC'::text, now()) - interval '20 days', timezone('UTC'::text, now())),

    -- March 1 (5 openings)
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-01-001', 'Product Manager', 1, 'JD for Product Manager', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050001'::uuid, 'FULL_TIME_OPENING', 5, 8, 'BACHELOR_EDUCATION', 'ACTIVE_OPENING_STATE', timezone('UTC'::text, now()) - interval '15 days', timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-01-002', 'UX Designer', 2, 'JD for UX Designer', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050001'::uuid, 'FULL_TIME_OPENING', 3, 6, 'BACHELOR_EDUCATION', 'DRAFT_OPENING_STATE', timezone('UTC'::text, now()) - interval '15 days', timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-01-003', 'DevOps Engineer', 1, 'JD for DevOps', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050001'::uuid, 'FULL_TIME_OPENING', 4, 8, 'BACHELOR_EDUCATION', 'SUSPENDED_OPENING_STATE', timezone('UTC'::text, now()) - interval '15 days', timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-01-004', 'QA Engineer', 2, 'JD for QA', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050001'::uuid, 'FULL_TIME_OPENING', 2, 5, 'BACHELOR_EDUCATION', 'CLOSED_OPENING_STATE', timezone('UTC'::text, now()) - interval '15 days', timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-01-005', 'Technical Writer', 1, 'JD for Tech Writer', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050001'::uuid, 'FULL_TIME_OPENING', 1, 3, 'BACHELOR_EDUCATION', 'ACTIVE_OPENING_STATE', timezone('UTC'::text, now()) - interval '15 days', timezone('UTC'::text, now())),

    -- March 6 (15 openings)
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-001', 'Sales Executive', 3, 'JD for Sales Executive', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050002'::uuid, 'FULL_TIME_OPENING', 1, 3, 'BACHELOR_EDUCATION', 'SUSPENDED_OPENING_STATE', timezone('UTC'::text, now()) - interval '10 days', timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-002', 'Sales Manager', 1, 'JD for Sales Manager', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050002'::uuid, 'FULL_TIME_OPENING', 5, 8, 'BACHELOR_EDUCATION', 'ACTIVE_OPENING_STATE', timezone('UTC'::text, now()) - interval '10 days', timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-003', 'Account Executive', 2, 'JD for Account Executive', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050002'::uuid, 'FULL_TIME_OPENING', 2, 4, 'BACHELOR_EDUCATION', 'DRAFT_OPENING_STATE', timezone('UTC'::text, now()) - interval '10 days', timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-004', 'Marketing Manager', 1, 'JD for Marketing Manager', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050002'::uuid, 'FULL_TIME_OPENING', 6, 10, 'BACHELOR_EDUCATION', 'CLOSED_OPENING_STATE', timezone('UTC'::text, now()) - interval '10 days', timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-005', 'Business Analyst', 2, 'JD for Business Analyst', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050002'::uuid, 'FULL_TIME_OPENING', 3, 6, 'BACHELOR_EDUCATION', 'ACTIVE_OPENING_STATE', timezone('UTC'::text, now()) - interval '10 days', timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-006', 'Data Analyst', 2, 'JD for Data Analyst', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050002'::uuid, 'FULL_TIME_OPENING', 2, 5, 'BACHELOR_EDUCATION', 'SUSPENDED_OPENING_STATE', timezone('UTC'::text, now()) - interval '10 days', timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-007', 'HR Manager', 1, 'JD for HR Manager', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050002'::uuid, 'FULL_TIME_OPENING', 5, 8, 'BACHELOR_EDUCATION', 'DRAFT_OPENING_STATE', timezone('UTC'::text, now()) - interval '10 days', timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-008', 'HR Executive', 2, 'JD for HR Executive', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050002'::uuid, 'FULL_TIME_OPENING', 1, 3, 'BACHELOR_EDUCATION', 'CLOSED_OPENING_STATE', timezone('UTC'::text, now()) - interval '10 days', timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-009', 'Office Manager', 1, 'JD for Office Manager', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050002'::uuid, 'FULL_TIME_OPENING', 3, 6, 'BACHELOR_EDUCATION', 'ACTIVE_OPENING_STATE', timezone('UTC'::text, now()) - interval '10 days', timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-010', 'Admin Assistant', 2, 'JD for Admin Assistant', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050002'::uuid, 'FULL_TIME_OPENING', 0, 2, 'BACHELOR_EDUCATION', 'SUSPENDED_OPENING_STATE', timezone('UTC'::text, now()) - interval '10 days', timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-011', 'Project Manager', 1, 'JD for Project Manager', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050002'::uuid, 'FULL_TIME_OPENING', 5, 8, 'BACHELOR_EDUCATION', 'ACTIVE_OPENING_STATE', timezone('UTC'::text, now()) - interval '10 days', timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-012', 'Technical Lead', 1, 'JD for Technical Lead', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050002'::uuid, 'FULL_TIME_OPENING', 8, 12, 'BACHELOR_EDUCATION', 'DRAFT_OPENING_STATE', timezone('UTC'::text, now()) - interval '10 days', timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-013', 'System Architect', 1, 'JD for System Architect', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050002'::uuid, 'FULL_TIME_OPENING', 10, 15, 'BACHELOR_EDUCATION', 'SUSPENDED_OPENING_STATE', timezone('UTC'::text, now()) - interval '10 days', timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-014', 'Database Administrator', 2, 'JD for DBA', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050002'::uuid, 'FULL_TIME_OPENING', 5, 8, 'BACHELOR_EDUCATION', 'ACTIVE_OPENING_STATE', timezone('UTC'::text, now()) - interval '10 days', timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-015', 'Security Engineer', 1, 'JD for Security Engineer', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050002'::uuid, 'FULL_TIME_OPENING', 3, 6, 'BACHELOR_EDUCATION', 'DRAFT_OPENING_STATE', timezone('UTC'::text, now()) - interval '10 days', timezone('UTC'::text, now())),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-11-001', 'Tech Lead', 1, 'JD for Tech Lead', '12345678-0006-0006-0006-000000040004'::uuid, '12345678-0006-0006-0006-000000040005'::uuid, '12345678-0006-0006-0006-000000050001'::uuid, 'FULL_TIME_OPENING', 8, 12, 'BACHELOR_EDUCATION', 'CLOSED_OPENING_STATE', timezone('UTC'::text, now()) - interval '5 days', timezone('UTC'::text, now()));

-- Insert opening locations
INSERT INTO public.opening_locations (employer_id, opening_id, location_id)
VALUES
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Feb-15-001', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Feb-25-001', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Feb-25-002', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-01-001', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-01-002', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-01-003', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-01-004', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-01-005', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-001', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-002', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-003', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-004', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-005', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-006', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-007', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-008', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-009', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-010', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-011', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-012', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-013', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-014', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-015', '12345678-0006-0006-0006-000000060001'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-11-001', '12345678-0006-0006-0006-000000060002'::uuid);

-- Insert opening watchers
INSERT INTO public.opening_watchers (employer_id, opening_id, watcher_id)
VALUES
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Feb-15-001', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Feb-15-001', '12345678-0006-0006-0006-000000040008'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Feb-25-001', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Feb-25-002', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-01-001', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-01-002', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-01-003', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-01-004', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-01-005', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-001', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-002', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-003', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-004', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-005', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-006', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-007', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-008', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-009', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-010', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-011', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-012', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-013', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-014', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-06-015', '12345678-0006-0006-0006-000000040007'::uuid),
    ('12345678-0006-0006-0006-000000000201'::uuid, '2024-Mar-11-001', '12345678-0006-0006-0006-000000040007'::uuid);

COMMIT;
