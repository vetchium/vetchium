BEGIN;

DELETE FROM candidacy_comments
WHERE candidacy_id IN (
    SELECT id FROM candidacies 
    WHERE employer_id = '12345678-0011-0011-0011-000000000201'::uuid
);

DELETE FROM candidacies
WHERE employer_id = '12345678-0011-0011-0011-000000000201'::uuid;

DELETE FROM applications
WHERE employer_id = '12345678-0011-0011-0011-000000000201'::uuid;

DELETE FROM opening_watchers
WHERE opening_id IN (
    SELECT id FROM openings 
    WHERE employer_id = '12345678-0011-0011-0011-000000000201'::uuid
);

DELETE FROM org_cost_centers
WHERE employer_id = '12345678-0011-0011-0011-000000000201'::uuid;

DELETE FROM openings 
WHERE employer_id = '12345678-0011-0011-0011-000000000201'::uuid;

DELETE FROM org_user_tokens
WHERE org_user_id IN (
    SELECT id FROM org_users 
    WHERE employer_id = '12345678-0011-0011-0011-000000000201'::uuid
);


DELETE FROM org_users
WHERE employer_id = '12345678-0011-0011-0011-000000000201'::uuid
AND email IN ('admin@candidacy-comments.example', 'hiringmanager@candidacy-comments.example', 'recruiter@candidacy-comments.example', 'watcher@candidacy-comments.example', 'regular@candidacy-comments.example');

DELETE FROM domains
WHERE employer_id = '12345678-0011-0011-0011-000000000201'::uuid;

DELETE FROM employers
WHERE id = '12345678-0011-0011-0011-000000000201'::uuid;

DELETE FROM emails
WHERE email_key = '12345678-0011-0011-0011-000000000011'::uuid;

DELETE FROM hub_user_tokens
WHERE hub_user_id IN (
    SELECT id FROM hub_users
    WHERE email IN ('0011-active@hub.example')
);

DELETE FROM hub_users
WHERE email IN ('0011-active@hub.example');

COMMIT;
