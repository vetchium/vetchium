BEGIN;

-- Delete applications
DELETE FROM applications 
WHERE employer_id IN (
    SELECT id FROM employers 
    WHERE onboard_admin_email LIKE '%@applied%.example'
);

-- Delete openings
DELETE FROM openings 
WHERE employer_id IN (
    SELECT id FROM employers 
    WHERE onboard_admin_email LIKE '%@applied%.example'
);

-- Delete org_user_tokens
DELETE FROM org_user_tokens 
WHERE org_user_id IN (
    SELECT id FROM org_users 
    WHERE employer_id IN (
        SELECT id FROM employers 
        WHERE onboard_admin_email LIKE '%@applied%.example'
    )
);

-- Delete org_users
DELETE FROM org_users 
WHERE employer_id IN (
    SELECT id FROM employers 
    WHERE onboard_admin_email LIKE '%@applied%.example'
);

-- Delete domains
DELETE FROM domains 
WHERE employer_id IN (
    SELECT id FROM employers 
    WHERE onboard_admin_email LIKE '%@applied%.example'
);

-- Delete employers
DELETE FROM employers 
WHERE onboard_admin_email LIKE '%@applied%.example';

-- Delete emails
DELETE FROM emails 
WHERE email_key IN (
    '12345678-0010-0010-0010-000000000011'::uuid,
    '12345678-0010-0010-0010-000000000012'::uuid,
    '12345678-0010-0010-0010-000000000013'::uuid,
    '12345678-0010-0010-0010-000000000014'::uuid,
    '12345678-0010-0010-0010-000000000015'::uuid
);

-- Delete hub user tokens
DELETE FROM hub_user_tokens 
WHERE hub_user_id IN (
    '12345678-0010-0010-0010-000000050001'::uuid,
    '12345678-0010-0010-0010-000000050002'::uuid
);

-- Delete hub users
DELETE FROM hub_users 
WHERE id IN (
    '12345678-0010-0010-0010-000000050001'::uuid,
    '12345678-0010-0010-0010-000000050002'::uuid
);

COMMIT;
