-- Seed data for 0035-incognito-comments_test.go

-- Insert test users with UUID pattern 12345678-0035-0035-0035-xxxxxxxxxxxx
INSERT INTO hub_users (id, full_name, handle, email, password_hash, state, tier, resident_country_code, resident_city, preferred_language, short_bio, long_bio, created_at)
VALUES 
    ('12345678-0035-0035-0035-000000000001', 'User 0035-1', 'user0035001', 'user0035-1@0035-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()),
    ('12345678-0035-0035-0035-000000000002', 'User 0035-2', 'user0035002', 'user0035-2@0035-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()),
    ('12345678-0035-0035-0035-000000000003', 'User 0035-3', 'user0035003', 'user0035-3@0035-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()),
    ('12345678-0035-0035-0035-000000000004', 'User 0035-4', 'user0035004', 'user0035-4@0035-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()),
    ('12345678-0035-0035-0035-000000000005', 'User 0035-5', 'user0035005', 'user0035-5@0035-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()),
    ('12345678-0035-0035-0035-000000000006', 'User 0035-6', 'user0035006', 'user0035-6@0035-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()),
    ('12345678-0035-0035-0035-000000000007', 'User 0035-7', 'user0035007', 'user0035-7@0035-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()),
    ('12345678-0035-0035-0035-000000000008', 'User 0035-8', 'user0035008', 'user0035-8@0035-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()); 