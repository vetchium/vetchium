-- Cleanup test data for employer password reset tests

-- Delete any tokens that might have been created during tests
DELETE FROM org_user_tokens 
WHERE org_user_id IN (
    SELECT id FROM org_users 
    WHERE employer_id = '02900000-0000-0000-0000-000000000000'
);

-- Delete any TFA codes that might have been created
DELETE FROM org_user_tfa_codes 
WHERE tfa_token IN (
    SELECT token FROM org_user_tokens 
    WHERE org_user_id IN (
        SELECT id FROM org_users 
        WHERE employer_id = '02900000-0000-0000-0000-000000000000'
    )
);

-- Delete org users
DELETE FROM org_users 
WHERE employer_id = '02900000-0000-0000-0000-000000000000';

-- Delete employer primary domain
DELETE FROM employer_primary_domains 
WHERE employer_id = '02900000-0000-0000-0000-000000000000';

-- Delete domain
DELETE FROM domains 
WHERE id = '02900000-0000-0000-0000-000000000001';

-- Delete employer
DELETE FROM employers 
WHERE id = '02900000-0000-0000-0000-000000000000';

-- Delete any emails that might have been sent during tests
DELETE FROM emails 
WHERE email_to && ARRAY['active@passwordreset.example', 'disabled@passwordreset.example', 'multiple-requests@passwordreset.example', 'valid-reset@passwordreset.example', 'token-expiry@passwordreset.example', 'token-reuse@passwordreset.example', 'cross-test@passwordreset.example', 'session-test@passwordreset.example']; 