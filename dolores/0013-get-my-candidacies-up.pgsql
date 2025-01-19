BEGIN;

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
) VALUES
    (
        '12345678-0013-0013-0013-000000050001'::uuid,
        'Hub User One',
        'hubuser1',
        'hubuser1@my-candidacies.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        'ACTIVE_HUB_USER',
        'IND',
        'Bangalore',
        'en',
        timezone('UTC'::text, now()),
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0013-0013-0013-000000050002'::uuid,
        'Hub User Two',
        'hubuser2',
        'hubuser2@my-candidacies.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        'ACTIVE_HUB_USER',
        'USA',
        'New York',
        'en',
        timezone('UTC'::text, now()),
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0013-0013-0013-000000050003'::uuid,
        'Hub User Three',
        'hubuser3',
        'hubuser3@my-candidacies.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        'ACTIVE_HUB_USER',
        'GBR',
        'London',
        'en',
        timezone('UTC'::text, now()),
        timezone('UTC'::text, now())
    );

-- Create employers (3 companies)
INSERT INTO emails (email_key, email_from, email_to, email_cc, email_bcc, email_subject, email_html_body, email_text_body, email_state, created_at, processed_at)
VALUES
    ('12345678-0013-0013-0013-000000000011'::uuid, 'no-reply@vetchi.org', ARRAY['admin@my-candidacies-1.example'], NULL, NULL, 'Welcome to Vetchi', 'Welcome HTML', 'Welcome Text', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now())),
    ('12345678-0013-0013-0013-000000000012'::uuid, 'no-reply@vetchi.org', ARRAY['admin@my-candidacies-2.example'], NULL, NULL, 'Welcome to Vetchi', 'Welcome HTML', 'Welcome Text', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now())),
    ('12345678-0013-0013-0013-000000000013'::uuid, 'no-reply@vetchi.org', ARRAY['admin@my-candidacies-3.example'], NULL, NULL, 'Welcome to Vetchi', 'Welcome HTML', 'Welcome Text', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now()));

INSERT INTO employers (id, client_id_type, employer_state, company_name, onboard_admin_email, onboard_secret_token, token_valid_till, onboard_email_id, created_at)
VALUES
    ('12345678-0013-0013-0013-000000000201'::uuid, 'DOMAIN', 'ONBOARDED', 'my-candidacies-1.example', 'admin@my-candidacies-1.example', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0013-0013-0013-000000000011'::uuid, timezone('UTC'::text, now())),
    ('12345678-0013-0013-0013-000000000202'::uuid, 'DOMAIN', 'ONBOARDED', 'my-candidacies-2.example', 'admin@my-candidacies-2.example', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0013-0013-0013-000000000012'::uuid, timezone('UTC'::text, now())),
    ('12345678-0013-0013-0013-000000000203'::uuid, 'DOMAIN', 'ONBOARDED', 'my-candidacies-3.example', 'admin@my-candidacies-3.example', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0013-0013-0013-000000000013'::uuid, timezone('UTC'::text, now()));

-- Create domains
INSERT INTO domains (id, domain_name, domain_state, employer_id, created_at)
VALUES
    ('12345678-0013-0013-0013-000000003001'::uuid, 'my-candidacies-1.example', 'VERIFIED', '12345678-0013-0013-0013-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0013-0013-0013-000000003002'::uuid, 'my-candidacies-2.example', 'VERIFIED', '12345678-0013-0013-0013-000000000202'::uuid, timezone('UTC'::text, now())),
    ('12345678-0013-0013-0013-000000003003'::uuid, 'my-candidacies-3.example', 'VERIFIED', '12345678-0013-0013-0013-000000000203'::uuid, timezone('UTC'::text, now()));

-- After creating domains, add:
INSERT INTO employer_primary_domains (employer_id, domain_id)
SELECT employer_id, id 
FROM domains 
WHERE domain_name IN (
    'my-candidacies-1.example',
    'my-candidacies-2.example',
    'my-candidacies-3.example'
);

-- Create org users
INSERT INTO org_users (id, email, name, password_hash, org_user_roles, org_user_state, employer_id, created_at)
VALUES
    -- Company 1 users
    ('12345678-0013-0013-0013-000000040001'::uuid, 'admin@my-candidacies-1.example', 'Admin User 1', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0013-0013-0013-000000000201'::uuid, timezone('UTC'::text, now())),
    -- Company 2 users
    ('12345678-0013-0013-0013-000000040002'::uuid, 'admin@my-candidacies-2.example', 'Admin User 2', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0013-0013-0013-000000000202'::uuid, timezone('UTC'::text, now())),
    -- Company 3 users
    ('12345678-0013-0013-0013-000000040003'::uuid, 'admin@my-candidacies-3.example', 'Admin User 3', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0013-0013-0013-000000000203'::uuid, timezone('UTC'::text, now()));

-- Create cost centers
INSERT INTO org_cost_centers (id, cost_center_name, cost_center_state, notes, employer_id, created_at)
VALUES
    ('12345678-0013-0013-0013-000000050001'::uuid, 'Engineering', 'ACTIVE_CC', 'Engineering department', '12345678-0013-0013-0013-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0013-0013-0013-000000050002'::uuid, 'Engineering', 'ACTIVE_CC', 'Engineering department', '12345678-0013-0013-0013-000000000202'::uuid, timezone('UTC'::text, now())),
    ('12345678-0013-0013-0013-000000050003'::uuid, 'Engineering', 'ACTIVE_CC', 'Engineering department', '12345678-0013-0013-0013-000000000203'::uuid, timezone('UTC'::text, now()));

-- Create openings (5 per company)
INSERT INTO openings (
    employer_id,
    id,
    title,
    positions,
    jd,
    recruiter,
    hiring_manager,
    cost_center_id,
    opening_type,
    yoe_min,
    yoe_max,
    min_education_level,
    salary_min,
    salary_max,
    salary_currency,
    remote_country_codes,
    remote_timezones,
    state,
    created_at,
    last_updated_at
)
SELECT
    employer_id,
    '2024-Mar-' || employer_num || '-' || LPAD(opening_num::text, 3, '0'),
    'Software Engineer ' || opening_num,
    2,
    'Job Description ' || opening_num,
    org_user_id,
    org_user_id,
    cost_center_id,
    'FULL_TIME_OPENING',
    2,
    5,
    'BACHELOR_EDUCATION',
    50000,
    100000,
    'USD',
    ARRAY['IND', 'USA'],
    ARRAY['IST Indian Standard Time GMT+0530'],
    'ACTIVE_OPENING_STATE',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
FROM (
    SELECT 
        e.id as employer_id,
        RIGHT(e.id::text, 1) as employer_num,
        generate_series as opening_num,
        (SELECT id FROM org_users WHERE employer_id = e.id LIMIT 1) as org_user_id,
        (SELECT id FROM org_cost_centers WHERE employer_id = e.id LIMIT 1) as cost_center_id
    FROM employers e
    CROSS JOIN generate_series(1, 5)
    WHERE e.id IN (
        '12345678-0013-0013-0013-000000000201'::uuid,
        '12345678-0013-0013-0013-000000000202'::uuid,
        '12345678-0013-0013-0013-000000000203'::uuid
    )
) t;

-- Create applications and candidacies for hubuser1
WITH inserted_applications AS (
    INSERT INTO applications (
        id,
        employer_id,
        opening_id,
        hub_user_id,
        application_state,
        cover_letter,
        resume_sha,
        created_at
    )
    SELECT
        '12345678-0013-0013-' || o.id,
        o.employer_id,
        o.id,
        '12345678-0013-0013-0013-000000050001'::uuid,
        'APPLIED'::application_states,
        'Cover Letter for ' || o.id,
        'sha-sha-sha',
        timezone('UTC'::text, now())
    FROM openings o
    RETURNING id, employer_id, opening_id
)
INSERT INTO candidacies (
    id,
    application_id,
    employer_id,
    opening_id,
    candidacy_state,
    created_by,
    created_at
)
SELECT
    '12345678-0013-0013-' || opening_id || '-candidacy',
    id,
    employer_id,
    opening_id,
    CASE (row_number() OVER (ORDER BY opening_id)) % 3
        WHEN 0 THEN 'INTERVIEWING'::candidacy_states
        WHEN 1 THEN 'OFFERED'::candidacy_states
        WHEN 2 THEN 'CANDIDATE_UNSUITABLE'::candidacy_states
    END,
    (SELECT id FROM org_users WHERE employer_id = a.employer_id LIMIT 1),
    timezone('UTC'::text, now())
FROM inserted_applications a;

-- Create applications and candidacies for hubuser2 (fewer applications)
WITH inserted_applications AS (
    INSERT INTO applications (
        id,
        employer_id,
        opening_id,
        hub_user_id,
        application_state,
        cover_letter,
        resume_sha,
        created_at
    )
    SELECT
        '12345678-0013-0013-hubuser2-' || o.id,
        o.employer_id,
        o.id,
        '12345678-0013-0013-0013-000000050002'::uuid,
        'APPLIED'::application_states,
        'Cover Letter for ' || o.id,
        'sha-sha-sha',
        timezone('UTC'::text, now())
    FROM openings o
    WHERE o.employer_id = '12345678-0013-0013-0013-000000000201'::uuid
    RETURNING id, employer_id, opening_id
)
INSERT INTO candidacies (
    id,
    application_id,
    employer_id,
    opening_id,
    candidacy_state,
    created_by,
    created_at
)
SELECT
    '12345678-0013-0013-hubuser2-' || opening_id || '-candidacy',
    id,
    employer_id,
    opening_id,
    'INTERVIEWING'::candidacy_states,
    (SELECT id FROM org_users WHERE employer_id = a.employer_id LIMIT 1),
    timezone('UTC'::text, now())
FROM inserted_applications a;

COMMIT;
