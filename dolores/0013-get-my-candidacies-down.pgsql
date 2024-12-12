BEGIN;

DELETE FROM candidacies
WHERE employer_id IN (
    SELECT id FROM employers
    WHERE onboard_admin_email LIKE '%@my-candidacies-%.example'
);

-- Delete applications
DELETE FROM applications
WHERE hub_user_id IN (
    SELECT id FROM hub_users
    WHERE email LIKE '%@my-candidacies.example'
);

-- Delete openings
DELETE FROM openings
WHERE employer_id IN (
    SELECT id FROM employers
    WHERE onboard_admin_email LIKE '%@my-candidacies-%.example'
);

-- Delete cost centers
DELETE FROM org_cost_centers
WHERE employer_id IN (
    SELECT id FROM employers
    WHERE onboard_admin_email LIKE '%@my-candidacies-%.example'
);

-- Delete org user tokens
DELETE FROM org_user_tokens
WHERE org_user_id IN (
    SELECT id FROM org_users
    WHERE employer_id IN (
        SELECT id FROM employers
        WHERE onboard_admin_email LIKE '%@my-candidacies-%.example'
    )
);

-- Delete org users
DELETE FROM org_users
WHERE employer_id IN (
    SELECT id FROM employers
    WHERE onboard_admin_email LIKE '%@my-candidacies-%.example'
);

-- Delete domains
DELETE FROM domains
WHERE employer_id IN (
    SELECT id FROM employers
    WHERE onboard_admin_email LIKE '%@my-candidacies-%.example'
);

-- Delete employers
DELETE FROM employers
WHERE onboard_admin_email LIKE '%@my-candidacies-%.example';

-- Delete emails
DELETE FROM emails
WHERE email_key IN (
    '12345678-0013-0013-0013-000000000011'::uuid,
    '12345678-0013-0013-0013-000000000012'::uuid,
    '12345678-0013-0013-0013-000000000013'::uuid
);

-- Delete hub user tokens
DELETE FROM hub_user_tokens
WHERE hub_user_id IN (
    SELECT id FROM hub_users
    WHERE email LIKE '%@my-candidacies.example'
);

-- Delete hub users
DELETE FROM hub_users
WHERE email LIKE '%@my-candidacies.example';

COMMIT;
