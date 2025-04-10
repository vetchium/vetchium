BEGIN;

DELETE FROM posts WHERE author_id IN (
    '12345678-0022-0022-0022-000000000001',
    '12345678-0022-0022-0022-000000000002'
);

-- Clean up hub user tokens and related tables first
DELETE FROM hub_user_tfa_codes
WHERE tfa_token IN (
    SELECT token FROM hub_user_tokens
    WHERE hub_user_id IN (
        '12345678-0022-0022-0022-000000000001',
        '12345678-0022-0022-0022-000000000002'
    )
);

DELETE FROM hub_user_tokens
WHERE hub_user_id IN (
    '12345678-0022-0022-0022-000000000001',
    '12345678-0022-0022-0022-000000000002'
);

-- Clean up hub users
DELETE FROM hub_users
WHERE id IN (
    '12345678-0022-0022-0022-000000000001',
    '12345678-0022-0022-0022-000000000002'
);


COMMIT;
