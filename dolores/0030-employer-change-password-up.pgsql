BEGIN;

-- Create an email record first (required for employers table)
INSERT INTO public.emails (
    email_key,
    email_from,
    email_to,
    email_cc,
    email_bcc,
    email_subject,
    email_html_body,
    email_text_body,
    email_state,
    created_at,
    processed_at
) VALUES (
    '12345678-0030-0030-0030-000000000011'::uuid,
    'no-reply@vetchi.org',
    ARRAY['admin@0030-changepassword.example'],
    NULL,
    NULL,
    'Welcome to Vetchium Subject',
    'Welcome to Vetchium HTML Body',
    'Welcome to Vetchium Text Body',
    'PROCESSED',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
);

-- Create an employer for change password tests
INSERT INTO public.employers (
    id,
    client_id_type,
    employer_state,
    company_name,
    onboard_admin_email,
    onboard_secret_token,
    token_valid_till,
    onboard_email_id,
    created_at
) VALUES (
    '12345678-0030-0030-0030-000000000001'::uuid,
    'DOMAIN',
    'ONBOARDED',
    'Change Password Test Company',
    'admin@0030-changepassword.example',
    'blah',
    timezone('UTC'::text, now()) + interval '1 day',
    '12345678-0030-0030-0030-000000000011'::uuid,
    timezone('UTC'::text, now())
);

-- Create domain for the employer
INSERT INTO public.domains (
    id,
    domain_name,
    domain_state,
    employer_id,
    created_at
) VALUES (
    '12345678-0030-0030-0030-000000003001'::uuid,
    '0030-changepassword.example',
    'VERIFIED',
    '12345678-0030-0030-0030-000000000001'::uuid,
    timezone('UTC'::text, now())
);

-- Set primary domain
INSERT INTO public.employer_primary_domains (
    employer_id,
    domain_id
) VALUES (
    '12345678-0030-0030-0030-000000000001'::uuid,
    '12345678-0030-0030-0030-000000003001'::uuid
);

-- Password hash for "NewPassword123$"
INSERT INTO public.org_users (
    id,
    employer_id,
    name,
    email,
    password_hash,
    org_user_roles,
    org_user_state,
    created_at
) VALUES 
    (
        '12345678-0030-0030-0030-000000000001'::uuid,
        '12345678-0030-0030-0030-000000000001'::uuid,
        'Change Password Test User 1',
        'change1@0030-changepassword.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        ARRAY['ADMIN']::org_user_roles[],
        'ACTIVE_ORG_USER',
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0030-0030-0030-000000000002'::uuid,
        '12345678-0030-0030-0030-000000000001'::uuid,
        'Change Password Test User 2',
        'change2@0030-changepassword.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        ARRAY['ADMIN']::org_user_roles[],
        'ACTIVE_ORG_USER',
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0030-0030-0030-000000000003'::uuid,
        '12345678-0030-0030-0030-000000000001'::uuid,
        'Change Password Test User 3',
        'change3@0030-changepassword.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        ARRAY['ADMIN']::org_user_roles[],
        'ACTIVE_ORG_USER',
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0030-0030-0030-000000000004'::uuid,
        '12345678-0030-0030-0030-000000000001'::uuid,
        'Change Password Test User 4',
        'change4@0030-changepassword.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        ARRAY['ADMIN']::org_user_roles[],
        'ACTIVE_ORG_USER',
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0030-0030-0030-000000000005'::uuid,
        '12345678-0030-0030-0030-000000000001'::uuid,
        'Change Password Test User 5',
        'change5@0030-changepassword.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        ARRAY['ADMIN']::org_user_roles[],
        'ACTIVE_ORG_USER',
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0030-0030-0030-000000000006'::uuid,
        '12345678-0030-0030-0030-000000000001'::uuid,
        'Change Password Test User 6',
        'change6@0030-changepassword.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        ARRAY['ADMIN']::org_user_roles[],
        'ACTIVE_ORG_USER',
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0030-0030-0030-000000000007'::uuid,
        '12345678-0030-0030-0030-000000000001'::uuid,
        'Change Password Test User 7',
        'change7@0030-changepassword.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        ARRAY['ADMIN']::org_user_roles[],
        'ACTIVE_ORG_USER',
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0030-0030-0030-000000000008'::uuid,
        '12345678-0030-0030-0030-000000000001'::uuid,
        'Change Password Test User 8',
        'change8@0030-changepassword.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        ARRAY['ADMIN']::org_user_roles[],
        'ACTIVE_ORG_USER',
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0030-0030-0030-000000000009'::uuid,
        '12345678-0030-0030-0030-000000000001'::uuid,
        'Session Test User',
        'session-test@0030-changepassword.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        ARRAY['ADMIN']::org_user_roles[],
        'ACTIVE_ORG_USER',
        timezone('UTC'::text, now())
    );

COMMIT; 