INSERT INTO users (id, name, email, secret, email_verify) VALUES
    ('jd', 'John Doe', 'john@example.com', 'secretpass', true);

INSERT INTO alerts (id, form_type, filer_cik, user_id) VALUES
    (DEFAULT, NULL, '0001418091', 'jd'),
    (DEFAULT, '10-Q', '0001318605', 'jd'),
    (DEFAULT, '10-Q', '0001403161', 'jd'),
    (DEFAULT, NULL, '0001652044', 'jd'),
    (DEFAULT, NULL, '0000320193', 'jd'),
    (DEFAULT, 'S-1', NULL, 'jd');