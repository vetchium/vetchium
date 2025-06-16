-- Seed data for 0036-incognito-voting_test.go

-- Insert test users with UUID pattern 12345678-0036-0036-0036-xxxxxxxxxxxx
INSERT INTO hub_users (id, full_name, handle, email, password_hash, state, tier, resident_country_code, resident_city, preferred_language, short_bio, long_bio, created_at)
VALUES 
    ('12345678-0036-0036-0036-000000000001', 'User 0036-1', 'user0036001', 'user0036-1@0036-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()),
    ('12345678-0036-0036-0036-000000000002', 'User 0036-2', 'user0036002', 'user0036-2@0036-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()),
    ('12345678-0036-0036-0036-000000000003', 'User 0036-3', 'user0036003', 'user0036-3@0036-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()),
    ('12345678-0036-0036-0036-000000000004', 'User 0036-4', 'user0036004', 'user0036-4@0036-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()),
    ('12345678-0036-0036-0036-000000000005', 'User 0036-5', 'user0036005', 'user0036-5@0036-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()),
    ('12345678-0036-0036-0036-000000000006', 'User 0036-6', 'user0036006', 'user0036-6@0036-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()),
    ('12345678-0036-0036-0036-000000000007', 'User 0036-7', 'user0036007', 'user0036-7@0036-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()),
    ('12345678-0036-0036-0036-000000000008', 'User 0036-8', 'user0036008', 'user0036-8@0036-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()); 