BEGIN;

-- Clean up work history entries for both employers and hub users
DELETE FROM work_history
WHERE employer_id IN (
    '12345678-0017-0017-0017-000000000003',
    '12345678-0017-0017-0017-000000000004'
) OR hub_user_id IN (
    '12345678-0017-0017-0017-000000000001',
    '12345678-0017-0017-0017-000000000002'
);

-- Clean up any opening related data
DELETE FROM opening_tag_mappings
WHERE employer_id IN (
    '12345678-0017-0017-0017-000000000003',
    '12345678-0017-0017-0017-000000000004'
);

DELETE FROM interview_interviewers
WHERE employer_id IN (
    '12345678-0017-0017-0017-000000000003',
    '12345678-0017-0017-0017-000000000004'
);

DELETE FROM interviews
WHERE employer_id IN (
    '12345678-0017-0017-0017-000000000003',
    '12345678-0017-0017-0017-000000000004'
);

DELETE FROM candidacy_comments
WHERE employer_id IN (
    '12345678-0017-0017-0017-000000000003',
    '12345678-0017-0017-0017-000000000004'
);

DELETE FROM candidacies
WHERE employer_id IN (
    '12345678-0017-0017-0017-000000000003',
    '12345678-0017-0017-0017-000000000004'
);

DELETE FROM applications
WHERE employer_id IN (
    '12345678-0017-0017-0017-000000000003',
    '12345678-0017-0017-0017-000000000004'
);

DELETE FROM opening_locations
WHERE employer_id IN (
    '12345678-0017-0017-0017-000000000003',
    '12345678-0017-0017-0017-000000000004'
);

DELETE FROM opening_hiring_team
WHERE employer_id IN (
    '12345678-0017-0017-0017-000000000003',
    '12345678-0017-0017-0017-000000000004'
);

DELETE FROM opening_watchers
WHERE employer_id IN (
    '12345678-0017-0017-0017-000000000003',
    '12345678-0017-0017-0017-000000000004'
);

DELETE FROM openings
WHERE employer_id IN (
    '12345678-0017-0017-0017-000000000003',
    '12345678-0017-0017-0017-000000000004'
);

-- Clean up locations
DELETE FROM locations
WHERE employer_id IN (
    '12345678-0017-0017-0017-000000000003',
    '12345678-0017-0017-0017-000000000004'
);

-- Clean up cost centers
DELETE FROM org_cost_centers
WHERE employer_id IN (
    '12345678-0017-0017-0017-000000000003',
    '12345678-0017-0017-0017-000000000004'
);

-- Clean up org users and related tables
DELETE FROM org_user_tfa_codes
WHERE tfa_token IN (
    SELECT token FROM org_user_tokens
    WHERE org_user_id IN (
        SELECT id FROM org_users
        WHERE employer_id IN (
            '12345678-0017-0017-0017-000000000003',
            '12345678-0017-0017-0017-000000000004'
        )
    )
);

DELETE FROM org_user_tokens
WHERE org_user_id IN (
    SELECT id FROM org_users
    WHERE employer_id IN (
        '12345678-0017-0017-0017-000000000003',
        '12345678-0017-0017-0017-000000000004'
    )
);

DELETE FROM org_user_invites
WHERE org_user_id IN (
    SELECT id FROM org_users
    WHERE employer_id IN (
        '12345678-0017-0017-0017-000000000003',
        '12345678-0017-0017-0017-000000000004'
    )
);

DELETE FROM org_users
WHERE employer_id IN (
    '12345678-0017-0017-0017-000000000003',
    '12345678-0017-0017-0017-000000000004'
);

-- Clean up primary domains
DELETE FROM employer_primary_domains
WHERE employer_id IN (
    '12345678-0017-0017-0017-000000000003',
    '12345678-0017-0017-0017-000000000004'
);

-- Clean up domains
DELETE FROM domains
WHERE employer_id IN (
    '12345678-0017-0017-0017-000000000003',
    '12345678-0017-0017-0017-000000000004'
);

-- Clean up employers
DELETE FROM employers
WHERE id IN (
    '12345678-0017-0017-0017-000000000003',
    '12345678-0017-0017-0017-000000000004'
);

-- Clean up hub user tokens and related tables first
DELETE FROM hub_user_tfa_codes
WHERE tfa_token IN (
    SELECT token FROM hub_user_tokens
    WHERE hub_user_id IN (
        '12345678-0017-0017-0017-000000000001',
        '12345678-0017-0017-0017-000000000002'
    )
);

DELETE FROM hub_user_tokens
WHERE hub_user_id IN (
    '12345678-0017-0017-0017-000000000001',
    '12345678-0017-0017-0017-000000000002'
);

-- Finally clean up hub users
DELETE FROM hub_users
WHERE id IN (
    '12345678-0017-0017-0017-000000000001',
    '12345678-0017-0017-0017-000000000002'
);

COMMIT;
