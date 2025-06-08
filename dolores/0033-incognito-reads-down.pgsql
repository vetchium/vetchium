BEGIN;

-- Cleanup script for incognito reads testing
-- Remove in reverse order of dependencies

-- Remove incognito post comment votes (created via APIs during tests)
DELETE FROM incognito_post_comment_votes 
WHERE comment_id IN (
    SELECT c.id FROM incognito_post_comments c
    JOIN incognito_posts p ON c.incognito_post_id = p.id
    WHERE p.author_id IN (
        SELECT id FROM hub_users WHERE id::text LIKE '12345678-0033-0033-0033-%'
    )
);

-- Remove incognito post votes (created via APIs during tests)
DELETE FROM incognito_post_votes 
WHERE incognito_post_id IN (
    SELECT id FROM incognito_posts 
    WHERE author_id IN (
        SELECT id FROM hub_users WHERE id::text LIKE '12345678-0033-0033-0033-%'
    )
);

-- Remove incognito post comments (created via APIs during tests)
DELETE FROM incognito_post_comments 
WHERE incognito_post_id IN (
    SELECT id FROM incognito_posts 
    WHERE author_id IN (
        SELECT id FROM hub_users WHERE id::text LIKE '12345678-0033-0033-0033-%'
    )
);

-- Remove incognito post tags (created via APIs during tests)
DELETE FROM incognito_post_tags 
WHERE incognito_post_id IN (
    SELECT id FROM incognito_posts 
    WHERE author_id IN (
        SELECT id FROM hub_users WHERE id::text LIKE '12345678-0033-0033-0033-%'
    )
);

-- Remove incognito posts (created via APIs during tests)
DELETE FROM incognito_posts 
WHERE author_id IN (
    SELECT id FROM hub_users WHERE id::text LIKE '12345678-0033-0033-0033-%'
);

-- Clean up hub user tokens
DELETE FROM hub_user_tfa_codes
WHERE tfa_token IN (
    SELECT token FROM hub_user_tokens
    WHERE hub_user_id IN (
        SELECT id FROM hub_users WHERE id::text LIKE '12345678-0033-0033-0033-%'
    )
);

DELETE FROM hub_user_tokens 
WHERE hub_user_id IN (
    SELECT id FROM hub_users WHERE id::text LIKE '12345678-0033-0033-0033-%'
);

-- Clean up hub user invites for test domains
DELETE FROM hub_user_invites 
WHERE email LIKE '%@test0033.com' OR email LIKE '%@company0033.com';

-- Clean up emails for test domains
DELETE FROM emails 
WHERE EXISTS (
    SELECT 1 FROM UNNEST(email_to) AS recipient 
    WHERE recipient LIKE '%@test0033.com' OR recipient LIKE '%@company0033.com'
);

-- Remove test hub users
DELETE FROM hub_users WHERE id::text LIKE '12345678-0033-0033-0033-%';

-- Remove test approved domains - MUST be last to avoid FK constraint issues
DELETE FROM hub_user_signup_approved_domains 
WHERE domain_name IN ('test0033.com', 'company0033.com');

COMMIT;
