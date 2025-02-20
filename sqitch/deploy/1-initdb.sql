-- Deploy vetchi:1-initdb to pg

BEGIN;

CREATE TYPE hub_user_states AS ENUM ('ACTIVE_HUB_USER', 'DISABLED_HUB_USER', 'DELETED_HUB_USER');
CREATE TABLE IF NOT EXISTS hub_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name TEXT NOT NULL,
    handle TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    state hub_user_states NOT NULL,
    resident_country_code TEXT NOT NULL,
    resident_city TEXT,
    preferred_language TEXT NOT NULL,
    short_bio TEXT NOT NULL,
    long_bio TEXT NOT NULL,
    profile_picture_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TYPE hub_user_token_types AS ENUM (
    -- Sent as response to the TFA API
    'HUB_USER_SESSION',
    'HUB_USER_LTS',

    -- Sent as response to the Login API
    'HUB_USER_TFA_TOKEN',

    -- Sent as response to the Reset Password API
    'HUB_USER_RESET_PASSWORD_TOKEN'
);

CREATE TABLE hub_user_tokens (
    token TEXT CONSTRAINT hub_user_tokens_pkey PRIMARY KEY,
    hub_user_id UUID REFERENCES hub_users(id) NOT NULL,
    token_type hub_user_token_types NOT NULL,
    token_valid_till TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TABLE hub_user_tfa_codes (
    code TEXT NOT NULL,
    tfa_token TEXT NOT NULL REFERENCES hub_user_tokens(token) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TYPE email_states AS ENUM ('PENDING', 'PROCESSED');
CREATE TABLE emails(
	email_key UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	email_from TEXT NOT NULL,
	email_to TEXT ARRAY NOT NULL,
	email_cc TEXT ARRAY,
	email_bcc TEXT ARRAY,
	email_subject TEXT NOT NULL,
	email_html_body TEXT NOT NULL,
	email_text_body TEXT NOT NULL,
	email_state email_states NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),
	processed_at TIMESTAMP WITH TIME ZONE
);

---

CREATE TYPE client_id_types AS ENUM ('DOMAIN');
CREATE TYPE employer_states AS ENUM (
    'ONBOARD_PENDING',
    'ONBOARDED',
    'DEBOARDED',
    'HUB_ADDED_EMPLOYER'
);
CREATE TABLE employers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id_type client_id_types NOT NULL,
    employer_state employer_states NOT NULL,
    company_name TEXT NOT NULL,

    onboard_admin_email TEXT NOT NULL,

    -- TODO: Perhaps we can move this to org_user_tokens ?
    onboard_secret_token TEXT,
    token_valid_till TIMESTAMP WITH TIME ZONE,

    --- Despite its name, it should not be confused with an email address.
    --- This is the rowid in the 'emails' table for the welcome email sent.
    onboard_email_id UUID REFERENCES emails(email_key),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

-- Create a function to check if employer has required records
CREATE OR REPLACE FUNCTION check_employer_required_records()
RETURNS TRIGGER AS $$
BEGIN
    -- Only check for ONBOARDED employers
    IF NEW.employer_state = 'ONBOARDED' THEN
        IF NOT EXISTS (
            SELECT 1 FROM domains
            WHERE employer_id = NEW.id
        ) THEN
            RAISE EXCEPTION 'Onboarded employer must have at least one domain record';
        END IF;

        IF NOT EXISTS (
            SELECT 1 FROM employer_primary_domains
            WHERE employer_id = NEW.id
        ) THEN
            RAISE EXCEPTION 'Onboarded employer must have a primary domain record';
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to enforce the constraint
CREATE CONSTRAINT TRIGGER enforce_employer_required_records
    AFTER INSERT OR UPDATE ON employers
    DEFERRABLE INITIALLY DEFERRED
    FOR EACH ROW
    EXECUTE FUNCTION check_employer_required_records();

---

CREATE TYPE domain_states AS ENUM (
    'UNVERIFIED',
    'VERIFIED'
);
CREATE TABLE domains (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    domain_name TEXT NOT NULL,
    CONSTRAINT uniq_domain_name UNIQUE (domain_name),

    domain_state domain_states NOT NULL,

    employer_id UUID REFERENCES employers(id),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),
    CONSTRAINT uniq_employer_domain_id UNIQUE (employer_id, id)
);

CREATE TABLE hub_users_official_emails (
    hub_user_id UUID REFERENCES hub_users(id) NOT NULL,
    -- TODO: When Domain Ownership changes, this may break
    domain_id UUID REFERENCES domains(id) NOT NULL,

    official_email TEXT NOT NULL PRIMARY KEY,

    last_verified_at TIMESTAMP WITH TIME ZONE,
    -- Remember to set verification_code to NULL when the email is verified
    verification_code TEXT,
    verification_code_expires_at TIMESTAMP WITH TIME ZONE,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TABLE employer_primary_domains(
    employer_id UUID NOT NULL REFERENCES employers(id) ON DELETE CASCADE,
    domain_id UUID NOT NULL REFERENCES domains(id) ON DELETE CASCADE,

    PRIMARY KEY (employer_id),
    CONSTRAINT fk_employer_domain_match FOREIGN KEY (employer_id, domain_id)
        REFERENCES domains(employer_id, id)
);
---

CREATE TYPE org_user_roles AS ENUM (
    'ADMIN',
    'APPLICATIONS_CRUD',
    'APPLICATIONS_VIEWER',
    'COST_CENTERS_CRUD',
    'COST_CENTERS_VIEWER',
    'LOCATIONS_CRUD',
    'LOCATIONS_VIEWER',
    'ORG_USERS_CRUD',
    'ORG_USERS_VIEWER',
    'OPENINGS_CRUD',
    'OPENINGS_VIEWER'
);
CREATE TYPE org_user_states AS ENUM (
    'ACTIVE_ORG_USER',
    'INVITED_ORG_USER',
    'ADDED_ORG_USER',
    'DISABLED_ORG_USER',
    'REPLICATED_ORG_USER'
);
CREATE TABLE org_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL,
    name TEXT NOT NULL,
    password_hash TEXT,
    org_user_roles org_user_roles[] NOT NULL,
    org_user_state org_user_states NOT NULL,

--- As of now, we have only one org per employer. This may change in future.
    employer_id UUID REFERENCES employers(id) NOT NULL,
    CONSTRAINT uniq_email_employer_id UNIQUE (email, employer_id),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TYPE org_user_token_types AS ENUM (
    -- Sent as response to the TFA API
    'EMPLOYER_SESSION',
    'EMPLOYER_LTS',

    -- Sent as response to the SignIn API
    'EMPLOYER_TFA_TOKEN'
);
CREATE TABLE org_user_tokens (
    token TEXT CONSTRAINT org_user_tokens_pkey PRIMARY KEY,
    org_user_id UUID REFERENCES org_users(id) NOT NULL,
    token_type org_user_token_types NOT NULL,
    token_valid_till TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TABLE org_user_tfa_codes (
    code TEXT NOT NULL,
    tfa_token TEXT NOT NULL REFERENCES org_user_tokens(token) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TABLE org_user_invites (
    token TEXT CONSTRAINT org_user_invites_pkey PRIMARY KEY,
    org_user_id UUID REFERENCES org_users(id) NOT NULL,
    token_valid_till TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

---

CREATE TYPE cost_center_states AS ENUM ('ACTIVE_CC', 'DEFUNCT_CC');
CREATE TABLE org_cost_centers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cost_center_name TEXT NOT NULL,
    cost_center_state cost_center_states NOT NULL,
    notes TEXT NOT NULL,

    employer_id UUID REFERENCES employers(id) NOT NULL,
    CONSTRAINT uniq_cost_center_name_employer_id UNIQUE (cost_center_name, employer_id),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

---

CREATE TYPE location_states AS ENUM ('ACTIVE_LOCATION', 'DEFUNCT_LOCATION');
CREATE TABLE locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    country_code TEXT NOT NULL,
    postal_address TEXT NOT NULL,
    postal_code TEXT NOT NULL,
    openstreetmap_url TEXT,
    city_aka TEXT ARRAY,

    location_state location_states NOT NULL,

    employer_id UUID REFERENCES employers(id) NOT NULL,
    CONSTRAINT uniq_location_title_employer_id UNIQUE (title, employer_id),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

---

CREATE TYPE opening_states AS ENUM ('DRAFT_OPENING_STATE', 'ACTIVE_OPENING_STATE', 'SUSPENDED_OPENING_STATE', 'CLOSED_OPENING_STATE');
CREATE TYPE opening_types AS ENUM ('FULL_TIME_OPENING', 'PART_TIME_OPENING', 'CONTRACT_OPENING', 'INTERNSHIP_OPENING', 'UNSPECIFIED_OPENING');
CREATE TYPE education_levels AS ENUM ('BACHELOR_EDUCATION', 'MASTER_EDUCATION', 'DOCTORATE_EDUCATION', 'NOT_MATTERS_EDUCATION', 'UNSPECIFIED_EDUCATION');
CREATE TABLE openings (
    employer_id UUID REFERENCES employers(id) NOT NULL,
    id TEXT NOT NULL,
    CONSTRAINT openings_pk PRIMARY KEY (employer_id, id),
    CONSTRAINT opening_id_format_check CHECK (id ~ '^[0-9]{4}-[A-Za-z]{3}-[0-9]{1,2}-[0-9]+$'),

    title TEXT NOT NULL,
    positions INTEGER NOT NULL,
    jd TEXT NOT NULL,
    recruiter UUID REFERENCES org_users(id) NOT NULL,
    hiring_manager UUID REFERENCES org_users(id) NOT NULL,
    cost_center_id UUID REFERENCES org_cost_centers(id),
    employer_notes TEXT,
    remote_country_codes TEXT[],
    remote_timezones TEXT[],
    opening_type opening_types NOT NULL,
    yoe_min INTEGER NOT NULL,
    yoe_max INTEGER NOT NULL,
    min_education_level education_levels NOT NULL,
    salary_min NUMERIC,
    salary_max NUMERIC,
    salary_currency TEXT,
    state opening_states NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),
    last_updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),

    pagination_key BIGSERIAL
);

CREATE TABLE opening_hiring_team(
    employer_id UUID NOT NULL,
    opening_id TEXT NOT NULL,
    CONSTRAINT fk_opening FOREIGN KEY (employer_id, opening_id) REFERENCES openings (employer_id, id),

    hiring_team_mate_id UUID REFERENCES org_users(id) NOT NULL,

    PRIMARY KEY (employer_id, opening_id, hiring_team_mate_id)
);

CREATE TABLE opening_locations(
    employer_id UUID NOT NULL,
    opening_id TEXT NOT NULL,
    CONSTRAINT fk_opening FOREIGN KEY (employer_id, opening_id) REFERENCES openings (employer_id, id),

    location_id UUID REFERENCES locations(id) NOT NULL,

    PRIMARY KEY (employer_id, opening_id, location_id)
);

CREATE TABLE opening_watchers(
    employer_id UUID NOT NULL,
    opening_id TEXT NOT NULL,
    CONSTRAINT fk_opening FOREIGN KEY (employer_id, opening_id) REFERENCES openings (employer_id, id),

    watcher_id UUID REFERENCES org_users(id) NOT NULL,
    PRIMARY KEY (employer_id, opening_id, watcher_id)
);

CREATE TYPE application_color_tags AS ENUM ('GREEN', 'YELLOW', 'RED');
CREATE TYPE application_states AS ENUM ('APPLIED', 'REJECTED', 'SHORTLISTED', 'WITHDRAWN', 'EXPIRED');
CREATE TABLE applications (
    id TEXT PRIMARY KEY,
    employer_id UUID REFERENCES employers(id) NOT NULL,
    opening_id TEXT NOT NULL,
    CONSTRAINT fk_opening FOREIGN KEY (employer_id, opening_id) REFERENCES openings (employer_id, id),
    cover_letter TEXT NOT NULL,
    resume_sha TEXT NOT NULL,
    application_state application_states NOT NULL,

    color_tag application_color_tags,

    -- The user who applied for the opening
    hub_user_id UUID REFERENCES hub_users(id) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TYPE candidacy_states AS ENUM (
    -- What should be the state when a position is filled but a different
    -- candidate is in pipeline ? Or if the opening is no longer available for
    -- budget reasons ? Should we have a new state for it ?
    'INTERVIEWING',
    'OFFERED', 'OFFER_DECLINED', 'OFFER_ACCEPTED',
    'CANDIDATE_UNSUITABLE',
    'CANDIDATE_NOT_RESPONDING',
    'EMPLOYER_DEFUNCT'
);
CREATE TABLE candidacies(
    id TEXT PRIMARY KEY,

    application_id TEXT REFERENCES applications(id) NOT NULL,
    CONSTRAINT fk_application FOREIGN KEY (application_id) REFERENCES applications(id),

    employer_id UUID REFERENCES employers(id) NOT NULL,
    CONSTRAINT fk_employer FOREIGN KEY (employer_id) REFERENCES employers(id),

    opening_id TEXT NOT NULL,
    CONSTRAINT fk_opening FOREIGN KEY (employer_id, opening_id) REFERENCES openings (employer_id, id),

    candidacy_state candidacy_states NOT NULL,

    created_by UUID REFERENCES org_users(id) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TYPE comment_author_types AS ENUM ('ORG_USER', 'HUB_USER');
CREATE TABLE candidacy_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Discriminator field to identify the type of user
    author_type comment_author_types NOT NULL,

    -- Only one of these will be populated based on author_type
    org_user_id UUID REFERENCES org_users(id),
    hub_user_id UUID REFERENCES hub_users(id),

    comment_text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),

    -- Ensure exactly one user type is specified
    CONSTRAINT check_single_author CHECK (
        (author_type = 'ORG_USER' AND org_user_id IS NOT NULL AND hub_user_id IS NULL) OR
        (author_type = 'HUB_USER' AND hub_user_id IS NOT NULL AND org_user_id IS NULL)
    ),

    candidacy_id TEXT REFERENCES candidacies(id) NOT NULL,
    CONSTRAINT fk_candidacy FOREIGN KEY (candidacy_id) REFERENCES candidacies(id),

    employer_id UUID REFERENCES employers(id) NOT NULL,
    CONSTRAINT fk_employer FOREIGN KEY (employer_id) REFERENCES employers(id)
);

-- Index for chronological fetching
CREATE INDEX idx_candidacy_comments_chronological ON candidacy_comments(candidacy_id, created_at DESC);

---

CREATE TYPE interview_types AS ENUM (
    'IN_PERSON',
    'VIDEO_CALL',
    'TAKE_HOME',
    'OTHER_INTERVIEW'
);
CREATE TYPE interview_states AS ENUM (
    'SCHEDULED_INTERVIEW',
    'COMPLETED_INTERVIEW',
    'CANCELLED_INTERVIEW'
);
CREATE TYPE interviewers_decisions AS ENUM (
    'STRONG_YES',
    'YES',
    'NEUTRAL',
    'NO',
    'STRONG_NO'
);
CREATE TYPE rsvp_status AS ENUM (
    'YES',
    'NO',
    'NOT_SET'
);
CREATE TABLE interviews(
    id TEXT PRIMARY KEY,

    interview_type interview_types NOT NULL,
    interview_state interview_states NOT NULL,
    candidate_rsvp rsvp_status NOT NULL DEFAULT 'NOT_SET',

    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    CONSTRAINT valid_interview_duration CHECK (end_time > start_time),

    description TEXT,

    created_by UUID REFERENCES org_users(id) NOT NULL,

    candidacy_id TEXT REFERENCES candidacies(id) NOT NULL,
    employer_id UUID REFERENCES employers(id) NOT NULL,

    CONSTRAINT fk_candidacy FOREIGN KEY (candidacy_id) REFERENCES candidacies(id),

    interviewers_decision interviewers_decisions,
    positives TEXT,
    negatives TEXT,
    overall_assessment TEXT,
    feedback_to_candidate TEXT,
    feedback_submitted_by UUID REFERENCES org_users(id),

    feedback_submitted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT valid_feedback CHECK (
        (feedback_submitted_by IS NULL AND feedback_submitted_at IS NULL AND interviewers_decision IS NULL) OR
        (feedback_submitted_by IS NOT NULL AND feedback_submitted_at IS NOT NULL AND interviewers_decision IS NOT NULL)
    ),

    completed_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT valid_completion CHECK (
        (interview_state = 'COMPLETED_INTERVIEW' AND completed_at IS NOT NULL) OR
        (interview_state != 'COMPLETED_INTERVIEW' AND completed_at IS NULL)
    ),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TABLE interview_interviewers (
    interview_id TEXT REFERENCES interviews(id) NOT NULL,
    interviewer_id UUID REFERENCES org_users(id) NOT NULL,
    employer_id UUID REFERENCES employers(id) NOT NULL,
    rsvp_status rsvp_status NOT NULL DEFAULT 'NOT_SET',

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),

    PRIMARY KEY (interview_id, interviewer_id)
);

-- Add indices for common query patterns
CREATE INDEX idx_interviews_candidacy ON interviews(candidacy_id);
CREATE INDEX idx_interviews_employer ON interviews(employer_id);
CREATE INDEX idx_interviews_state ON interviews(interview_state);

CREATE TABLE opening_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TABLE opening_tag_mappings (
    employer_id UUID NOT NULL,
    opening_id TEXT NOT NULL,
    CONSTRAINT fk_opening FOREIGN KEY (employer_id, opening_id) REFERENCES openings (employer_id, id),
    tag_id UUID REFERENCES opening_tags(id) NOT NULL,
    PRIMARY KEY (employer_id, opening_id, tag_id)
);

-- Seed data for common opening tags
INSERT INTO opening_tags (name) VALUES
    ('DevOps'),
    ('Golang'),
    ('Database Administrator'),
    ('Frontend Developer'),
    ('Backend Developer'),
    ('Full Stack Developer'),
    ('Site Reliability Engineer'),
    ('Cloud Engineer'),
    ('Data Scientist'),
    ('Machine Learning Engineer'),
    ('Product Manager'),
    ('UI/UX Designer'),
    ('QA Engineer'),
    ('Security Engineer'),
    ('Mobile Developer'),
    ('Technical Writer'),
    ('Engineering Manager'),
    ('Technical Support'),
    ('Business Analyst'),
    ('System Administrator');

CREATE OR REPLACE FUNCTION get_or_create_dummy_employer(p_domain_name text)
RETURNS UUID AS $$
DECLARE
    employer_id UUID;
    domain_id UUID;
BEGIN
    -- First check if domain already exists and has an employer
    SELECT d.employer_id INTO employer_id
    FROM domains d
    WHERE d.domain_name = p_domain_name AND d.employer_id IS NOT NULL;

    IF FOUND THEN
        RETURN employer_id;
    END IF;

    -- Create dummy employer if not exists
    INSERT INTO employers (
        client_id_type,
        employer_state,
        company_name,
        onboard_admin_email
    ) VALUES (
        'DOMAIN',
        'HUB_ADDED_EMPLOYER',
        p_domain_name,  -- Use domain name as company name for dummy record
        'admin@' || p_domain_name -- TODO: Perhaps this could just be NULL ?
    )
    RETURNING id INTO employer_id;

    -- Create or update domain
    INSERT INTO domains (
        domain_name,
        domain_state,
        employer_id
    ) VALUES (
        p_domain_name,
        'UNVERIFIED',
        employer_id
    )
    ON CONFLICT (domain_name) DO UPDATE
    SET employer_id = EXCLUDED.employer_id
    RETURNING id INTO domain_id;

    -- Set as primary domain
    INSERT INTO employer_primary_domains (employer_id, domain_id)
    VALUES (employer_id, domain_id);

    RETURN employer_id;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE work_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hub_user_id UUID REFERENCES hub_users(id) NOT NULL,
    employer_id UUID REFERENCES employers(id) NOT NULL,
    title TEXT NOT NULL,
    start_date TEXT NOT NULL,
    end_date TEXT,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

-- Colleague Connection Scenarios:
--
-- 1. Initial Connection Request:
--    - HubUser1 sends invitation to HubUser2 (COLLEAGUING_PENDING)
--    - Only one pending invitation can exist between two users at a time
--
-- 2. Accepting Connection:
--    - HubUser2 accepts HubUser1's invitation (COLLEAGUING_ACCEPTED)
--    - Both users become colleagues
--
-- 3. Rejecting Connection:
--    - HubUser2 rejects HubUser1's invitation (COLLEAGUING_REJECTED)
--    - HubUser1 cannot send another invitation to HubUser2
--    - HubUser2 can still send an invitation to HubUser1
--
-- 4. Unlinking Connection:
--    - After COLLEAGUING_ACCEPTED, either user can unlink (COLLEAGUING_UNLINKED)
--    - If HubUser2 unlinks, HubUser1 cannot send another invitation
--    - If HubUser1 unlinks, HubUser2 cannot send another invitation
--
-- 5. Colleaguable Rules:
--    - A user cannot connect with themselves
--    - A user cannot send invitation if there's a pending invitation
--    - A user cannot send invitation if they are already connected
--    - A user cannot send invitation if the target user previously rejected/unlinked their invitation
--    - A user CAN send invitation even if they previously rejected/unlinked the target user's invitation
--
CREATE TYPE colleaguing_states AS ENUM (
    'COLLEAGUING_PENDING',           -- Initial state when invitation is sent
    'COLLEAGUING_ACCEPTED',         -- Both users are connected
    'COLLEAGUING_REJECTED',         -- The invitation was rejected
    'COLLEAGUING_UNLINKED'         -- The connection was unlinked after being accepted
);

-- Drop the existing colleagues table
DROP TABLE IF EXISTS colleagues;

CREATE TABLE colleague_connections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    requester_id UUID REFERENCES hub_users(id) NOT NULL,
    requested_id UUID REFERENCES hub_users(id) NOT NULL,
    state colleaguing_states NOT NULL,

    -- Track who performed which action and when
    rejected_by UUID REFERENCES hub_users(id),
    rejected_at TIMESTAMP WITH TIME ZONE,
    unlinked_by UUID REFERENCES hub_users(id),
    unlinked_at TIMESTAMP WITH TIME ZONE,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),

    CONSTRAINT check_self_reference CHECK (requester_id != requested_id),
    CONSTRAINT unique_connection UNIQUE (requester_id, requested_id),

    -- Ensure rejected_by and rejected_at are set together
    CONSTRAINT reject_action_integrity CHECK (
        (rejected_by IS NULL AND rejected_at IS NULL) OR
        (rejected_by IS NOT NULL AND rejected_at IS NOT NULL)
    ),

    -- Ensure unlinked_by and unlinked_at are set together
    CONSTRAINT unlink_action_integrity CHECK (
        (unlinked_by IS NULL AND unlinked_at IS NULL) OR
        (unlinked_by IS NOT NULL AND unlinked_at IS NOT NULL)
    ),

    -- Ensure state matches the action fields
    CONSTRAINT state_action_integrity CHECK (
        (state = 'COLLEAGUING_REJECTED' AND rejected_by IS NOT NULL AND unlinked_by IS NULL) OR
        (state = 'COLLEAGUING_UNLINKED' AND unlinked_by IS NOT NULL AND rejected_by IS NULL) OR
        (state IN ('COLLEAGUING_PENDING', 'COLLEAGUING_ACCEPTED') AND rejected_by IS NULL AND unlinked_by IS NULL)
    )
);

-- Create indexes for common queries
CREATE INDEX idx_colleague_connections_requester ON colleague_connections(requester_id);
CREATE INDEX idx_colleague_connections_requested ON colleague_connections(requested_id);

-- Function to check if a user can send a colleague request to another user
CREATE OR REPLACE FUNCTION is_colleaguable(seeking_user UUID, target_user UUID)
RETURNS BOOLEAN AS $$
BEGIN
    -- Prevent self-connections
    IF seeking_user = target_user THEN
        RETURN FALSE;
    END IF;

    -- Check if there's a pending invitation in either direction
    IF EXISTS (
        SELECT 1 FROM colleague_connections
        WHERE (requester_id = seeking_user AND requested_id = target_user
               OR requester_id = target_user AND requested_id = seeking_user)
        AND state = 'COLLEAGUING_PENDING'
    ) THEN
        RETURN FALSE;
    END IF;

    -- Check if there's an existing accepted connection
    IF EXISTS (
        SELECT 1 FROM colleague_connections
        WHERE (requester_id = seeking_user AND requested_id = target_user
               OR requester_id = target_user AND requested_id = seeking_user)
        AND state = 'COLLEAGUING_ACCEPTED'
    ) THEN
        RETURN FALSE;
    END IF;

    -- Check if target_user has previously rejected/unlinked a connection from seeking_user
    IF EXISTS (
        SELECT 1 FROM colleague_connections
        WHERE requester_id = seeking_user
        AND requested_id = target_user
        AND (
            (state = 'COLLEAGUING_REJECTED' AND rejected_by = target_user) OR
            (state = 'COLLEAGUING_UNLINKED' AND unlinked_by = target_user)
        )
    ) THEN
        RETURN FALSE;
    END IF;

    -- Check if seeking_user has previously rejected/unlinked a connection from target_user
    IF EXISTS (
        SELECT 1 FROM colleague_connections
        WHERE requester_id = target_user
        AND requested_id = seeking_user
        AND (
            (state = 'COLLEAGUING_REJECTED' AND rejected_by = seeking_user) OR
            (state = 'COLLEAGUING_UNLINKED' AND unlinked_by = seeking_user)
        )
    ) THEN
        RETURN FALSE;
    END IF;

    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;

-- Function to check if two users are connected
CREATE OR REPLACE FUNCTION are_colleagues(user1_id UUID, user2_id UUID)
RETURNS BOOLEAN AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1 FROM colleague_connections
        WHERE (requester_id = user1_id AND requested_id = user2_id
               OR requester_id = user2_id AND requested_id = user1_id)
        AND state = 'COLLEAGUING_ACCEPTED'
    );
END;
$$ LANGUAGE plpgsql;

-- Table to track old files that need cleanup
CREATE TABLE stale_files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),
    cleaned_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT unique_file_path UNIQUE (file_path)
);

-- Index for finding unprocessed stale files
CREATE INDEX idx_stale_files_unprocessed
    ON stale_files(cleaned_at)
    WHERE cleaned_at IS NULL;

-- Drop old tables if they exist
DROP TABLE IF EXISTS stale_profile_pictures;
DROP TYPE IF EXISTS file_type;

COMMIT;
