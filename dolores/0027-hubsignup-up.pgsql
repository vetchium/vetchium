BEGIN;

-- Create domains for testing
INSERT INTO domains (id, domain_name, domain_state, employer_id) VALUES
-- Approved domain
('12345678-0027-0027-0027-000000000001', '0027-example.com', 'VERIFIED', NULL),
-- Unapproved domain
('12345678-0027-0027-0027-000000000002', 'unapproved-0027-example.com', 'UNVERIFIED', NULL);

-- Create hub users for testing existing users
INSERT INTO hub_users (
    id, full_name, handle, email, password_hash,
    state, tier, resident_country_code, resident_city, preferred_language,
    short_bio, long_bio, created_at
) VALUES
-- Existing user for testing duplicate signup
(
    '12345678-0027-0027-0027-000000000003',
    'Existing Test User',
    'existinguser',
    'existing@0027-example.com',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- Password: NewPassword123$
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'IND',
    'Bangalore',
    'en',
    'Existing Test User',
    'Existing Test User for testing duplicate signup attempts',
    timezone('UTC'::text, now())
);

-- Create hub user invites for testing existing invites
INSERT INTO hub_user_invites (email, token, token_valid_till) VALUES
-- Existing invite
('invited@0027-example.com', 'existing-invite-token', timezone('UTC', now()) + interval '24 hours');

-- Create approved domains list for testing domain validation
INSERT INTO hub_user_signup_approved_domains (domain_name, notes) VALUES
('0027-example.com', 'Test domain for 0027 test'),
('another-0027-example.com', 'Another test domain for 0027 test');

COMMIT;
