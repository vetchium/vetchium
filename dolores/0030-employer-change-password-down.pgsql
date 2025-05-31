BEGIN;

-- Clean up test data in correct order (foreign key dependencies)
DELETE FROM public.org_user_tokens WHERE org_user_id IN (
    SELECT id FROM public.org_users WHERE employer_id = '12345678-0030-0030-0030-000000000001'::uuid
);
DELETE FROM public.employer_primary_domains WHERE employer_id = '12345678-0030-0030-0030-000000000001'::uuid;
DELETE FROM public.org_users WHERE employer_id = '12345678-0030-0030-0030-000000000001'::uuid;
DELETE FROM public.domains WHERE employer_id = '12345678-0030-0030-0030-000000000001'::uuid;
DELETE FROM public.employers WHERE id = '12345678-0030-0030-0030-000000000001'::uuid;
DELETE FROM public.emails WHERE email_key = '12345678-0030-0030-0030-000000000011'::uuid;

COMMIT; 