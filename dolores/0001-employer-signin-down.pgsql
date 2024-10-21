BEGIN;

DELETE FROM domains WHERE employer_id = '12345678-0001-0001-0001-000000000201'::UUID;
DELETE FROM org_user_tokens WHERE org_user_id IN (
    SELECT id FROM org_users WHERE employer_id = '12345678-0001-0001-0001-000000000201'::UUID
);
DELETE FROM org_users WHERE employer_id = '12345678-0001-0001-0001-000000000201'::UUID;
DELETE FROM employers WHERE id = '12345678-0001-0001-0001-000000000201'::UUID;
DELETE FROM emails WHERE email_key = '12345678-0001-0001-0001-000000000011'::UUID;

DELETE FROM domains WHERE domain_name IN ('secretsapp.com', 'aadal.in');
DELETE FROM org_user_tokens WHERE org_user_id IN (
    SELECT id FROM org_users WHERE employer_id IN (
        SELECT id FROM employers WHERE onboard_admin_email IN ('secretsapp@example.com', 'aadal@example.com')
    )
);
DELETE FROM org_users WHERE employer_id IN (
    SELECT id FROM employers WHERE onboard_admin_email IN ('secretsapp@example.com', 'aadal@example.com')
);
DELETE FROM employers WHERE onboard_admin_email IN ('secretsapp@example.com', 'aadal@example.com');

COMMIT;