BEGIN;

-- First clean up opening_tag_mappings
DELETE FROM opening_tag_mappings WHERE employer_id = '12345678-0015-0015-0015-000000000201'::uuid;

-- Then clean up opening_locations
DELETE FROM opening_locations WHERE employer_id = '12345678-0015-0015-0015-000000000201'::uuid;

-- Then clean up openings
DELETE FROM openings WHERE employer_id = '12345678-0015-0015-0015-000000000201'::uuid;

-- Clean up hub user
DELETE FROM hub_users WHERE id = '12345678-0015-0015-0015-000000050002'::uuid;

-- Clean up location
DELETE FROM locations WHERE id = '12345678-0015-0015-0015-000000060001'::uuid;

-- Clean up cost center
DELETE FROM org_cost_centers WHERE id = '12345678-0015-0015-0015-000000050001'::uuid;

-- Clean up org user tokens first
DELETE FROM org_user_tokens WHERE org_user_id = '12345678-0015-0015-0015-000000040001'::uuid;

-- Clean up org user
DELETE FROM org_users WHERE id = '12345678-0015-0015-0015-000000040001'::uuid;

-- Clean up domain mapping
DELETE FROM employer_primary_domains WHERE domain_id = '12345678-0015-0015-0015-000000003001'::uuid;

-- Clean up domain
DELETE FROM domains WHERE id = '12345678-0015-0015-0015-000000003001'::uuid;

-- Clean up employer
DELETE FROM employers WHERE id = '12345678-0015-0015-0015-000000000201'::uuid;

-- Clean up email
DELETE FROM emails WHERE email_key = '12345678-0015-0015-0015-000000000011'::uuid;

COMMIT;
