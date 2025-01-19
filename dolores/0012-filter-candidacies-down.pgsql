BEGIN;

-- Delete candidacies
DELETE FROM candidacies
WHERE employer_id = '12345678-0011-0011-0011-000000000201'::uuid;

-- Delete applications
DELETE FROM applications
WHERE employer_id = '12345678-0011-0011-0011-000000000201'::uuid;

-- Delete openings
DELETE FROM openings
WHERE employer_id = '12345678-0011-0011-0011-000000000201'::uuid;

-- Delete cost centers
DELETE FROM org_cost_centers
WHERE employer_id = '12345678-0011-0011-0011-000000000201'::uuid;

-- Delete org user tokens
DELETE FROM org_user_tokens
WHERE org_user_id IN (
    SELECT id FROM org_users
    WHERE employer_id = '12345678-0011-0011-0011-000000000201'::uuid
);

-- Delete org users
DELETE FROM org_users
WHERE employer_id = '12345678-0011-0011-0011-000000000201'::uuid;

-- Delete domains
DELETE FROM domains
WHERE employer_id = '12345678-0011-0011-0011-000000000201'::uuid;

-- Delete employer
DELETE FROM employers
WHERE id = '12345678-0011-0011-0011-000000000201'::uuid;

-- Delete emails
DELETE FROM emails
WHERE email_key = '12345678-0011-0011-0011-000000000011'::uuid;

-- Delete hub users
DELETE FROM hub_users
WHERE id IN (
    '12345678-0011-0011-0011-000000060001'::uuid,
    '12345678-0011-0011-0011-000000060002'::uuid,
    '12345678-0011-0011-0011-000000060003'::uuid
);

COMMIT;
