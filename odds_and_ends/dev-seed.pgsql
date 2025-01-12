BEGIN;


INSERT INTO public.emails (email_key, email_from, email_to, email_cc, email_bcc, email_subject, email_html_body, email_text_body, email_state, created_at, processed_at)
    VALUES ('87654321-8336-8336-8336-000000000011'::uuid, 'no-reply@vetchi.org', ARRAY['admin@vetchidev.example'], NULL, NULL, 'Welcome to Vetchi Subject', 'Welcome to Vetchi HTML Body', 'Welcome to Vetchi Text Body', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now()));


INSERT INTO public.employers (id, client_id_type, employer_state, company_name, onboard_admin_email, onboard_secret_token, token_valid_till, onboard_email_id, created_at)
    VALUES ('87654321-8336-8336-8336-000000000201'::uuid, 'DOMAIN', 'ONBOARDED', 'vetchidev.example', 'admin@vetchidev.example', 'blah', timezone('UTC'::text, now()) + interval '1 day', '87654321-8336-8336-8336-000000000011'::uuid, timezone('UTC'::text, now()));


INSERT INTO public.domains (id, domain_name, domain_state, employer_id, created_at)
    VALUES ('87654321-8336-8336-8336-000000003001'::uuid, 'vetchidev.example', 'VERIFIED', '87654321-8336-8336-8336-000000000201'::uuid, timezone('UTC'::text, now()));

INSERT INTO public.employer_primary_domains (employer_id, domain_id)
    VALUES ('87654321-8336-8336-8336-000000000201'::uuid, '87654321-8336-8336-8336-000000003001'::uuid);

INSERT INTO public.org_users (id, name, email, password_hash, org_user_roles, org_user_state, employer_id, created_at)
    VALUES ('87654321-8336-8336-8336-000000040001'::uuid, 'admin', 'admin@vetchidev.example', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', '87654321-8336-8336-8336-000000000201'::uuid, timezone('UTC'::text, now()));

INSERT INTO public.org_cost_centers (id, cost_center_name, cost_center_state, notes, employer_id, created_at)
    VALUES 
    ('87654321-8336-8336-8336-000000050001'::uuid, 'Engineering', 'ACTIVE_CC', 'Engineering Cost Center', '87654321-8336-8336-8336-000000000201'::uuid, timezone('UTC'::text, now())),
    ('87654321-8336-8336-8336-000000050002'::uuid, 'Marketing', 'ACTIVE_CC', 'Marketing Cost Center', '87654321-8336-8336-8336-000000000201'::uuid, timezone('UTC'::text, now())),
    ('87654321-8336-8336-8336-000000050003'::uuid, 'Sales', 'ACTIVE_CC', 'Sales Cost Center', '87654321-8336-8336-8336-000000000201'::uuid, timezone('UTC'::text, now())),
    ('87654321-8336-8336-8336-000000050004'::uuid, 'HR', 'ACTIVE_CC', 'HR Cost Center', '87654321-8336-8336-8336-000000000201'::uuid, timezone('UTC'::text, now()));

INSERT INTO public.locations (id, title, country_code, postal_address, postal_code, openstreetmap_url, city_aka, location_state, employer_id, created_at)
    VALUES 
    ('87654321-8336-8336-8336-000000060001'::uuid, 'ABC Tech Park', 'USA', '123 Main St, Nob Hills, San Francisco, California, United States of America', '12345', 'https://www.openstreetmap.org/node/12345', ARRAY['San Francisco', 'SFO']::text[], 'ACTIVE_LOCATION', '87654321-8336-8336-8336-000000000201'::uuid, timezone('UTC'::text, now())),
    ('87654321-8336-8336-8336-000000060002'::uuid, 'Mormonix Tech Park', 'USA', '5, Joseph Smith Street, 1st Block, Provo, Utah, United States of America', '23456', 'https://www.openstreetmap.org/node/23456', ARRAY['Provo']::text[], 'ACTIVE_LOCATION', '87654321-8336-8336-8336-000000000201'::uuid, timezone('UTC'::text, now())),
    ('87654321-8336-8336-8336-000000060003'::uuid, 'Tidel Park', 'IND', '18, Old Mahabalipuram Road, Kanchipuram, Tamil Nadu, India', '600042', 'https://www.openstreetmap.org/node/34567', ARRAY['Kanchipuram', 'Madras', 'Chennai']::text[], 'ACTIVE_LOCATION', '87654321-8336-8336-8336-000000000201'::uuid, timezone('UTC'::text, now())),
    ('87654321-8336-8336-8336-000000060004'::uuid, 'Manjunatha Tech Park', 'IND', '1, Hosur Road, Bengalooru, Karnataka, India', '560029', 'https://www.openstreetmap.org/node/560029', ARRAY['Bengalooru', 'Bengaluru', 'Bangalore']::text[], 'ACTIVE_LOCATION', '87654321-8336-8336-8336-000000000201'::uuid, timezone('UTC'::text, now()))
    ;


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
