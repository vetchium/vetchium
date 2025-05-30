BEGIN;
DELETE FROM org_user_invites
WHERE org_user_id IN (
    SELECT
        id
    FROM
        org_users
    WHERE
        employer_id = '12345678-0004-0004-0004-000000000201'::uuid
);
DELETE FROM org_user_tokens
WHERE org_user_id IN (
    SELECT
        id
    FROM
        org_users
    WHERE
        employer_id = '12345678-0004-0004-0004-000000000201'::uuid);
DELETE FROM public.org_users
WHERE employer_id = '12345678-0004-0004-0004-000000000201'::uuid;
DELETE FROM public.domains
WHERE employer_id = '12345678-0004-0004-0004-000000000201'::uuid;
DELETE FROM public.employers
WHERE id = '12345678-0004-0004-0004-000000000201'::uuid;
DELETE FROM public.emails
WHERE email_key = '12345678-0004-0004-0004-000000000011'::uuid;
COMMIT;
