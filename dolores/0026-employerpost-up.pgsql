BEGIN;

-- Create employer for testing
INSERT INTO emails (email_key, email_from, email_to, email_cc, email_bcc, email_subject, email_html_body, email_text_body, email_state, created_at, processed_at)
VALUES ('12345678-0026-0026-0026-000000000011'::uuid, 'no-reply@vetchi.org', ARRAY['admin@0026-employerposts.example.com'], NULL, NULL, 'Welcome to Vetchium Subject', 'Welcome HTML', 'Welcome Text', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now()));

INSERT INTO employers (id, client_id_type, employer_state, company_name, onboard_admin_email, onboard_secret_token, token_valid_till, onboard_email_id, created_at)
VALUES ('12345678-0026-0026-0026-000000000201'::uuid, 'DOMAIN', 'ONBOARDED', '0026-employerposts.example.com', 'admin@0026-employerposts.example.com', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0026-0026-0026-000000000011'::uuid, timezone('UTC'::text, now()));

INSERT INTO domains (id, domain_name, domain_state, employer_id, created_at)
VALUES ('12345678-0026-0026-0026-000000003001'::uuid, '0026-employerposts.example.com', 'VERIFIED', '12345678-0026-0026-0026-000000000201'::uuid, timezone('UTC'::text, now()));

-- Set primary domain
INSERT INTO employer_primary_domains (employer_id, domain_id)
VALUES ('12345678-0026-0026-0026-000000000201'::uuid, '12345678-0026-0026-0026-000000003001'::uuid);

-- Create org users with different roles
INSERT INTO org_users (id, email, name, password_hash, org_user_roles, org_user_state, employer_id, created_at)
VALUES
    ('12345678-0026-0026-0026-000000040001'::uuid, 'admin@0026-employerposts.example.com', 'Admin User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0026-0026-0026-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000040002'::uuid, 'marketing@0026-employerposts.example.com', 'Marketing User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['EMPLOYER_POSTS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0026-0026-0026-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0026-0026-0026-000000040003'::uuid, 'regular@0026-employerposts.example.com', 'Regular User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['APPLICATIONS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0026-0026-0026-000000000201'::uuid, timezone('UTC'::text, now()));

-- Add tags for testing
INSERT INTO tags (id, name)
VALUES
    ('12345678-0026-0026-0026-000000050001'::uuid, 'engineering'),
    ('12345678-0026-0026-0026-000000050002'::uuid, 'marketing'),
    ('12345678-0026-0026-0026-000000050003'::uuid, 'golang'),
    ('12345678-0026-0026-0026-000000050004'::uuid, 'react'),
    ('12345678-0026-0026-0026-000000050005'::uuid, 'testing')
ON CONFLICT (name) DO NOTHING;

-- Create some employer posts for testing
INSERT INTO employer_posts (id, content, employer_id, created_at, updated_at)
VALUES
    ('12345678-0026-0026-0026-000000060001', 'First employer post for testing', '12345678-0026-0026-0026-000000000201'::uuid, timezone('UTC'::text, now() - interval '3 days'), timezone('UTC'::text, now() - interval '3 days')),
    ('12345678-0026-0026-0026-000000060002', 'Second employer post with tags', '12345678-0026-0026-0026-000000000201'::uuid, timezone('UTC'::text, now() - interval '2 days'), timezone('UTC'::text, now() - interval '2 days')),
    ('12345678-0026-0026-0026-000000060003', 'Third employer post for pagination testing', '12345678-0026-0026-0026-000000000201'::uuid, timezone('UTC'::text, now() - interval '1 day'), timezone('UTC'::text, now() - interval '1 day')),
    ('12345678-0026-0026-0026-000000060004', 'Fourth employer post, most recent', '12345678-0026-0026-0026-000000000201'::uuid, timezone('UTC'::text, now()), timezone('UTC'::text, now()));

-- Add tags to posts
INSERT INTO employer_post_tags (employer_post_id, tag_id)
VALUES
    ('12345678-0026-0026-0026-000000060002', '12345678-0026-0026-0026-000000050001'::uuid),  -- engineering
    ('12345678-0026-0026-0026-000000060002', '12345678-0026-0026-0026-000000050003'::uuid),  -- golang
    ('12345678-0026-0026-0026-000000060003', '12345678-0026-0026-0026-000000050002'::uuid),  -- marketing
    ('12345678-0026-0026-0026-000000060004', '12345678-0026-0026-0026-000000050003'::uuid),  -- golang
    ('12345678-0026-0026-0026-000000060004', '12345678-0026-0026-0026-000000050004'::uuid),  -- react
    ('12345678-0026-0026-0026-000000060004', '12345678-0026-0026-0026-000000050005'::uuid);  -- testing

COMMIT;