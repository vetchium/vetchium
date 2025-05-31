BEGIN;

-- Clean up test data
DELETE FROM hub_user_signup_approved_domains WHERE domain_name IN ('registration.example', 'test-registration.example');

-- Clean up any test users that might have been created
DELETE FROM hub_users WHERE email LIKE '%@registration.example';

COMMIT;