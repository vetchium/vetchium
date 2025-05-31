BEGIN;

-- Add approved domains for hub user registration testing
INSERT INTO hub_user_signup_approved_domains (domain_name, notes) VALUES
('0031-registration.example', 'Test domain for hub user registration tests'),
('0031-test-registration.example', 'Additional test domain for registration tests');

COMMIT; 