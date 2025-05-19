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

COMMIT;