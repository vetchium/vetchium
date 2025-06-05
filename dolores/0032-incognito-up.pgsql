BEGIN;

-- Test data for incognito posts testing

-- Insert test approved domains
INSERT INTO hub_user_signup_approved_domains (domain_name, notes)
VALUES 
    ('test0032.com', 'Test domain for incognito posts testing'),
    ('company0032.com', 'Another test domain for incognito posts testing');

-- Insert test hub users
INSERT INTO hub_users (
    id, full_name, handle, email, password_hash, state, tier, 
    resident_country_code, resident_city, preferred_language, short_bio, long_bio
) VALUES 
    ('12345678-0032-0032-0032-000000000001', 'Alice Test0032', 'alice0032test', 'alice@test0032.com', 
     '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio'),
    ('12345678-0032-0032-0032-000000000002', 'Bob Test0032', 'bob0032test', 'bob@test0032.com',
     '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'PAID_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio'),
    ('12345678-0032-0032-0032-000000000003', 'Charlie Test0032', 'charlie0032test', 'charlie@company0032.com',
     '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio'),
    ('12345678-0032-0032-0032-000000000004', 'Diana Test0032', 'diana0032test', 'diana@test0032.com',
     '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'DISABLED_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio'),
    ('12345678-0032-0032-0032-000000000005', 'Eve Test0032', 'eve0032test', 'eve@test0032.com',
     '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio'),
    ('12345678-0032-0032-0032-000000000006', 'Frank Test0032', 'frank0032test', 'frank@test0032.com',
     '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio');

-- Insert test tags for incognito posts testing (if they don't already exist)
INSERT INTO tags (id, display_name) VALUES 
    ('technology', 'Technology'),
    ('careers', 'Careers'),
    ('personal-development', 'Personal Development'),
    ('mentorship', 'Mentorship')
ON CONFLICT (id) DO NOTHING;

COMMIT; 