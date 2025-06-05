BEGIN;

-- Cleanup script for incognito posts testing
-- Remove in reverse order of dependencies

-- Remove incognito post votes (if any were created during tests)
DELETE FROM incognito_post_votes 
WHERE incognito_post_id IN (
    SELECT id FROM incognito_posts 
    WHERE author_id IN (
        SELECT id FROM hub_users WHERE id::text LIKE '12345678-0032-0032-0032-%'
    )
);

-- Remove incognito post comment votes (if any were created during tests)
DELETE FROM incognito_post_comment_votes 
WHERE comment_id IN (
    SELECT c.id FROM incognito_post_comments c
    JOIN incognito_posts p ON c.incognito_post_id = p.id
    WHERE p.author_id IN (
        SELECT id FROM hub_users WHERE id::text LIKE '12345678-0032-0032-0032-%'
    )
);

-- Remove incognito post comments (if any were created during tests)
DELETE FROM incognito_post_comments 
WHERE incognito_post_id IN (
    SELECT id FROM incognito_posts 
    WHERE author_id IN (
        SELECT id FROM hub_users WHERE id::text LIKE '12345678-0032-0032-0032-%'
    )
);

-- Remove incognito post tags (if any were created during tests)
DELETE FROM incognito_post_tags 
WHERE incognito_post_id IN (
    SELECT id FROM incognito_posts 
    WHERE author_id IN (
        SELECT id FROM hub_users WHERE id::text LIKE '12345678-0032-0032-0032-%'
    )
);

-- Remove incognito posts (if any were created during tests)
DELETE FROM incognito_posts 
WHERE author_id IN (
    SELECT id FROM hub_users WHERE id::text LIKE '12345678-0032-0032-0032-%'
);

-- Clean up any hub user tokens that might reference our test users
DELETE FROM hub_user_tfa_codes
WHERE tfa_token IN (
    SELECT token FROM hub_user_tokens
    WHERE hub_user_id IN (
        SELECT id FROM hub_users WHERE id::text LIKE '12345678-0032-0032-0032-%'
    )
);

DELETE FROM hub_user_tokens 
WHERE hub_user_id IN (
    SELECT id FROM hub_users WHERE id::text LIKE '12345678-0032-0032-0032-%'
);

-- Clean up any hub user invites for our test domains
DELETE FROM hub_user_invites 
WHERE email LIKE '%@test0032.com' OR email LIKE '%@company0032.com';

-- Clean up any emails sent to our test domains
DELETE FROM emails 
WHERE EXISTS (
    SELECT 1 FROM UNNEST(email_to) AS recipient 
    WHERE recipient LIKE '%@test0032.com' OR recipient LIKE '%@company0032.com'
);

-- Remove test hub users
DELETE FROM hub_users WHERE id::text LIKE '12345678-0032-0032-0032-%';

-- Remove test approved domains - MUST be last to avoid FK constraint issues
DELETE FROM hub_user_signup_approved_domains 
WHERE domain_name IN ('test0032.com', 'company0032.com');

COMMIT; 