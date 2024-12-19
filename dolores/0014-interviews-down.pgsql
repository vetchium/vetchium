BEGIN;
DELETE FROM interview_assessments
WHERE employer_id = '12345678-0014-0014-000000000201'::uuid;

DELETE FROM interview_interviewers
WHERE employer_id = '12345678-0014-0014-000000000201'::uuid;

DELETE FROM interviews
WHERE employer_id = '12345678-0014-0014-000000000201'::uuid;

DELETE FROM candidacy_comments
WHERE employer_id = '12345678-0014-0014-000000000201'::uuid;

DELETE FROM candidacies
WHERE employer_id = '12345678-0014-0014-000000000201'::uuid;

DELETE FROM applications
WHERE employer_id = '12345678-0014-0014-000000000201'::uuid;

DELETE FROM opening_locations
WHERE employer_id = '12345678-0014-0014-000000000201'::uuid;

DELETE FROM openings
WHERE employer_id = '12345678-0014-0014-000000000201'::uuid;

DELETE FROM locations
WHERE employer_id = '12345678-0014-0014-000000000201'::uuid;

DELETE FROM org_cost_centers
WHERE employer_id = '12345678-0014-0014-000000000201'::uuid;

DELETE FROM org_user_tokens
WHERE org_user_id IN (
    SELECT id FROM org_users 
    WHERE employer_id = '12345678-0014-0014-000000000201'::uuid
);

DELETE FROM org_users
WHERE employer_id = '12345678-0014-0014-000000000201'::uuid;

DELETE FROM domains
WHERE employer_id = '12345678-0014-0014-000000000201'::uuid;

DELETE FROM employers
WHERE id = '12345678-0014-0014-000000000201'::uuid;

DELETE FROM emails
WHERE email_key = '12345678-0014-0014-000000000011'::uuid;

DELETE FROM hub_user_tokens
WHERE hub_user_id IN (
    SELECT id FROM hub_users
    WHERE email = 'interview@hub.example'
);

DELETE FROM hub_users
WHERE email = 'interview@hub.example';

COMMIT; 