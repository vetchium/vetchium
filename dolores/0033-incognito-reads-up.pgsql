BEGIN;

-- Test data for incognito reads testing - only basic user setup
-- All posts, comments, and votes will be created via APIs in tests

-- Insert test approved domains
INSERT INTO hub_user_signup_approved_domains (domain_name, notes)
VALUES 
    ('test0033.com', 'Test domain for incognito reads testing'),
    ('company0033.com', 'Another test domain for incognito reads testing');

-- Insert test hub users with varying tiers and states
INSERT INTO hub_users (
    id, full_name, handle, email, password_hash, state, tier, 
    resident_country_code, resident_city, preferred_language, short_bio, long_bio
) VALUES 
    ('12345678-0033-0033-0033-000000000001', 'Alice Reads0033', 'alice0033reads', 'alice@test0033.com', 
     '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Alice reads bio', 'Alice long bio'),
    ('12345678-0033-0033-0033-000000000002', 'Bob Reads0033', 'bob0033reads', 'bob@test0033.com',
     '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'PAID_HUB_USER', 'US', 'Test City', 'en', 'Bob reads bio', 'Bob long bio'),
    ('12345678-0033-0033-0033-000000000003', 'Charlie Reads0033', 'charlie0033reads', 'charlie@company0033.com',
     '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Charlie reads bio', 'Charlie long bio'),
    ('12345678-0033-0033-0033-000000000005', 'Eve Reads0033', 'eve0033reads', 'eve@test0033.com',
     '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Eve reads bio', 'Eve long bio'),
    ('12345678-0033-0033-0033-000000000006', 'Frank Reads0033', 'frank0033reads', 'frank@test0033.com',
     '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Frank reads bio', 'Frank long bio'),
    ('12345678-0033-0033-0033-000000000007', 'Grace Reads0033', 'grace0033reads', 'grace@test0033.com',
     '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Grace reads bio', 'Grace long bio');

COMMIT;
