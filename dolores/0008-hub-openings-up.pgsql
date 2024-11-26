BEGIN;

-- Create hub user for testing with both resident country and city
INSERT INTO hub_users (
    id,
    full_name,
    handle,
    email,
    password_hash,
    state,
    resident_country_code,
    resident_city,
    created_at,
    updated_at
) VALUES (
    '12345678-0008-0008-0008-000000050001'::uuid,
    'Hub Opening Test User',
    'hub_opening_test',
    'hubopening@hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'IND',
    'Bangalore',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
);

-- Create additional test users with different combinations
INSERT INTO hub_users (
    id,
    full_name,
    handle,
    email,
    password_hash,
    state,
    resident_country_code,
    resident_city,
    created_at,
    updated_at
) VALUES 
-- User with only resident country
('12345678-0008-0008-0008-000000050002'::uuid,
 'India Only User',
 'india_only',
 'indiaonly@hub.example',
 '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
 'ACTIVE_HUB_USER',
 'IND',
 NULL,
 timezone('UTC'::text, now()),
 timezone('UTC'::text, now())),

-- User with different country/city combination
('12345678-0008-0008-0008-000000050003'::uuid,
 'US User',
 'us_user',
 'ususer@hub.example',
 '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
 'ACTIVE_HUB_USER',
 'USA',
 'New York',
 timezone('UTC'::text, now()),
 timezone('UTC'::text, now()));

-- Create test employers (5 companies)
INSERT INTO emails (email_key, email_from, email_to, email_cc, email_bcc, email_subject, email_html_body, email_text_body, email_state, created_at, processed_at)
VALUES
    ('12345678-0008-0008-0008-000000000011'::uuid, 'no-reply@vetchi.org', ARRAY['admin@hubopening1.example'], NULL, NULL, 'Welcome to Vetchi', 'Welcome HTML', 'Welcome Text', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now())),
    ('12345678-0008-0008-0008-000000000012'::uuid, 'no-reply@vetchi.org', ARRAY['admin@hubopening2.example'], NULL, NULL, 'Welcome to Vetchi', 'Welcome HTML', 'Welcome Text', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now())),
    ('12345678-0008-0008-0008-000000000013'::uuid, 'no-reply@vetchi.org', ARRAY['admin@hubopening3.example'], NULL, NULL, 'Welcome to Vetchi', 'Welcome HTML', 'Welcome Text', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now())),
    ('12345678-0008-0008-0008-000000000014'::uuid, 'no-reply@vetchi.org', ARRAY['admin@hubopening4.example'], NULL, NULL, 'Welcome to Vetchi', 'Welcome HTML', 'Welcome Text', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now())),
    ('12345678-0008-0008-0008-000000000015'::uuid, 'no-reply@vetchi.org', ARRAY['admin@hubopening5.example'], NULL, NULL, 'Welcome to Vetchi', 'Welcome HTML', 'Welcome Text', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now()));

INSERT INTO employers (id, client_id_type, employer_state, company_name, onboard_admin_email, onboard_secret_token, token_valid_till, onboard_email_id, created_at)
VALUES
    ('12345678-0008-0008-0008-000000000201'::uuid, 'DOMAIN', 'ONBOARDED', 'hubopening1.example', 'admin@hubopening1.example', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0008-0008-0008-000000000011'::uuid, timezone('UTC'::text, now())),
    ('12345678-0008-0008-0008-000000000202'::uuid, 'DOMAIN', 'ONBOARDED', 'hubopening2.example', 'admin@hubopening2.example', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0008-0008-0008-000000000012'::uuid, timezone('UTC'::text, now())),
    ('12345678-0008-0008-0008-000000000203'::uuid, 'DOMAIN', 'ONBOARDED', 'hubopening3.example', 'admin@hubopening3.example', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0008-0008-0008-000000000013'::uuid, timezone('UTC'::text, now())),
    ('12345678-0008-0008-0008-000000000204'::uuid, 'DOMAIN', 'ONBOARDED', 'hubopening4.example', 'admin@hubopening4.example', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0008-0008-0008-000000000014'::uuid, timezone('UTC'::text, now())),
    ('12345678-0008-0008-0008-000000000205'::uuid, 'DOMAIN', 'ONBOARDED', 'hubopening5.example', 'admin@hubopening5.example', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0008-0008-0008-000000000015'::uuid, timezone('UTC'::text, now()));

-- Create domains for each employer
INSERT INTO domains (id, domain_name, domain_state, employer_id, created_at)
VALUES
    -- Company 1 domains
    ('12345678-0008-0008-0008-000000003001'::uuid, 'hubopening1.example', 'VERIFIED', '12345678-0008-0008-0008-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0008-0008-0008-000000003002'::uuid, 'hubopening1-alt.example', 'VERIFIED', '12345678-0008-0008-0008-000000000201'::uuid, timezone('UTC'::text, now())),
    
    -- Company 2 domains
    ('12345678-0008-0008-0008-000000003003'::uuid, 'hubopening2.example', 'VERIFIED', '12345678-0008-0008-0008-000000000202'::uuid, timezone('UTC'::text, now())),
    
    -- Company 3 domains
    ('12345678-0008-0008-0008-000000003004'::uuid, 'hubopening3.example', 'VERIFIED', '12345678-0008-0008-0008-000000000203'::uuid, timezone('UTC'::text, now())),
    ('12345678-0008-0008-0008-000000003005'::uuid, 'hubopening3-alt.example', 'VERIFIED', '12345678-0008-0008-0008-000000000203'::uuid, timezone('UTC'::text, now())),
    ('12345678-0008-0008-0008-000000003006'::uuid, 'hubopening3-alt2.example', 'VERIFIED', '12345678-0008-0008-0008-000000000203'::uuid, timezone('UTC'::text, now())),
    
    -- Company 4 domains
    ('12345678-0008-0008-0008-000000003007'::uuid, 'hubopening4.example', 'VERIFIED', '12345678-0008-0008-0008-000000000204'::uuid, timezone('UTC'::text, now())),
    ('12345678-0008-0008-0008-000000003008'::uuid, 'hubopening4-alt.example', 'VERIFIED', '12345678-0008-0008-0008-000000000204'::uuid, timezone('UTC'::text, now())),
    
    -- Company 5 domains
    ('12345678-0008-0008-0008-000000003009'::uuid, 'hubopening5.example', 'VERIFIED', '12345678-0008-0008-0008-000000000205'::uuid, timezone('UTC'::text, now()));

-- Create org users for each employer
INSERT INTO org_users (id, email, name, password_hash, org_user_roles, org_user_state, employer_id, created_at)
SELECT 
    gen_random_uuid(),
    'admin@' || domain_name,
    'Admin User',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    ARRAY['ADMIN']::org_user_roles[],
    'ACTIVE_ORG_USER',
    employer_id,
    timezone('UTC'::text, now())
FROM domains;

-- Create locations for each employer
INSERT INTO locations (id, title, country_code, postal_address, postal_code, city_aka, location_state, employer_id, created_at)
VALUES
    -- Company 1 locations
    ('12345678-0008-0008-0008-000000060001'::uuid, 'Bangalore Office', 'IND', '123 MG Road', '560001', ARRAY['Bengaluru', 'Bangalore'], 'ACTIVE_LOCATION', '12345678-0008-0008-0008-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0008-0008-0008-000000060002'::uuid, 'Chennai Office', 'IND', '456 Anna Salai', '600002', ARRAY['Madras'], 'ACTIVE_LOCATION', '12345678-0008-0008-0008-000000000201'::uuid, timezone('UTC'::text, now())),
    
    -- Company 2 locations
    ('12345678-0008-0008-0008-000000060003'::uuid, 'New York Office', 'USA', '789 Broadway', '10013', ARRAY['NYC'], 'ACTIVE_LOCATION', '12345678-0008-0008-0008-000000000202'::uuid, timezone('UTC'::text, now())),
    ('12345678-0008-0008-0008-000000060004'::uuid, 'San Francisco Office', 'USA', '123 Market St', '94105', ARRAY['SF'], 'ACTIVE_LOCATION', '12345678-0008-0008-0008-000000000202'::uuid, timezone('UTC'::text, now())),
    
    -- Company 3 locations
    ('12345678-0008-0008-0008-000000060005'::uuid, 'London Office', 'GBR', '456 Oxford St', 'W1D 1BS', ARRAY['Central London'], 'ACTIVE_LOCATION', '12345678-0008-0008-0008-000000000203'::uuid, timezone('UTC'::text, now())),
    
    -- Company 4 locations
    ('12345678-0008-0008-0008-000000060006'::uuid, 'Singapore Office', 'SGP', '789 Orchard Rd', '238839', ARRAY['SG'], 'ACTIVE_LOCATION', '12345678-0008-0008-0008-000000000204'::uuid, timezone('UTC'::text, now())),
    ('12345678-0008-0008-0008-000000060007'::uuid, 'Tokyo Office', 'JPN', '123 Shibuya', '150-0002', ARRAY['TYO'], 'ACTIVE_LOCATION', '12345678-0008-0008-0008-000000000204'::uuid, timezone('UTC'::text, now())),
    
    -- Company 5 locations
    ('12345678-0008-0008-0008-000000060008'::uuid, 'Sydney Office', 'AUS', '456 George St', '2000', ARRAY['SYD'], 'ACTIVE_LOCATION', '12345678-0008-0008-0008-000000000205'::uuid, timezone('UTC'::text, now())),
    ('12345678-0008-0008-0008-000000060009'::uuid, 'Melbourne Office', 'AUS', '789 Collins St', '3000', ARRAY['MEL'], 'ACTIVE_LOCATION', '12345678-0008-0008-0008-000000000205'::uuid, timezone('UTC'::text, now()));

-- Create cost centers for each employer
INSERT INTO org_cost_centers (id, cost_center_name, cost_center_state, notes, employer_id, created_at)
SELECT 
    gen_random_uuid(),
    'Engineering',
    'ACTIVE_CC',
    'Engineering department',
    id,
    timezone('UTC'::text, now())
FROM employers;

-- Create openings for each employer
-- Company 1: Focus on software engineering roles
INSERT INTO openings (
    employer_id, 
    id,
    title,
    positions,
    jd,
    recruiter,
    hiring_manager,
    cost_center_id,
    opening_type,
    yoe_min,
    yoe_max,
    min_education_level,
    salary_min,
    salary_max,
    salary_currency,
    remote_country_codes,
    remote_timezones,
    state,
    created_at,
    last_updated_at
)
SELECT
    '12345678-0008-0008-0008-000000000201'::uuid,
    '2024-Mar-01-' || LPAD(CAST(generate_series AS text), 3, '0'),
    CASE (generate_series % 5)
        WHEN 0 THEN 'Senior Software Engineer'
        WHEN 1 THEN 'Full Stack Developer'
        WHEN 2 THEN 'Backend Engineer'
        WHEN 3 THEN 'Frontend Developer'
        WHEN 4 THEN 'DevOps Engineer'
    END,
    CASE (generate_series % 3)
        WHEN 0 THEN 1
        WHEN 1 THEN 2
        WHEN 2 THEN 3
    END,
    'Looking for talented engineers...',
    (SELECT id FROM org_users WHERE employer_id = '12345678-0008-0008-0008-000000000201'::uuid LIMIT 1),
    (SELECT id FROM org_users WHERE employer_id = '12345678-0008-0008-0008-000000000201'::uuid LIMIT 1),
    (SELECT id FROM org_cost_centers WHERE employer_id = '12345678-0008-0008-0008-000000000201'::uuid LIMIT 1),
    'FULL_TIME_OPENING',
    CASE (generate_series % 3)
        WHEN 0 THEN 0
        WHEN 1 THEN 2
        WHEN 2 THEN 5
    END,
    CASE (generate_series % 3)
        WHEN 0 THEN 3
        WHEN 1 THEN 5
        WHEN 2 THEN 8
    END,
    'BACHELOR_EDUCATION',
    50000,
    100000,
    'USD',
    ARRAY['IND', 'USA'],
    ARRAY['IST Indian Standard Time GMT+0530'],
    'ACTIVE_OPENING_STATE',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
FROM generate_series(1, 10);

-- Company 2: Focus on data science roles
INSERT INTO openings (
    employer_id, 
    id,
    title,
    positions,
    jd,
    recruiter,
    hiring_manager,
    cost_center_id,
    opening_type,
    yoe_min,
    yoe_max,
    min_education_level,
    salary_min,
    salary_max,
    salary_currency,
    remote_country_codes,
    remote_timezones,
    state,
    created_at,
    last_updated_at
)
SELECT
    '12345678-0008-0008-0008-000000000202'::uuid,
    '2024-Mar-02-' || LPAD(CAST(generate_series AS text), 3, '0'),
    CASE (generate_series % 5)
        WHEN 0 THEN 'Data Scientist'
        WHEN 1 THEN 'Machine Learning Engineer'
        WHEN 2 THEN 'AI Researcher'
        WHEN 3 THEN 'Data Analyst'
        WHEN 4 THEN 'Research Scientist'
    END,
    CASE (generate_series % 3)
        WHEN 0 THEN 1
        WHEN 1 THEN 2
        WHEN 2 THEN 3
    END,
    'Join our data science team...',
    (SELECT id FROM org_users WHERE employer_id = '12345678-0008-0008-0008-000000000202'::uuid LIMIT 1),
    (SELECT id FROM org_users WHERE employer_id = '12345678-0008-0008-0008-000000000202'::uuid LIMIT 1),
    (SELECT id FROM org_cost_centers WHERE employer_id = '12345678-0008-0008-0008-000000000202'::uuid LIMIT 1),
    'FULL_TIME_OPENING',
    CASE (generate_series % 3)
        WHEN 0 THEN 2
        WHEN 1 THEN 4
        WHEN 2 THEN 6
    END,
    CASE (generate_series % 3)
        WHEN 0 THEN 5
        WHEN 1 THEN 7
        WHEN 2 THEN 10
    END,
    'MASTER_EDUCATION',
    80000,
    150000,
    'USD',
    ARRAY['USA'],
    ARRAY['PST Pacific Standard Time (North America) GMT-0800'],
    'ACTIVE_OPENING_STATE',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
FROM generate_series(1, 10);

-- Company 3: Focus on product/design roles
INSERT INTO openings (
    employer_id,
    id,
    title,
    positions,
    jd,
    recruiter,
    hiring_manager,
    cost_center_id,
    opening_type,
    yoe_min,
    yoe_max,
    min_education_level,
    salary_min,
    salary_max,
    salary_currency,
    remote_country_codes,
    remote_timezones,
    state,
    created_at,
    last_updated_at
)
SELECT
    '12345678-0008-0008-0008-000000000203'::uuid,
    '2024-Mar-03-' || LPAD(CAST(generate_series AS text), 3, '0'),
    CASE (generate_series % 5)
        WHEN 0 THEN 'Product Manager'
        WHEN 1 THEN 'UX Designer'
        WHEN 2 THEN 'Product Designer'
        WHEN 3 THEN 'UI Developer'
        WHEN 4 THEN 'Design Systems Engineer'
    END,
    CASE (generate_series % 3)
        WHEN 0 THEN 1
        WHEN 1 THEN 2
        WHEN 2 THEN 3
    END,
    'Join our product team...',
    (SELECT id FROM org_users WHERE employer_id = '12345678-0008-0008-0008-000000000203'::uuid LIMIT 1),
    (SELECT id FROM org_users WHERE employer_id = '12345678-0008-0008-0008-000000000203'::uuid LIMIT 1),
    (SELECT id FROM org_cost_centers WHERE employer_id = '12345678-0008-0008-0008-000000000203'::uuid LIMIT 1),
    'FULL_TIME_OPENING',
    CASE (generate_series % 3)
        WHEN 0 THEN 1
        WHEN 1 THEN 3
        WHEN 2 THEN 5
    END,
    CASE (generate_series % 3)
        WHEN 0 THEN 4
        WHEN 1 THEN 6
        WHEN 2 THEN 8
    END,
    'BACHELOR_EDUCATION',
    70000,
    130000,
    'USD',
    ARRAY['GBR'],
    ARRAY['GMT Greenwich Mean Time GMT+0000'],
    'ACTIVE_OPENING_STATE',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
FROM generate_series(1, 10);

-- Link openings to locations
INSERT INTO opening_locations (employer_id, opening_id, location_id)
SELECT 
    o.employer_id,
    o.id,
    l.id
FROM openings o
CROSS JOIN locations l
WHERE o.employer_id = l.employer_id;

COMMIT;
