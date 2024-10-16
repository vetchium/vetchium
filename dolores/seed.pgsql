BEGIN;

--- email table primary key ids should start from 11 and in 2 digits
--- employer table primary key ids should start from 201 and in 3 digits
--- domain table primary key ids should start from 3001 and in 4 digits
--- org_users table primary key ids should start from 40001 and in 5 digits

INSERT INTO public.emails (
    id,
    email_from,
    email_to,
    email_cc,
    email_bcc,
    email_subject,
    email_html_body,
    email_text_body,
    email_state,
    created_at,
    processed_at
) VALUES (
    11,
    'no-reply@vetchi.org',
    ARRAY['admin@domain-onboarded.example'],
    NULL,
    NULL,
    'Welcome to Vetchi',
    '<h1>Welcome to Vetchi</h1>',
    'Welcome to Vetchi',
    'PROCESSED',
    timezone('UTC'::text, now()), timezone('UTC'::text, now())
);

INSERT INTO public.employers (
    id,
    client_id_type,
    employer_state,
    onboard_admin_email,
    onboard_secret_token,
    token_valid_till,
    onboard_email_id,
    created_at
) VALUES (
    201,
    'DOMAIN',
    'ONBOARDED',
    'admin@domain-onboarded.example',
    '$2y$05$nTpbRp.SqiP0baLK/Am40.rbMkJItfiGJD7E9.n7k/d9b4LAAh2P6',
    timezone('UTC'::text, now()) + interval '1 day', 
    11, 
    timezone('UTC'::text, now())
);

INSERT INTO public.domains (
    id,
    domain_name,
    domain_state,
    employer_id,
    created_at
) VALUES(
    1001,
    'domain-onboarded.example',
    'VERIFIED',
    201,
    timezone('UTC'::text, now())
);

INSERT INTO public.org_users(
    id,
    email,
    password_hash,
    org_user_role,
    org_user_state,
    employer_id,
    created_at
) VALUES(
    1001,
    'admin@domain-onboarded.example',
    --- password is Password123$
    '$2y$05$nTpbRp.SqiP0baLK/Am40.rbMkJItfiGJD7E9.n7k/d9b4LAAh2P6',
    'ADMIN',
    'ACTIVE',
    201,
    timezone('UTC'::text, now())
);

COMMIT;