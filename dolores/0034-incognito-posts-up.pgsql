-- Seed data for 0034-incognito-posts_test.go

-- Insert test users with UUID pattern 12345678-0034-0034-0034-xxxxxxxxxxxx
INSERT INTO hub_users (id, full_name, handle, email, password_hash, state, tier, resident_country_code, resident_city, preferred_language, short_bio, long_bio, created_at)
VALUES 
    ('12345678-0034-0034-0034-000000000001', 'User 0034-1', 'user0034001', 'user0034-1@0034-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()),
    ('12345678-0034-0034-0034-000000000002', 'User 0034-2', 'user0034002', 'user0034-2@0034-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()),
    ('12345678-0034-0034-0034-000000000003', 'User 0034-3', 'user0034003', 'user0034-3@0034-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()),
    ('12345678-0034-0034-0034-000000000004', 'User 0034-4', 'user0034004', 'user0034-4@0034-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()),
    ('12345678-0034-0034-0034-000000000005', 'User 0034-5', 'user0034005', 'user0034-5@0034-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()),
    ('12345678-0034-0034-0034-000000000006', 'User 0034-6', 'user0034006', 'user0034-6@0034-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()),
    ('12345678-0034-0034-0034-000000000007', 'User 0034-7', 'user0034007', 'user0034-7@0034-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()),
    ('12345678-0034-0034-0034-000000000008', 'User 0034-8', 'user0034008', 'user0034-8@0034-test.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'Test City', 'en', 'Short bio', 'Long bio', NOW()); 