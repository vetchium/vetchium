BEGIN;


INSERT INTO public.emails (email_key, email_from, email_to, email_cc, email_bcc, email_subject, email_html_body, email_text_body, email_state, created_at, processed_at)
    VALUES ('87654321-8336-8336-8336-000000000011'::uuid, 'no-reply@vetchi.org', ARRAY['admin@vetchidev.example'], NULL, NULL, 'Welcome to Vetchi Subject', 'Welcome to Vetchi HTML Body', 'Welcome to Vetchi Text Body', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now()));


INSERT INTO public.employers (id, client_id_type, employer_state, company_name, onboard_admin_email, onboard_secret_token, token_valid_till, onboard_email_id, created_at)
    VALUES ('87654321-8336-8336-8336-000000000201'::uuid, 'DOMAIN', 'ONBOARDED', 'vetchidev.example', 'admin@vetchidev.example', 'blah', timezone('UTC'::text, now()) + interval '1 day', '87654321-8336-8336-8336-000000000011'::uuid, timezone('UTC'::text, now()));


INSERT INTO public.domains (id, domain_name, domain_state, employer_id, created_at)
    VALUES ('87654321-8336-8336-8336-000000003001'::uuid, 'vetchidev.example', 'VERIFIED', '87654321-8336-8336-8336-000000000201'::uuid, timezone('UTC'::text, now()));


INSERT INTO public.org_users (id, name, email, password_hash, org_user_roles, org_user_state, employer_id, created_at)
    VALUES ('87654321-8336-8336-8336-000000040001'::uuid, 'admin', 'admin@vetchidev.example', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', '87654321-8336-8336-8336-000000000201'::uuid, timezone('UTC'::text, now()));


INSERT INTO public.hub_users (
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
)
VALUES 
    (
        '56781234-5678-5678-5678-000000050001'::uuid,
        'User 1',
        'user1',
        'user1@hub.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        'ACTIVE_HUB_USER',
        'USA',
        'New York',
        'en',
        timezone('UTC'::text, now()),
        timezone('UTC'::text, now())
    ),
    (
        '56781234-5678-5678-5678-000000050002'::uuid,
        'User 2',
        'user2',
        'user2@hub.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        'ACTIVE_HUB_USER',
        'India',
        'Chennai',
        'en',
        timezone('UTC'::text, now()),
        timezone('UTC'::text, now())
    ),
    (
        '56781234-5678-5678-5678-000000050003'::uuid,
        'User 3',
        'user3',
        'user3@hub.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        'ACTIVE_HUB_USER',
        'France',
        'Some hippie beach',
        'fr',
        timezone('UTC'::text, now()),
        timezone('UTC'::text, now())
    ),
    (
        '56781234-5678-5678-5678-000000050004'::uuid,
        'தமிழ்',
        'sometamilaccount',
        'தமிழ்@hub.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        'ACTIVE_HUB_USER',
        'India',
        'Chennai',
        'ta',
        timezone('UTC'::text, now()),
        timezone('UTC'::text, now())
    );

COMMIT;
