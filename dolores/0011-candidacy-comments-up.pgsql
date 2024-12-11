BEGIN;

-- Create employer for testing
INSERT INTO emails (email_key, email_from, email_to, email_cc, email_bcc, email_subject, email_html_body, email_text_body, email_state, created_at, processed_at)
VALUES ('12345678-0011-0011-0011-000000000011'::uuid, 'no-reply@vetchi.org', ARRAY['admin@candidacy-comments.example'], NULL, NULL, 'Welcome to Vetchi Subject', 'Welcome HTML', 'Welcome Text', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now()));

INSERT INTO employers (id, client_id_type, employer_state, company_name, onboard_admin_email, onboard_secret_token, token_valid_till, onboard_email_id, created_at)
VALUES ('12345678-0011-0011-0011-000000000201'::uuid, 'DOMAIN', 'ONBOARDED', 'candidacy-comments.example', 'admin@candidacy-comments.example', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0011-0011-0011-000000000011'::uuid, timezone('UTC'::text, now()));

INSERT INTO domains (id, domain_name, domain_state, employer_id, created_at)
VALUES ('12345678-0011-0011-0011-000000003001'::uuid, 'candidacy-comments.example', 'VERIFIED', '12345678-0011-0011-0011-000000000201'::uuid, timezone('UTC'::text, now()));

-- Create org users with different roles
INSERT INTO org_users (id, email, name, password_hash, org_user_roles, org_user_state, employer_id, created_at)
VALUES 
    ('12345678-0011-0011-0011-000000040001'::uuid, 'admin@candidacy-comments.example', 'Admin User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0011-0011-0011-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0011-0011-0011-000000040002'::uuid, 'crud@candidacy-comments.example', 'CRUD User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['APPLICATIONS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0011-0011-0011-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0011-0011-0011-000000040003'::uuid, 'viewer@candidacy-comments.example', 'Viewer User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['APPLICATIONS_VIEWER']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0011-0011-0011-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0011-0011-0011-000000040004'::uuid, 'disabled@candidacy-comments.example', 'Disabled User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['APPLICATIONS_CRUD']::org_user_roles[], 'DISABLED_ORG_USER', '12345678-0011-0011-0011-000000000201'::uuid, timezone('UTC'::text, now()));

-- Create hub users for testing
INSERT INTO hub_users (id, full_name, handle, email, password_hash, state, resident_country_code, resident_city, preferred_language, created_at, updated_at)
VALUES
    ('12345678-0011-0011-0011-000000050001'::uuid, 'Active Hub User', 'active_hub_user', '0011-active@hub.example', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'USA', 'New York', 'en', timezone('UTC'::text, now()), timezone('UTC'::text, now())),
    ('12345678-0011-0011-0011-000000050002'::uuid, 'Disabled Hub User', 'disabled_hub_user', '0011-disabled@hub.example', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'DISABLED_HUB_USER', 'USA', 'New York', 'en', timezone('UTC'::text, now()), timezone('UTC'::text, now())),
    ('12345678-0011-0011-0011-000000050003'::uuid, 'Deleted Hub User', 'deleted_hub_user', '0011-deleted@hub.example', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'DELETED_HUB_USER', 'USA', 'New York', 'en', timezone('UTC'::text, now()), timezone('UTC'::text, now()));

-- Create test candidacies
INSERT INTO candidacies (id, employer_id, hub_user_id, created_at)
VALUES
    ('12345678-0011-0011-0011-000000060001'::uuid, '12345678-0011-0011-0011-000000000201'::uuid, '12345678-0011-0011-0011-000000050001'::uuid, timezone('UTC'::text, now())),
    ('12345678-0011-0011-0011-000000060002'::uuid, '12345678-0011-0011-0011-000000000201'::uuid, '12345678-0011-0011-0011-000000050002'::uuid, timezone('UTC'::text, now())),
    ('12345678-0011-0011-0011-000000060003'::uuid, '12345678-0011-0011-0011-000000000201'::uuid, '12345678-0011-0011-0011-000000050003'::uuid, timezone('UTC'::text, now()));

-- Create some initial comments
INSERT INTO candidacy_comments (id, candidacy_id, commenter_type, commenter_id, comment_text, created_at)
VALUES
    ('12345678-0011-0011-0011-000000070001'::uuid, '12345678-0011-0011-0011-000000060001'::uuid, 'EMPLOYER', '12345678-0011-0011-0011-000000040001'::uuid, 'Initial employer comment', timezone('UTC'::text, now())),
    ('12345678-0011-0011-0011-000000070002'::uuid, '12345678-0011-0011-0011-000000060001'::uuid, 'HUB_USER', '12345678-0011-0011-0011-000000050001'::uuid, 'Initial hub user comment', timezone('UTC'::text, now()));

COMMIT;
