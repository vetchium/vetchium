-- Set up test data for 0029 employer password reset tests

-- Create 0029 employer domain
INSERT INTO employers (
    id,
    client_id_type,
    employer_state,
    company_name,
    onboard_admin_email,
    created_at
) VALUES (
    '02900000-0029-0029-0029-000000000000',
    'DOMAIN',
    'ONBOARDED',
    '0029 Password Reset Test Company',
    'admin@0029-passwordreset.example',
    timezone('UTC', now())
);

INSERT INTO domains (
    id,
    domain_name,
    domain_state,
    employer_id,
    created_at
) VALUES (
    '02900000-0029-0029-0029-000000000001',
    '0029-passwordreset.example',
    'VERIFIED',
    '02900000-0029-0029-0029-000000000000',
    timezone('UTC', now())
);

INSERT INTO employer_primary_domains (
    employer_id,
    domain_id
) VALUES (
    '02900000-0029-0029-0029-000000000000',
    '02900000-0029-0029-0029-000000000001'
);

-- Create test org users for password reset scenarios (each test case gets unique email)
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
-- User for "should handle various forgot password scenarios" test
(
    '02900001-0029-0029-0029-000000000000',
    'test001-forgot-scenarios@0029-passwordreset.example',
    '0029 Forgot Scenarios User',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    ARRAY['ADMIN']::org_user_roles[],
    'ACTIVE_ORG_USER',
    '02900000-0029-0029-0029-000000000000',
    timezone('UTC', now())
),
-- User for "should invalidate previous tokens when multiple requests are made" test
(
    '02900002-0029-0029-0029-000000000000',
    'test002-multiple-requests@0029-passwordreset.example',
    '0029 Multiple Requests User',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    ARRAY['ORG_USERS_CRUD']::org_user_roles[],
    'ACTIVE_ORG_USER',
    '02900000-0029-0029-0029-000000000000',
    timezone('UTC', now())
),
-- User for "should handle various reset password scenarios" test
(
    '02900003-0029-0029-0029-000000000000',
    'test003-reset-scenarios@0029-passwordreset.example',
    '0029 Reset Scenarios User',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    ARRAY['ORG_USERS_VIEWER']::org_user_roles[],
    'ACTIVE_ORG_USER',
    '02900000-0029-0029-0029-000000000000',
    timezone('UTC', now())
),
-- User for "should handle token expiry" test
(
    '02900004-0029-0029-0029-000000000000',
    'test004-token-expiry@0029-passwordreset.example',
    '0029 Token Expiry User',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    ARRAY['ORG_USERS_VIEWER']::org_user_roles[],
    'ACTIVE_ORG_USER',
    '02900000-0029-0029-0029-000000000000',
    timezone('UTC', now())
),
-- User for "should prevent token reuse" test
(
    '02900005-0029-0029-0029-000000000000',
    'test005-token-reuse@0029-passwordreset.example',
    '0029 Token Reuse User',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    ARRAY['ORG_USERS_VIEWER']::org_user_roles[],
    'ACTIVE_ORG_USER',
    '02900000-0029-0029-0029-000000000000',
    timezone('UTC', now())
),
-- User for "should handle cross-employer token attempts" test
(
    '02900006-0029-0029-0029-000000000000',
    'test006-cross-employer@0029-passwordreset.example',
    '0029 Cross Employer User',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    ARRAY['ORG_USERS_VIEWER']::org_user_roles[],
    'ACTIVE_ORG_USER',
    '02900000-0029-0029-0029-000000000000',
    timezone('UTC', now())
),
-- User for "should maintain session validity after password reset" test
(
    '02900007-0029-0029-0029-000000000000',
    'test007-session-validity@0029-passwordreset.example',
    '0029 Session Validity User',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    ARRAY['ORG_USERS_VIEWER']::org_user_roles[],
    'ACTIVE_ORG_USER',
    '02900000-0029-0029-0029-000000000000',
    timezone('UTC', now())
),
-- Disabled user for testing that disabled users can still reset password
(
    '02900008-0029-0029-0029-000000000000',
    'test008-disabled@0029-passwordreset.example',
    '0029 Disabled User',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    ARRAY['ORG_USERS_VIEWER']::org_user_roles[],
    'DISABLED_ORG_USER',
    '02900000-0029-0029-0029-000000000000',
    timezone('UTC', now())
); 