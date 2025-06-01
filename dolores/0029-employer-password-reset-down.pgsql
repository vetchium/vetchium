-- Cleanup test data for 0029 employer password reset tests

-- Delete any tokens that might have been created during tests
DELETE FROM org_user_tokens
WHERE org_user_id IN (
    SELECT id FROM org_users
    WHERE employer_id = '02900000-0029-0029-0029-000000000000'
);

-- Delete any TFA codes that might have been created
DELETE FROM org_user_tfa_codes
WHERE tfa_token IN (
    SELECT token FROM org_user_tokens
    WHERE org_user_id IN (
        SELECT id FROM org_users
        WHERE employer_id = '02900000-0029-0029-0029-000000000000'
    )
);

-- Delete org users
DELETE FROM org_users
WHERE employer_id = '02900000-0029-0029-0029-000000000000';

-- Delete employer primary domain
DELETE FROM employer_primary_domains
WHERE employer_id = '02900000-0029-0029-0029-000000000000';

-- Delete domain
DELETE FROM domains
WHERE id = '02900000-0029-0029-0029-000000000001';

-- Delete employer
DELETE FROM employers
WHERE id = '02900000-0029-0029-0029-000000000000';

-- Delete any emails that might have been sent during tests to 0029 domain
DELETE FROM emails
WHERE email_to && ARRAY[
    'test001-forgot-scenarios@0029-passwordreset.example',
    'test002-multiple-requests@0029-passwordreset.example',
    'test003-reset-scenarios@0029-passwordreset.example',
    'test004-token-expiry@0029-passwordreset.example',
    'test005-token-reuse@0029-passwordreset.example',
    'test006-cross-employer@0029-passwordreset.example',
    'test007-session-validity@0029-passwordreset.example',
    'test008-disabled@0029-passwordreset.example'
];