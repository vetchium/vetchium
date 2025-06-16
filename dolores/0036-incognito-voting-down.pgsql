-- Cleanup data for 0036-incognito-voting_test.go

-- Delete in proper order to respect foreign key constraints

-- Delete incognito post comment votes
DELETE FROM incognito_post_comment_votes 
WHERE comment_id IN (
    SELECT ipc.id 
    FROM incognito_post_comments ipc 
    JOIN incognito_posts ip ON ipc.incognito_post_id = ip.id 
    WHERE ip.author_id LIKE '12345678-0036-%'
);

-- Delete incognito post votes
DELETE FROM incognito_post_votes 
WHERE incognito_post_id IN (
    SELECT id FROM incognito_posts WHERE author_id LIKE '12345678-0036-%'
);

-- Delete incognito post comments
DELETE FROM incognito_post_comments 
WHERE incognito_post_id IN (
    SELECT id FROM incognito_posts WHERE author_id LIKE '12345678-0036-%'
);

-- Delete incognito post tags
DELETE FROM incognito_post_tags 
WHERE incognito_post_id IN (
    SELECT id FROM incognito_posts WHERE author_id LIKE '12345678-0036-%'
);

-- Delete incognito posts
DELETE FROM incognito_posts WHERE author_id LIKE '12345678-0036-%';

-- Delete hub users created for this test
DELETE FROM hub_users WHERE id LIKE '12345678-0036-%'; 