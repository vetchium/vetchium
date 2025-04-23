BEGIN;

-- Create hub users for timeline testing
INSERT INTO hub_users (
    id, 
    full_name, 
    handle, 
    email, 
    password_hash, 
    state, 
    tier, 
    resident_country_code, 
    preferred_language, 
    short_bio, 
    long_bio
) VALUES 
    ('12345678-0024-0024-0024-000000000001', 'Timeline User 1', 'timeline-user1-0024', 'user1@0024-timeline-test.example.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'en', 'Short bio 1', 'Long bio 1'),
    ('12345678-0024-0024-0024-000000000002', 'Timeline User 2', 'timeline-user2-0024', 'user2@0024-timeline-test.example.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'en', 'Short bio 2', 'Long bio 2'),
    ('12345678-0024-0024-0024-000000000003', 'Timeline User 3', 'timeline-user3-0024', 'user3@0024-timeline-test.example.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'en', 'Short bio 3', 'Long bio 3'),
    ('12345678-0024-0024-0024-000000000004', 'Timeline User 4', 'timeline-user4-0024', 'user4@0024-timeline-test.example.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'en', 'Short bio 4', 'Long bio 4'),
    ('12345678-0024-0024-0024-000000000005', 'Timeline User 5', 'timeline-user5-0024', 'user5@0024-timeline-test.example.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'en', 'Short bio 5', 'Long bio 5'),
    ('12345678-0024-0024-0024-000000000006', 'Timeline User 6', 'timeline-user6-0024', 'user6@0024-timeline-test.example.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'en', 'Short bio 6', 'Long bio 6'),
    ('12345678-0024-0024-0024-000000000007', 'Timeline User 7', 'timeline-user7-0024', 'user7@0024-timeline-test.example.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'en', 'Short bio 7', 'Long bio 7'),
    ('12345678-0024-0024-0024-000000000008', 'Timeline User 8', 'timeline-user8-0024', 'user8@0024-timeline-test.example.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'en', 'Short bio 8', 'Long bio 8'),
    ('12345678-0024-0024-0024-000000000009', 'Timeline User 9', 'timeline-user9-0024', 'user9@0024-timeline-test.example.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'en', 'Short bio 9', 'Long bio 9'),
    ('12345678-0024-0024-0024-000000000010', 'Timeline User 10', 'timeline-user10-0024', 'user10@0024-timeline-test.example.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'en', 'Short bio 10', 'Long bio 10'),
    ('12345678-0024-0024-0024-000000000011', 'Timeline User 11', 'timeline-user11-0024', 'user11@0024-timeline-test.example.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'en', 'Short bio 11', 'Long bio 11'),
    ('12345678-0024-0024-0024-000000000012', 'Timeline User 12', 'timeline-user12-0024', 'user12@0024-timeline-test.example.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'en', 'Short bio 12', 'Long bio 12'),
    ('12345678-0024-0024-0024-000000000013', 'Timeline User 13', 'timeline-user13-0024', 'user13@0024-timeline-test.example.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'en', 'Short bio 13', 'Long bio 13'),
    ('12345678-0024-0024-0024-000000000014', 'Timeline User 14', 'timeline-user14-0024', 'user14@0024-timeline-test.example.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'en', 'Short bio 14', 'Long bio 14'),
    ('12345678-0024-0024-0024-000000000015', 'Timeline User 15', 'timeline-user15-0024', 'user15@0024-timeline-test.example.com', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'US', 'en', 'Short bio 15', 'Long bio 15');

-- Create some initial follow relationships for testing
INSERT INTO following_relationships (consuming_hub_user_id, producing_hub_user_id)
VALUES 
    -- Initial follow relationships (user1 follows users 2, 3, 4)
    ('12345678-0024-0024-0024-000000000001', '12345678-0024-0024-0024-000000000002'),
    ('12345678-0024-0024-0024-000000000001', '12345678-0024-0024-0024-000000000003'),
    ('12345678-0024-0024-0024-000000000001', '12345678-0024-0024-0024-000000000004'),
    
    -- user5 follows users 6, 7
    ('12345678-0024-0024-0024-000000000005', '12345678-0024-0024-0024-000000000006'),
    ('12345678-0024-0024-0024-000000000005', '12345678-0024-0024-0024-000000000007'),
    
    -- user11 follows users 12, 13, 14
    ('12345678-0024-0024-0024-000000000011', '12345678-0024-0024-0024-000000000012'),
    ('12345678-0024-0024-0024-000000000011', '12345678-0024-0024-0024-000000000013'),
    ('12345678-0024-0024-0024-000000000011', '12345678-0024-0024-0024-000000000014');

-- Create some test tags
INSERT INTO tags (id, name) 
VALUES 
    ('12345678-0024-0024-0024-000000000001', 'timeline-test'),
    ('12345678-0024-0024-0024-000000000002', 'testing'),
    ('12345678-0024-0024-0024-000000000003', 'examples'),
    ('12345678-0024-0024-0024-000000000004', 'automation');

-- Create some initial posts for testing with pre-populated timelines
INSERT INTO posts (id, content, author_id, created_at)
VALUES
    -- Posts by user2 (followed by user1)
    ('post-0024-000000000001', 'This is post 1 from timeline-user2-0024', '12345678-0024-0024-0024-000000000002', NOW() - INTERVAL '3 days'),
    ('post-0024-000000000002', 'This is post 2 from timeline-user2-0024', '12345678-0024-0024-0024-000000000002', NOW() - INTERVAL '2 days'),
    
    -- Posts by user3 (followed by user1)
    ('post-0024-000000000003', 'This is post 1 from timeline-user3-0024', '12345678-0024-0024-0024-000000000003', NOW() - INTERVAL '3 days'),
    ('post-0024-000000000004', 'This is post 2 from timeline-user3-0024', '12345678-0024-0024-0024-000000000003', NOW() - INTERVAL '1 day'),

    -- Posts by user6 (followed by user5)
    ('post-0024-000000000005', 'This is post 1 from timeline-user6-0024', '12345678-0024-0024-0024-000000000006', NOW() - INTERVAL '2 days'),
    ('post-0024-000000000006', 'This is post 2 from timeline-user6-0024', '12345678-0024-0024-0024-000000000006', NOW() - INTERVAL '1 day'),
    
    -- Posts by user12 (followed by user11)
    ('post-0024-000000000007', 'This is post 1 from timeline-user12-0024', '12345678-0024-0024-0024-000000000012', NOW() - INTERVAL '3 days'),
    ('post-0024-000000000008', 'This is post 2 from timeline-user12-0024', '12345678-0024-0024-0024-000000000012', NOW() - INTERVAL '2 days'),
    ('post-0024-000000000009', 'This is post 3 from timeline-user12-0024', '12345678-0024-0024-0024-000000000012', NOW() - INTERVAL '1 day'),
    
    -- Posts by user13 (followed by user11)
    ('post-0024-000000000010', 'This is post 1 from timeline-user13-0024', '12345678-0024-0024-0024-000000000013', NOW() - INTERVAL '2 days'),
    ('post-0024-000000000011', 'This is post 2 from timeline-user13-0024', '12345678-0024-0024-0024-000000000013', NOW() - INTERVAL '1 day');

-- Add tags to some posts
INSERT INTO post_tags (post_id, tag_id)
VALUES
    ('post-0024-000000000001', '12345678-0024-0024-0024-000000000001'),
    ('post-0024-000000000001', '12345678-0024-0024-0024-000000000002'),
    ('post-0024-000000000002', '12345678-0024-0024-0024-000000000003'),
    ('post-0024-000000000004', '12345678-0024-0024-0024-000000000004'),
    ('post-0024-000000000007', '12345678-0024-0024-0024-000000000001'),
    ('post-0024-000000000010', '12345678-0024-0024-0024-000000000002');

-- Create some initial timelines for testing
-- For user1 (following users 2,3,4)
INSERT INTO hu_active_home_timelines (hub_user_id, last_refreshed_at, last_accessed_at)
VALUES ('12345678-0024-0024-0024-000000000001', NOW(), NOW());

INSERT INTO hu_home_timelines (hub_user_id, post_id)
SELECT '12345678-0024-0024-0024-000000000001', id
FROM posts
WHERE author_id IN (
    '12345678-0024-0024-0024-000000000002',
    '12345678-0024-0024-0024-000000000003',
    '12345678-0024-0024-0024-000000000004'
);

-- For user5 (following users 6,7)
INSERT INTO hu_active_home_timelines (hub_user_id, last_refreshed_at, last_accessed_at)
VALUES ('12345678-0024-0024-0024-000000000005', NOW(), NOW());

INSERT INTO hu_home_timelines (hub_user_id, post_id)
SELECT '12345678-0024-0024-0024-000000000005', id
FROM posts
WHERE author_id IN (
    '12345678-0024-0024-0024-000000000006',
    '12345678-0024-0024-0024-000000000007'
);

-- For user11 (following users 12,13,14)
INSERT INTO hu_active_home_timelines (hub_user_id, last_refreshed_at, last_accessed_at)
VALUES ('12345678-0024-0024-0024-000000000011', NOW(), NOW());

INSERT INTO hu_home_timelines (hub_user_id, post_id)
SELECT '12345678-0024-0024-0024-000000000011', id
FROM posts
WHERE author_id IN (
    '12345678-0024-0024-0024-000000000012',
    '12345678-0024-0024-0024-000000000013',
    '12345678-0024-0024-0024-000000000014'
);

COMMIT;
