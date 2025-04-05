BEGIN;
--- email table primary key uuids should end in 2 digits, 11, 12, 13, etc
--- employer table primary key uuids should end in 3 digits, 201, 202, 203, etc
--- domain table primary key uuids should end in 4 digits, 3001, 3002, 3003, etc
--- org_users table primary key uuids should end in 5 digits, 40001, 40002, 40003, etc
INSERT INTO public.emails (email_key, email_from, email_to, email_cc, email_bcc, email_subject, email_html_body, email_text_body, email_state, created_at, processed_at)
    VALUES ('12345678-0001-0001-0001-000000000011'::uuid, 'no-reply@vetchi.org', ARRAY['admin@domain-onboarded.example'], NULL, NULL, 'Welcome to Vetchium Subject', 'Welcome to Vetchium HTML Body', 'Welcome to Vetchium Text Body', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now()));
INSERT INTO public.employers (id, client_id_type, employer_state, company_name, onboard_admin_email, onboard_secret_token, token_valid_till, onboard_email_id, created_at)
    VALUES ('12345678-0001-0001-0001-000000000201'::uuid, 'DOMAIN', 'ONBOARDED', 'domain-onboarded.example', 'admin@domain-onboarded.example', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0001-0001-0001-000000000011'::uuid, timezone('UTC'::text, now()));
INSERT INTO public.domains (id, domain_name, domain_state, employer_id, created_at)
    VALUES ('12345678-0001-0001-0001-000000003001'::uuid, 'domain-onboarded.example', 'VERIFIED', '12345678-0001-0001-0001-000000000201'::uuid, timezone('UTC'::text, now()));

-- Set primary domain
INSERT INTO public.employer_primary_domains (employer_id, domain_id)
    VALUES ('12345678-0001-0001-0001-000000000201'::uuid, '12345678-0001-0001-0001-000000003001'::uuid);

INSERT INTO public.org_users (id, name, email, password_hash, org_user_roles, org_user_state, employer_id, created_at)
    VALUES ('12345678-0001-0001-0001-000000040001'::uuid, 'admin', 'admin@domain-onboarded.example', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0001-0001-0001-000000000201'::uuid, timezone('UTC'::text, now()));
COMMIT;
