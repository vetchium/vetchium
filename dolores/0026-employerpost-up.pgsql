BEGIN;

-- Create employers for testing
INSERT INTO emails (email_key, email_from, email_to, email_cc, email_bcc, email_subject, email_html_body, email_text_body, email_state, created_at, processed_at)
VALUES 
    ('12345678-0026-0026-0026-000000000011'::uuid, 'no-reply@vetchi.org', ARRAY['admin@0026-employerposts.example.com'], NULL, NULL, 'Welcome to Vetchium Subject', 'Welcome HTML', 'Welcome Text', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000000012'::uuid, 'no-reply@vetchi.org', ARRAY['admin@0026-employerposts2.example.com'], NULL, NULL, 'Welcome to Vetchium Subject', 'Welcome HTML', 'Welcome Text', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000000013'::uuid, 'no-reply@vetchi.org', ARRAY['admin@0026-employerposts3.example.com'], NULL, NULL, 'Welcome to Vetchium Subject', 'Welcome HTML', 'Welcome Text', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000000014'::uuid, 'no-reply@vetchi.org', ARRAY['admin@0026-employerposts4.example.com'], NULL, NULL, 'Welcome to Vetchium Subject', 'Welcome HTML', 'Welcome Text', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now()));

INSERT INTO employers (id, client_id_type, employer_state, company_name, onboard_admin_email, onboard_secret_token, token_valid_till, onboard_email_id, created_at)
VALUES 
    ('12345678-0026-0026-0026-000000000201'::uuid, 'DOMAIN', 'ONBOARDED', '0026-employerposts.example.com', 'admin@0026-employerposts.example.com', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0026-0026-0026-000000000011'::uuid, timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000000202'::uuid, 'DOMAIN', 'ONBOARDED', '0026-employerposts2.example.com', 'admin@0026-employerposts2.example.com', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0026-0026-0026-000000000012'::uuid, timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000000203'::uuid, 'DOMAIN', 'ONBOARDED', '0026-employerposts3.example.com', 'admin@0026-employerposts3.example.com', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0026-0026-0026-000000000013'::uuid, timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000000204'::uuid, 'DOMAIN', 'ONBOARDED', '0026-employerposts4.example.com', 'admin@0026-employerposts4.example.com', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0026-0026-0026-000000000014'::uuid, timezone('UTC'::text, now()));

INSERT INTO domains (id, domain_name, domain_state, employer_id, created_at)
VALUES 
    ('12345678-0026-0026-0026-000000003001'::uuid, '0026-employerposts.example.com', 'VERIFIED', '12345678-0026-0026-0026-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000003002'::uuid, '0026-employerposts2.example.com', 'VERIFIED', '12345678-0026-0026-0026-000000000202'::uuid, timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000003003'::uuid, '0026-employerposts3.example.com', 'VERIFIED', '12345678-0026-0026-0026-000000000203'::uuid, timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000003004'::uuid, '0026-employerposts4.example.com', 'VERIFIED', '12345678-0026-0026-0026-000000000204'::uuid, timezone('UTC'::text, now()));

-- Set primary domains
INSERT INTO employer_primary_domains (employer_id, domain_id)
VALUES 
    ('12345678-0026-0026-0026-000000000201'::uuid, '12345678-0026-0026-0026-000000003001'::uuid),
    ('12345678-0026-0026-0026-000000000202'::uuid, '12345678-0026-0026-0026-000000003002'::uuid),
    ('12345678-0026-0026-0026-000000000203'::uuid, '12345678-0026-0026-0026-000000003003'::uuid),
    ('12345678-0026-0026-0026-000000000204'::uuid, '12345678-0026-0026-0026-000000003004'::uuid);

-- Create org users with different roles for each employer
INSERT INTO org_users (id, email, name, password_hash, org_user_roles, org_user_state, employer_id, created_at)
VALUES
    -- Employer 1 users
    ('12345678-0026-0026-0026-000000040001'::uuid, 'admin@0026-employerposts.example.com', 'Admin User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0026-0026-0026-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000040002'::uuid, 'marketing@0026-employerposts.example.com', 'Marketing User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['EMPLOYER_POSTS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0026-0026-0026-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000040003'::uuid, 'regular@0026-employerposts.example.com', 'Regular User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['APPLICATIONS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0026-0026-0026-000000000201'::uuid, timezone('UTC'::text, now())),
    
    -- Employer 2 users
    ('12345678-0026-0026-0026-000000040004'::uuid, 'admin@0026-employerposts2.example.com', 'Admin User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0026-0026-0026-000000000202'::uuid, timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000040005'::uuid, 'marketing@0026-employerposts2.example.com', 'Marketing User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['EMPLOYER_POSTS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0026-0026-0026-000000000202'::uuid, timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000040006'::uuid, 'regular@0026-employerposts2.example.com', 'Regular User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['APPLICATIONS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0026-0026-0026-000000000202'::uuid, timezone('UTC'::text, now())),
    
    -- Employer 3 users
    ('12345678-0026-0026-0026-000000040007'::uuid, 'admin@0026-employerposts3.example.com', 'Admin User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0026-0026-0026-000000000203'::uuid, timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000040008'::uuid, 'marketing@0026-employerposts3.example.com', 'Marketing User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['EMPLOYER_POSTS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0026-0026-0026-000000000203'::uuid, timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000040009'::uuid, 'regular@0026-employerposts3.example.com', 'Regular User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['APPLICATIONS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0026-0026-0026-000000000203'::uuid, timezone('UTC'::text, now())),
    
    -- Employer 4 users
    ('12345678-0026-0026-0026-000000040010'::uuid, 'admin@0026-employerposts4.example.com', 'Admin User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0026-0026-0026-000000000204'::uuid, timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000040011'::uuid, 'marketing@0026-employerposts4.example.com', 'Marketing User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['EMPLOYER_POSTS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0026-0026-0026-000000000204'::uuid, timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000040012'::uuid, 'regular@0026-employerposts4.example.com', 'Regular User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['APPLICATIONS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0026-0026-0026-000000000204'::uuid, timezone('UTC'::text, now()));

-- Add tags for testing
INSERT INTO tags (id, name)
VALUES
    ('12345678-0026-0026-0026-000000050001'::uuid, '0026-engineering'),
    ('12345678-0026-0026-0026-000000050002'::uuid, '0026-marketing'),
    ('12345678-0026-0026-0026-000000050003'::uuid, '0026-golang'),
    ('12345678-0026-0026-0026-000000050004'::uuid, '0026-react'),
    ('12345678-0026-0026-0026-000000050005'::uuid, '0026-testing')
ON CONFLICT (name) DO NOTHING;

-- Create test hub users for org following tests
INSERT INTO hub_users (id, full_name, handle, email, password_hash, state, tier, resident_country_code, preferred_language, short_bio, long_bio, created_at)
VALUES 
    ('12345678-0026-0026-0026-000000060001'::uuid, 'Test Hub User 1', 'testhub1', 'test1@0026-hubuser.example.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'en', 'Test bio', 'Test long bio', timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000060002'::uuid, 'Test Hub User 2', 'testhub2', 'test2@0026-hubuser.example.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'en', 'Test bio', 'Test long bio', timezone('UTC'::text, now()));

-- Create test hub user tokens
INSERT INTO hub_user_tokens (token, hub_user_id, token_type, token_valid_till, created_at)
VALUES 
    ('0026-test-token-1', '12345678-0026-0026-0026-000000060001'::uuid, 'HUB_USER_SESSION', timezone('UTC'::text, now()) + interval '1 day', timezone('UTC'::text, now())),
    ('0026-test-token-2', '12345678-0026-0026-0026-000000060002'::uuid, 'HUB_USER_SESSION', timezone('UTC'::text, now()) + interval '1 day', timezone('UTC'::text, now()));

-- Create test employer for org following tests
INSERT INTO emails (email_key, email_from, email_to, email_cc, email_bcc, email_subject, email_html_body, email_text_body, email_state, created_at, processed_at)
VALUES 
    ('12345678-0026-0026-0026-000000000015'::uuid, 'no-reply@vetchi.org', ARRAY['admin@0026-orgfollow.example.com'], NULL, NULL, 'Welcome to Vetchium Subject', 'Welcome HTML', 'Welcome Text', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000000016'::uuid, 'no-reply@vetchi.org', ARRAY['admin@0026-hubtest-different.example.com'], NULL, NULL, 'Welcome to Vetchium Subject', 'Welcome HTML', 'Welcome Text', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now()));

INSERT INTO employers (id, client_id_type, employer_state, company_name, onboard_admin_email, onboard_secret_token, token_valid_till, onboard_email_id, created_at)
VALUES 
    ('12345678-0026-0026-0026-000000000205'::uuid, 'DOMAIN', 'ONBOARDED', '0026-orgfollow.example.com', 'admin@0026-orgfollow.example.com', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0026-0026-0026-000000000015'::uuid, timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000000206'::uuid, 'DOMAIN', 'ONBOARDED', '0026-hubtest-different.example.com', 'admin@0026-hubtest-different.example.com', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0026-0026-0026-000000000016'::uuid, timezone('UTC'::text, now()));

INSERT INTO domains (id, domain_name, domain_state, employer_id, created_at)
VALUES 
    ('12345678-0026-0026-0026-000000003005'::uuid, '0026-orgfollow.example.com', 'VERIFIED', '12345678-0026-0026-0026-000000000205'::uuid, timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000003006'::uuid, '0026-hubtest-different.example.com', 'VERIFIED', '12345678-0026-0026-0026-000000000206'::uuid, timezone('UTC'::text, now()));

INSERT INTO employer_primary_domains (employer_id, domain_id)
VALUES 
    ('12345678-0026-0026-0026-000000000205'::uuid, '12345678-0026-0026-0026-000000003005'::uuid),
    ('12345678-0026-0026-0026-000000000206'::uuid, '12345678-0026-0026-0026-000000003006'::uuid);

INSERT INTO org_users (id, email, name, password_hash, org_user_roles, org_user_state, employer_id, created_at)
VALUES
    ('12345678-0026-0026-0026-000000040013'::uuid, 'admin@0026-orgfollow.example.com', 'Admin User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0026-0026-0026-000000000205'::uuid, timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000040014'::uuid, 'admin@0026-hubtest-different.example.com', 'Different Employer Admin', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0026-0026-0026-000000000206'::uuid, timezone('UTC'::text, now()));

COMMIT;