BEGIN;
DELETE FROM hub_user_tokens 
WHERE hub_user_id IN (
    SELECT id FROM hub_users 
    WHERE email IN (
        'active@hub.example',
        'disabled@hub.example', 
        'deleted@hub.example',
        'password-change@hub.example'
    )
);

DELETE FROM hub_users 
WHERE email IN (
    'active@hub.example',
    'disabled@hub.example',
    'deleted@hub.example',
    'password-change@hub.example'
);
COMMIT;
