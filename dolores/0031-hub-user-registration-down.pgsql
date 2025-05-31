BEGIN;

-- Clean up test data - be very thorough
-- First clean up tokens that might reference hub users
DELETE FROM hub_user_tokens WHERE hub_user_id IN (SELECT id FROM hub_users WHERE email LIKE '%@0031-registration.example' OR email LIKE '%@0031-test-registration.example');

-- Clean up any invites for test domains  
DELETE FROM hub_user_invites WHERE email LIKE '%@0031-registration.example' OR email LIKE '%@0031-test-registration.example';

-- Clean up any test users that might have been created
DELETE FROM hub_users WHERE email LIKE '%@0031-registration.example' OR email LIKE '%@0031-test-registration.example';

-- Clean up emails sent to test domains
DELETE FROM emails WHERE EXISTS (
    SELECT 1 FROM UNNEST(email_to) AS recipient 
    WHERE recipient LIKE '%@0031-registration.example' OR recipient LIKE '%@0031-test-registration.example'
);

-- Clean up the approved domains
DELETE FROM hub_user_signup_approved_domains WHERE domain_name IN ('0031-registration.example', '0031-test-registration.example');

COMMIT;