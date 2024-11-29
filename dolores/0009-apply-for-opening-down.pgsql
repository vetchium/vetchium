BEGIN;

-- Delete applications
DELETE FROM applications
WHERE employer_id = '12345678-0009-0009-0009-000000000201'::uuid;

-- Delete opening locations
DELETE FROM opening_locations
WHERE employer_id = '12345678-0009-0009-0009-000000000201'::uuid;

-- Delete openings
DELETE FROM openings
WHERE employer_id = '12345678-0009-0009-0009-000000000201'::uuid;

-- Delete locations
DELETE FROM locations
WHERE employer_id = '12345678-0009-0009-0009-000000000201'::uuid;

-- Delete cost centers
DELETE FROM org_cost_centers
WHERE employer_id = '12345678-0009-0009-0009-000000000201'::uuid;

-- Delete org user tokens
DELETE FROM org_user_tokens
WHERE org_user_id IN (
    SELECT id FROM org_users
    WHERE employer_id = '12345678-0009-0009-0009-000000000201'::uuid
);

-- Delete org users
DELETE FROM org_users
WHERE employer_id = '12345678-0009-0009-0009-000000000201'::uuid;

-- Delete domains
DELETE FROM domains
WHERE employer_id = '12345678-0009-0009-0009-000000000201'::uuid;

-- Delete employer
DELETE FROM employers
WHERE id = '12345678-0009-0009-0009-000000000201'::uuid;

-- Delete email
DELETE FROM emails
WHERE email_key = '12345678-0009-0009-0009-000000000011'::uuid;

-- Delete hub user tokens
DELETE FROM hub_user_tokens
WHERE hub_user_id IN (
    '12345678-0009-0009-0009-000000050001'::uuid,
    '12345678-0009-0009-0009-000000050002'::uuid,
    '12345678-0009-0009-0009-000000050003'::uuid
);

-- Delete hub users
DELETE FROM hub_users
WHERE id IN (
    '12345678-0009-0009-0009-000000050001'::uuid,
    '12345678-0009-0009-0009-000000050002'::uuid,
    '12345678-0009-0009-0009-000000050003'::uuid
);

COMMIT;
