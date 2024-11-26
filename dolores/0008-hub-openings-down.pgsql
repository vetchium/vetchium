BEGIN;

-- Delete opening locations
DELETE FROM opening_locations
WHERE employer_id IN (
    SELECT id FROM employers
    WHERE onboard_admin_email LIKE '%@hubopening%.example'
);

-- Delete openings
DELETE FROM openings
WHERE employer_id IN (
    SELECT id FROM employers
    WHERE onboard_admin_email LIKE '%@hubopening%.example'
);

-- Delete locations
DELETE FROM locations
WHERE employer_id IN (
    SELECT id FROM employers
    WHERE onboard_admin_email LIKE '%@hubopening%.example'
);

-- Delete org_cost_centers
DELETE FROM org_cost_centers
WHERE employer_id IN (
    SELECT id FROM employers
    WHERE onboard_admin_email LIKE '%@hubopening%.example'
);

-- Delete org_user_tokens
DELETE FROM org_user_tokens
WHERE org_user_id IN (
    SELECT id FROM org_users
    WHERE employer_id IN (
        SELECT id FROM employers
        WHERE onboard_admin_email LIKE '%@hubopening%.example'
    )
);

-- Delete org_users
DELETE FROM org_users
WHERE employer_id IN (
    SELECT id FROM employers
    WHERE onboard_admin_email LIKE '%@hubopening%.example'
);

-- Delete domains
DELETE FROM domains
WHERE employer_id IN (
    SELECT id FROM employers
    WHERE onboard_admin_email LIKE '%@hubopening%.example'
);

-- Delete employers
DELETE FROM employers
WHERE id IN (
    '12345678-0008-0008-0008-000000000201'::uuid,
    '12345678-0008-0008-0008-000000000202'::uuid,
    '12345678-0008-0008-0008-000000000203'::uuid,
    '12345678-0008-0008-0008-000000000204'::uuid,
    '12345678-0008-0008-0008-000000000205'::uuid
);

-- Delete emails
DELETE FROM emails
WHERE email_key IN (
    '12345678-0008-0008-0008-000000000011'::uuid,
    '12345678-0008-0008-0008-000000000012'::uuid,
    '12345678-0008-0008-0008-000000000013'::uuid,
    '12345678-0008-0008-0008-000000000014'::uuid,
    '12345678-0008-0008-0008-000000000015'::uuid
);

-- Delete hub user tokens
DELETE FROM hub_user_tokens
WHERE hub_user_id IN (
    SELECT id FROM hub_users
    WHERE email = 'hubopening@hub.example'
);

-- Delete hub users
DELETE FROM hub_users
WHERE id IN (
    '12345678-0008-0008-0008-000000050001'::uuid,
    '12345678-0008-0008-0008-000000050002'::uuid,
    '12345678-0008-0008-0008-000000050003'::uuid
);

COMMIT;
