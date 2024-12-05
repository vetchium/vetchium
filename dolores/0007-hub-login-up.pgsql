BEGIN;

INSERT INTO public.hub_users (
    id,
    full_name,
    handle,
    email,
    password_hash,
    state,
    resident_country_code,
    resident_city,
    preferred_language,
    created_at,
    updated_at
)
VALUES 
    (
        '12345678-0007-0007-0007-000000050001'::uuid,
        'Active User',
        'active_user',
        'active@hub.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        'ACTIVE_HUB_USER',
        'USA',
        'New York',
        'en',
        timezone('UTC'::text, now()),
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0007-0007-0007-000000050002'::uuid,
        'Disabled User',
        'disabled_user',
        'disabled@hub.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        'DISABLED_HUB_USER',
        'USA',
        'New York',
        'en',
        timezone('UTC'::text, now()),
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0007-0007-0007-000000050003'::uuid,
        'Deleted User',
        'deleted_user',
        'deleted@hub.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        'DELETED_HUB_USER',
        'USA',
        'San Francisco',
        'en',
        timezone('UTC'::text, now()),
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0007-0007-0007-000000050004'::uuid,
        'Password Change User',
        'password_change_user',
        'password-change@hub.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        'ACTIVE_HUB_USER',
        'IND',
        'Bangalore',
        'en',
        timezone('UTC'::text, now()),
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0007-0007-0007-000000050005'::uuid,
        'Password Reset User',
        'password_reset_user',
        'password-reset@hub.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        'ACTIVE_HUB_USER',
        'IND',
        'Chennai',
        'en',
        timezone('UTC'::text, now()),
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0007-0007-0007-000000050006'::uuid,
        'Token Expiry Test User',
        'token_expiry_user',
        'token-expiry@hub.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        'ACTIVE_HUB_USER',
        'GER',
        'Nuremberg',
        'en',
        timezone('UTC'::text, now()),
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0007-0007-0007-000000050007'::uuid,
        'Token Reuse Test User',
        'token_reuse_user',
        'token-reuse@hub.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        'ACTIVE_HUB_USER',
        'PRC',
        'Shanghai',
        'en',
        timezone('UTC'::text, now()),
        timezone('UTC'::text, now())
    ),(
        '12345678-0007-0007-0007-000000050008'::uuid,
        'Remember Me Test User',
        'rememberme_user',
        'rememberme@hub.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        'ACTIVE_HUB_USER',
        'FRA',
        'Paris',
        'en',
        timezone('UTC'::text, now()),
        timezone('UTC'::text, now())
    ),
    (
        '12345678-0007-0007-0007-000000050009'::uuid,
        'TFA Test User',
        'tfatest_user',
        'tfatest@hub.example',
        '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
        'ACTIVE_HUB_USER',
        'JPN',
        'Tokyo',
        'en',
        timezone('UTC'::text, now()),
        timezone('UTC'::text, now())
    );
COMMIT;
