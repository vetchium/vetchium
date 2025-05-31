BEGIN;

-- Add approved domains for hub user registration testing
INSERT INTO hub_user_signup_approved_domains (domain_name, notes) VALUES
('registration.example', 'Test domain for hub user registration tests'),
('test-registration.example', 'Additional test domain for registration tests');

COMMIT; 