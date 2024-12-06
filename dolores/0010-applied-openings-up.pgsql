BEGIN;

-- Create emails for employers
INSERT INTO emails (
    email_key,
    email_from,
    email_to,
    email_subject,
    email_html_body,
    email_text_body,
    email_state,
    created_at,
    processed_at
) VALUES (
    '12345678-0010-0010-0010-000000000011'::uuid,
    'no-reply@vetchi.org',
    ARRAY['admin@applied1.example'],
    'Welcome to Vetchi',
    'Welcome HTML',
    'Welcome Text',
    'PROCESSED',
    NOW(),
    NOW()
), (
    '12345678-0010-0010-0010-000000000012'::uuid,
    'no-reply@vetchi.org',
    ARRAY['admin@applied2.example'],
    'Welcome to Vetchi',
    'Welcome HTML',
    'Welcome Text',
    'PROCESSED',
    NOW(),
    NOW()
);

-- Create employers
INSERT INTO employers (
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
    '12345678-0010-0010-0010-000000000201'::uuid,
    'DOMAIN',
    'ONBOARDED',
    'Applied1 Inc',
    'admin@applied1.example',
    'secret1',
    NOW() + INTERVAL '1 day',
    '12345678-0010-0010-0010-000000000011'::uuid,
    NOW()
), (
    '12345678-0010-0010-0010-000000000202'::uuid,
    'DOMAIN',
    'ONBOARDED',
    'Applied2 Inc',
    'admin@applied2.example',
    'secret2',
    NOW() + INTERVAL '1 day',
    '12345678-0010-0010-0010-000000000012'::uuid,
    NOW()
);

-- Create domains
INSERT INTO domains (
    id,
    domain_name,
    domain_state,
    employer_id,
    created_at
) VALUES (
    '12345678-0010-0010-0010-000000000301'::uuid,
    'applied1.example',
    'VERIFIED',
    '12345678-0010-0010-0010-000000000201'::uuid,
    NOW()
), (
    '12345678-0010-0010-0010-000000000302'::uuid,
    'applied2.example',
    'VERIFIED',
    '12345678-0010-0010-0010-000000000202'::uuid,
    NOW()
);

-- Create org_users
INSERT INTO org_users (
    id,
    email,
    name,
    password_hash,
    org_user_roles,
    org_user_state,
    employer_id,
    created_at
) VALUES (
    '12345678-0010-0010-0010-000000000401'::uuid,
    'admin@applied1.example',
    'Admin User',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    ARRAY['ADMIN']::org_user_roles[],
    'ACTIVE_ORG_USER',
    '12345678-0010-0010-0010-000000000201'::uuid,
    NOW()
), (
    '12345678-0010-0010-0010-000000000402'::uuid,
    'viewer@applied1.example',
    'Viewer User',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    ARRAY['OPENINGS_VIEWER']::org_user_roles[],
    'ACTIVE_ORG_USER',
    '12345678-0010-0010-0010-000000000201'::uuid,
    NOW()
), (
    '12345678-0010-0010-0010-000000000403'::uuid,
    'non-app@applied1.example',
    'Non App User',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    ARRAY[]::org_user_roles[],
    'ACTIVE_ORG_USER',
    '12345678-0010-0010-0010-000000000201'::uuid,
    NOW()
);

-- Create org_user_tokens
INSERT INTO org_user_tokens (
    token,
    org_user_id,
    token_type,
    token_valid_till,
    created_at
) VALUES (
    'admin-token',
    '12345678-0010-0010-0010-000000000401'::uuid,
    'EMPLOYER_SESSION',
    NOW() + INTERVAL '1 day',
    NOW()
), (
    'viewer-token',
    '12345678-0010-0010-0010-000000000402'::uuid,
    'EMPLOYER_SESSION',
    NOW() + INTERVAL '1 day',
    NOW()
), (
    'non-app-token',
    '12345678-0010-0010-0010-000000000403'::uuid,
    'EMPLOYER_SESSION',
    NOW() + INTERVAL '1 day',
    NOW()
);

-- Create hub users
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
    created_at,
    updated_at
) VALUES (
    '12345678-0010-0010-0010-000000050001'::uuid,
    'Hub User 1',
    'hub_user_1',
    'hub1@applied1.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'IND',
    'Bangalore',
    'en',
    NOW(),
    NOW()
), (
    '12345678-0010-0010-0010-000000050002'::uuid,
    'Hub User 2',
    'hub_user_2',
    'hub2@applied1.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'IND',
    'Bangalore',
    'en',
    NOW(),
    NOW()
);

-- Create hub user tokens
INSERT INTO hub_user_tokens (
    token,
    hub_user_id,
    token_type,
    token_valid_till,
    created_at
) VALUES (
    'hub1-token',
    '12345678-0010-0010-0010-000000050001'::uuid,
    'HUB_USER_SESSION',
    NOW() + INTERVAL '1 day',
    NOW()
), (
    'hub2-token',
    '12345678-0010-0010-0010-000000050002'::uuid,
    'HUB_USER_SESSION',
    NOW() + INTERVAL '1 day',
    NOW()
);

-- Create openings
INSERT INTO openings (
    employer_id,
    id,
    title,
    positions,
    jd,
    recruiter,
    hiring_manager,
    opening_type,
    yoe_min,
    yoe_max,
    min_education_level,
    state,
    created_at,
    last_updated_at
) VALUES (
    '12345678-0010-0010-0010-000000000201'::uuid,
    '2024-Mar-06-1',
    'Software Engineer',
    1,
    'Software Engineer position',
    '12345678-0010-0010-0010-000000000401'::uuid,
    '12345678-0010-0010-0010-000000000401'::uuid,
    'FULL_TIME_OPENING',
    2,
    5,
    'BACHELOR_EDUCATION',
    'ACTIVE_OPENING_STATE',
    NOW(),
    NOW()
), (
    '12345678-0010-0010-0010-000000000201'::uuid,
    '2024-Mar-06-2',
    'Senior Software Engineer',
    1,
    'Senior Software Engineer position',
    '12345678-0010-0010-0010-000000000401'::uuid,
    '12345678-0010-0010-0010-000000000401'::uuid,
    'FULL_TIME_OPENING',
    5,
    8,
    'BACHELOR_EDUCATION',
    'ACTIVE_OPENING_STATE',
    NOW(),
    NOW()
);

-- Create applications
INSERT INTO applications (
    id,
    employer_id,
    opening_id,
    cover_letter,
    original_filename,
    internal_filename,
    application_state,
    color_tag,
    hub_user_id,
    created_at
) VALUES (
    'APP-12345678-0010-0010-0010-000000000201-1',
    '12345678-0010-0010-0010-000000000201'::uuid,
    '2024-Mar-06-1',
    'Cover Letter 1',
    'resume1.pdf',
    'internal-1.pdf',
    'APPLIED',
    NULL,
    '12345678-0010-0010-0010-000000050001'::uuid,
    NOW()
), (
    'APP-12345678-0010-0010-0010-000000000201-2',
    '12345678-0010-0010-0010-000000000201'::uuid,
    '2024-Mar-06-1',
    'Cover Letter 2',
    'resume2.pdf',
    'internal-2.pdf',
    'APPLIED',
    NULL,
    '12345678-0010-0010-0010-000000050002'::uuid,
    NOW()
), (
    'APP-12345678-0010-0010-0010-000000000202-1',
    '12345678-0010-0010-0010-000000000201'::uuid,
    '2024-Mar-06-2',
    'Cover Letter 3',
    'resume3.pdf',
    'internal-3.pdf',
    'APPLIED',
    NULL,
    '12345678-0010-0010-0010-000000050001'::uuid,
    NOW()
);

COMMIT;
