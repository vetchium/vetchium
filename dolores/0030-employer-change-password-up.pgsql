BEGIN;

-- Create an employer for change password tests
INSERT INTO public.employers (
    id,
    name,
    domain,
    company_type,
    industry,
    state,
    employee_count_range,
    created_at
) VALUES (
    '12345678-0030-0030-0030-000000000001'::uuid,
    'Change Password Test Company',
    'changepassword.example',
    'PRIVATE_LIMITED',
    'TECHNOLOGY',
    'ACTIVE_EMPLOYER',
    'RANGE_11_50',
    timezone('UTC'::text, now())
);

-- Password hash for "CurrentPassword123$"
INSERT INTO public.org_users (
    id,
    employer_id,
    name,
    email,
    password_hash,
    roles,
    state,
    created_at
) VALUES 
    (
        '12345678-0030-0030-0030-000000000001'::uuid,
        '12345678-0030-0030-0030-000000000001'::uuid,
        'Change Password Test User 1',
        'change1@changepassword.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        ARRAY['ADMIN'],
        'ACTIVE_ORG_USER',
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0030-0030-0030-000000000002'::uuid,
        '12345678-0030-0030-0030-000000000001'::uuid,
        'Change Password Test User 2',
        'change2@changepassword.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        ARRAY['ADMIN'],
        'ACTIVE_ORG_USER',
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0030-0030-0030-000000000003'::uuid,
        '12345678-0030-0030-0030-000000000001'::uuid,
        'Change Password Test User 3',
        'change3@changepassword.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        ARRAY['ADMIN'],
        'ACTIVE_ORG_USER',
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0030-0030-0030-000000000004'::uuid,
        '12345678-0030-0030-0030-000000000001'::uuid,
        'Change Password Test User 4',
        'change4@changepassword.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        ARRAY['ADMIN'],
        'ACTIVE_ORG_USER',
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0030-0030-0030-000000000005'::uuid,
        '12345678-0030-0030-0030-000000000001'::uuid,
        'Change Password Test User 5',
        'change5@changepassword.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        ARRAY['ADMIN'],
        'ACTIVE_ORG_USER',
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0030-0030-0030-000000000006'::uuid,
        '12345678-0030-0030-0030-000000000001'::uuid,
        'Change Password Test User 6',
        'change6@changepassword.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        ARRAY['ADMIN'],
        'ACTIVE_ORG_USER',
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0030-0030-0030-000000000007'::uuid,
        '12345678-0030-0030-0030-000000000001'::uuid,
        'Change Password Test User 7',
        'change7@changepassword.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        ARRAY['ADMIN'],
        'ACTIVE_ORG_USER',
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0030-0030-0030-000000000008'::uuid,
        '12345678-0030-0030-0030-000000000001'::uuid,
        'Change Password Test User 8',
        'change8@changepassword.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        ARRAY['ADMIN'],
        'ACTIVE_ORG_USER',
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0030-0030-0030-000000000009'::uuid,
        '12345678-0030-0030-0030-000000000001'::uuid,
        'Session Test User',
        'session-test@changepassword.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        ARRAY['ADMIN'],
        'ACTIVE_ORG_USER',
        timezone('UTC'::text, now())
    );

COMMIT; 