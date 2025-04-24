BEGIN;

-- First clean up timeline data since it references posts
DELETE FROM hu_home_timelines
WHERE hub_user_id::text LIKE '12345678-0025-0025-0025-%'
   OR post_id IN (SELECT id FROM posts WHERE author_id::text LIKE '12345678-0025-0025-0025-%');

DELETE FROM hu_active_home_timelines
WHERE hub_user_id::text LIKE '12345678-0025-0025-0025-%';

-- Clean up follow relationships
DELETE FROM following_relationships
WHERE consuming_hub_user_id::text LIKE '12345678-0025-0025-0025-%'
   OR producing_hub_user_id::text LIKE '12345678-0025-0025-0025-%';

-- Clean up post votes
DELETE FROM post_votes
WHERE user_id::text LIKE '12345678-0025-0025-0025-%'
   OR post_id IN (SELECT id FROM posts WHERE author_id::text LIKE '12345678-0025-0025-0025-%');

-- Clean up post tags
DELETE FROM post_tags
WHERE post_id IN (SELECT id FROM posts WHERE author_id::text LIKE '12345678-0025-0025-0025-%');

-- Clean up posts
DELETE FROM posts
WHERE author_id::text LIKE '12345678-0025-0025-0025-%';

-- Clean up hub user tokens and related tables
DELETE FROM hub_user_tfa_codes
WHERE tfa_token IN (
    SELECT token FROM hub_user_tokens
    WHERE hub_user_id::text LIKE '12345678-0025-0025-0025-%'
);

DELETE FROM hub_user_tokens
WHERE hub_user_id::text LIKE '12345678-0025-0025-0025-%';

-- Finally clean up users
DELETE FROM hub_users
WHERE id::text LIKE '12345678-0025-0025-0025-%';

COMMIT;
