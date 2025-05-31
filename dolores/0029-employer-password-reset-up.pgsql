-- Set up test data for employer password reset tests

-- Create employer domain
INSERT INTO employers (
    id,
    client_id_type,
    employer_state,
    company_name,
    onboard_admin_email,
    created_at
) VALUES (
    '02900000-0000-0000-0000-000000000000',
    'DOMAIN',
    'ONBOARDED',
    'Password Reset Test Company',
    'admin@passwordreset.example',
    timezone('UTC', now())
);

INSERT INTO domains (
    id,
    domain_name,
    domain_state,
    employer_id,
    created_at
) VALUES (
    '02900000-0000-0000-0000-000000000001',
    'passwordreset.example',
    'VERIFIED',
    '02900000-0000-0000-0000-000000000000',
    timezone('UTC', now())
);

INSERT INTO employer_primary_domains (
    employer_id,
    domain_id
) VALUES (
    '02900000-0000-0000-0000-000000000000',
    '02900000-0000-0000-0000-000000000001'
);

-- Create test org users for password reset scenarios
INSERT INTO org_users (
    id,
    email,
    name,
    password_hash,
    org_user_roles,
    org_user_state,
    employer_id,
    created_at
) VALUES
-- Active user for basic password reset tests
(
    '02900001-0000-0000-0000-000000000000',
    'active@passwordreset.example',
    'Active User',
    '$2a$10$rFJF1QKnxNn7YLCOqVkN9.xHVP5Z8z9K5rJQ3vV7X1uG3K9Q8J5Dm', -- NewPassword123$
    ARRAY['ADMIN']::org_user_roles[],
    'ACTIVE_ORG_USER',
    '02900000-0000-0000-0000-000000000000',
    timezone('UTC', now())
),
-- Disabled user - should still be able to reset password
(
    '02900002-0000-0000-0000-000000000000',
    'disabled@passwordreset.example',
    'Disabled User',
    '$2a$10$rFJF1QKnxNn7YLCOqVkN9.xHVP5Z8z9K5rJQ3vV7X1uG3K9Q8J5Dm', -- NewPassword123$
    ARRAY['ORG_USERS_VIEWER']::org_user_roles[],
    'DISABLED_ORG_USER',
    '02900000-0000-0000-0000-000000000000',
    timezone('UTC', now())
),
-- User for multiple requests test
(
    '02900003-0000-0000-0000-000000000000',
    'multiple-requests@passwordreset.example',
    'Multiple Requests User',
    '$2a$10$rFJF1QKnxNn7YLCOqVkN9.xHVP5Z8z9K5rJQ3vV7X1uG3K9Q8J5Dm', -- NewPassword123$
    ARRAY['ORG_USERS_CRUD']::org_user_roles[],
    'ACTIVE_ORG_USER',
    '02900000-0000-0000-0000-000000000000',
    timezone('UTC', now())
),
-- User for valid reset test
(
    '02900004-0000-0000-0000-000000000000',
    'valid-reset@passwordreset.example',
    'Valid Reset User',
    '$2a$10$rFJF1QKnxNn7YLCOqVkN9.xHVP5Z8z9K5rJQ3vV7X1uG3K9Q8J5Dm', -- NewPassword123$
    ARRAY['ORG_USERS_VIEWER']::org_user_roles[],
    'ACTIVE_ORG_USER',
    '02900000-0000-0000-0000-000000000000',
    timezone('UTC', now())
),
-- User for token expiry test
(
    '02900005-0000-0000-0000-000000000000',
    'token-expiry@passwordreset.example',
    'Token Expiry User',
    '$2a$10$rFJF1QKnxNn7YLCOqVkN9.xHVP5Z8z9K5rJQ3vV7X1uG3K9Q8J5Dm', -- NewPassword123$
    ARRAY['ORG_USERS_VIEWER']::org_user_roles[],
    'ACTIVE_ORG_USER',
    '02900000-0000-0000-0000-000000000000',
    timezone('UTC', now())
),
-- User for token reuse test
(
    '02900006-0000-0000-0000-000000000000',
    'token-reuse@passwordreset.example',
    'Token Reuse User',
    '$2a$10$rFJF1QKnxNn7YLCOqVkN9.xHVP5Z8z9K5rJQ3vV7X1uG3K9Q8J5Dm', -- NewPassword123$
    ARRAY['ORG_USERS_VIEWER']::org_user_roles[],
    'ACTIVE_ORG_USER',
    '02900000-0000-0000-0000-000000000000',
    timezone('UTC', now())
),
-- User for cross-employer test
(
    '02900007-0000-0000-0000-000000000000',
    'cross-test@passwordreset.example',
    'Cross Test User',
    '$2a$10$rFJF1QKnxNn7YLCOqVkN9.xHVP5Z8z9K5rJQ3vV7X1uG3K9Q8J5Dm', -- NewPassword123$
    ARRAY['ORG_USERS_VIEWER']::org_user_roles[],
    'ACTIVE_ORG_USER',
    '02900000-0000-0000-0000-000000000000',
    timezone('UTC', now())
),
-- User for session validity test
(
    '02900008-0000-0000-0000-000000000000',
    'session-test@passwordreset.example',
    'Session Test User',
    '$2a$10$rFJF1QKnxNn7YLCOqVkN9.xHVP5Z8z9K5rJQ3vV7X1uG3K9Q8J5Dm', -- NewPassword123$
    ARRAY['ORG_USERS_VIEWER']::org_user_roles[],
    'ACTIVE_ORG_USER',
    '02900000-0000-0000-0000-000000000000',
    timezone('UTC', now())
); 