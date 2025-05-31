BEGIN;

-- Clean up test data
DELETE FROM public.org_users WHERE employer_id = '12345678-0030-0030-0030-000000000001'::uuid;
DELETE FROM public.employers WHERE id = '12345678-0030-0030-0030-000000000001'::uuid;

COMMIT; 