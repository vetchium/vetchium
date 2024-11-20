BEGIN;
DELETE FROM hub_user_tokens 
WHERE hub_user_id IN (
    SELECT id FROM hub_users 
    WHERE email IN (
        'active@hub.example',
        'disabled@hub.example', 
        'deleted@hub.example',
        'password-change@hub.example',
        'password-reset@hub.example',
        'token-expiry@hub.example',
        'token-reuse@hub.example',
        'rememberme@hub.example',
        'tfatest@hub.example'
    )
);

DELETE FROM hub_users 
WHERE email IN (
    'active@hub.example',
    'disabled@hub.example',
    'deleted@hub.example',
    'password-change@hub.example',
    'password-reset@hub.example',
    'token-expiry@hub.example',
    'token-reuse@hub.example',
    'rememberme@hub.example',
    'tfatest@hub.example'
);
COMMIT;
