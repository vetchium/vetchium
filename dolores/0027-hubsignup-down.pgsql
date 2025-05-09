BEGIN;

-- Clean up emails table
DELETE FROM emails
WHERE email_to && ARRAY[
    'new@0027-example.com',
    'existing@0027-example.com',
    'invited@0027-example.com',
    'invalid@unapproved-0027-example.com',
    'another@0027-example.com'
];

-- Clean up hub user invites
DELETE FROM hub_user_invites
WHERE email IN (
    'new@0027-example.com',
    'existing@0027-example.com',
    'invited@0027-example.com',
    'invalid@unapproved-0027-example.com',
    'another@0027-example.com'
);

-- Clean up hub user tokens and related tables first
DELETE FROM hub_user_tfa_codes
WHERE tfa_token IN (
    SELECT token FROM hub_user_tokens
    WHERE hub_user_id = '12345678-0027-0027-0027-000000000003'
);

DELETE FROM hub_user_tokens
WHERE hub_user_id = '12345678-0027-0027-0027-000000000003';

-- Clean up hub users
DELETE FROM hub_users
WHERE id = '12345678-0027-0027-0027-000000000003';

-- Clean up approved domains
DELETE FROM approved_domains
WHERE domain IN ('0027-example.com', 'another-0027-example.com');

-- Clean up domains
DELETE FROM domains
WHERE id IN (
    '12345678-0027-0027-0027-000000000001',
    '12345678-0027-0027-0027-000000000002'
);

COMMIT;
