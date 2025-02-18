BEGIN;

-- Create welcome email first
INSERT INTO emails (email_key, email_from, email_to, email_cc, email_bcc, email_subject, email_html_body, email_text_body, email_state, created_at, processed_at)
VALUES ('12345678-0015-0015-0015-000000000011'::uuid, 'no-reply@vetchi.org', ARRAY['tags.test@openingtags.example'], NULL, NULL, 'Welcome to Vetchi Subject', 'Welcome to Vetchi HTML Body', 'Welcome to Vetchi Text Body', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now()));

-- Create employer with proper onboarding fields
INSERT INTO employers (id, client_id_type, employer_state, company_name, onboard_admin_email, onboard_secret_token, token_valid_till, onboard_email_id, created_at)
VALUES ('12345678-0015-0015-0015-000000000201'::uuid, 'DOMAIN', 'ONBOARDED', 'openingtags.example', 'tags.test@openingtags.example', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0015-0015-0015-000000000011'::uuid, timezone('UTC'::text, now()));

-- Create domain
INSERT INTO domains (id, domain_name, domain_state, employer_id, created_at)
VALUES ('12345678-0015-0015-0015-000000003001'::uuid, 'openingtags.example', 'VERIFIED', '12345678-0015-0015-0015-000000000201'::uuid, timezone('UTC'::text, now()));

-- Set primary domain
INSERT INTO employer_primary_domains (employer_id, domain_id)
VALUES ('12345678-0015-0015-0015-000000000201'::uuid, '12345678-0015-0015-0015-000000003001'::uuid);

-- Create org user with admin role
INSERT INTO org_users (id, email, name, password_hash, org_user_roles, org_user_state, employer_id, created_at)
VALUES ('12345678-0015-0015-0015-000000040001'::uuid, 'tags.test@openingtags.example', 'Tags Test Admin', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0015-0015-0015-000000000201'::uuid, timezone('UTC'::text, now()));

-- Create cost center
INSERT INTO org_cost_centers (id, cost_center_name, cost_center_state, notes, employer_id, created_at)
VALUES ('12345678-0015-0015-0015-000000050001'::uuid, 'Engineering', 'ACTIVE_CC', 'Engineering department', '12345678-0015-0015-0015-000000000201'::uuid, timezone('UTC'::text, now()));

-- Create location
INSERT INTO locations (id, title, country_code, postal_address, postal_code, city_aka, location_state, employer_id, created_at)
VALUES ('12345678-0015-0015-0015-000000060001'::uuid, 'Main Office', 'IND', '123 Main St', '600001', ARRAY['Chennai', 'Madras'], 'ACTIVE_LOCATION', '12345678-0015-0015-0015-000000000201'::uuid, timezone('UTC'::text, now()));

-- Create hub user
INSERT INTO hub_users (
    id,
    full_name,
    handle,
    email,
    password_hash,
    state,
    resident_country_code,
    resident_city,
    preferred_language,
    short_bio,
    long_bio,
    created_at,
    updated_at
) VALUES (
    '12345678-0015-0015-0015-000000050002'::uuid,
    'Tags Test Hub User',
    'tags_test',
    'tags.test.hub@openingtags.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'IND',
    'Chennai',
    'en',
    'Tags Test Hub User is methodical',
    'Tags Test Hub User was born in India and finished education at PSG Tech and has 5 years as experience.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
);

-- Create pre-existing tags for testing
INSERT INTO opening_tags (id, name, created_at)
VALUES 
    ('12345678-0015-0015-0015-000000070001'::uuid, 'Go', timezone('UTC'::text, now())),
    ('12345678-0015-0015-0015-000000070002'::uuid, 'Java', timezone('UTC'::text, now())),
    ('12345678-0015-0015-0015-000000070003'::uuid, 'Python', timezone('UTC'::text, now())),
    ('12345678-0015-0015-0015-000000070004'::uuid, 'PostgreSQL', timezone('UTC'::text, now())),
    ('12345678-0015-0015-0015-000000070005'::uuid, 'React', timezone('UTC'::text, now()));

COMMIT;
