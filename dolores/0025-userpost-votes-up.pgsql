-- Create test users for voting tests
INSERT INTO hub_users (id, full_name, handle, email, password_hash, state, tier, resident_country_code, preferred_language, short_bio, long_bio, created_at)
VALUES
  ('12345678-0025-0025-0025-000000000001', 'Voter One', 'voter-one', 'voter1@0025-votes.example.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'PAID_HUB_USER', 'IN', 'en', 'Test voter one', 'Test voter one - long bio', NOW()),
  ('12345678-0025-0025-0025-000000000002', 'Voter Two', 'voter-two', 'voter2@0025-votes.example.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'PAID_HUB_USER', 'IN', 'en', 'Test voter two', 'Test voter two - long bio', NOW()),
  ('12345678-0025-0025-0025-000000000003', 'Post Author', 'post-author', 'author@0025-votes.example.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'PAID_HUB_USER', 'IN', 'en', 'Test post author', 'Test post author - long bio', NOW());
