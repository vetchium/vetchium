-- Cleanup script for 0035-incognito-comments_test.go

-- Delete incognito post comment votes for test users
DELETE FROM incognito_post_comment_votes
WHERE comment_id IN (
    SELECT ipc.id FROM incognito_post_comments ipc
    JOIN incognito_posts ip ON ipc.incognito_post_id = ip.id
    WHERE ip.author_id::text LIKE '12345678-0035-%'
);

-- Delete incognito post votes for test users
DELETE FROM incognito_post_votes
WHERE incognito_post_id IN (
    SELECT id FROM incognito_posts WHERE author_id::text LIKE '12345678-0035-%'
);

-- Delete incognito post tags for test users
DELETE FROM incognito_post_tags
WHERE incognito_post_id IN (
    SELECT id FROM incognito_posts WHERE author_id::text LIKE '12345678-0035-%'
);

-- Delete incognito post comments for test users
DELETE FROM incognito_post_comments
WHERE incognito_post_id IN (
    SELECT id FROM incognito_posts WHERE author_id::text LIKE '12345678-0035-%'
);

-- Delete incognito posts for test users
DELETE FROM incognito_posts WHERE author_id::text LIKE '12345678-0035-%';

-- Delete hub_user_tfa_codes for test users (references hub_user_tokens)
DELETE FROM hub_user_tfa_codes
WHERE tfa_token IN (
    SELECT token FROM hub_user_tokens WHERE hub_user_id::text LIKE '12345678-0035-%'
);

-- Delete hub_user_tokens for test users (references hub_users)
DELETE FROM hub_user_tokens WHERE hub_user_id::text LIKE '12345678-0035-%';

-- Delete hub_users_official_emails for test users (references hub_users)
DELETE FROM hub_users_official_emails WHERE hub_user_id::text LIKE '12345678-0035-%';

-- Delete test users
DELETE FROM hub_users WHERE id::text LIKE '12345678-0035-%';