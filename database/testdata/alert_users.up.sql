INSERT INTO users (id, name, email, secret, email_verify) VALUES
    (DEFAULT, 'abc', 'abc@me.com', 'secret', true),
    (DEFAULT, 'def', 'def@me.com', 'secret', true),
    (DEFAULT, 'ghi', 'ghi@me.com', 'secret', true);

INSERT INTO alerts (id, form_type, filer_cik, user_id) VALUES
    (DEFAULT, '8-K', '123', (SELECT id from users WHERE name='abc')),
    (DEFAULT, '4', '456', (SELECT id from users WHERE name='abc')),
    (DEFAULT, NULL, '123', (SELECT id from users WHERE name='def')),
    (DEFAULT, '4', NULL, (SELECT id from users WHERE name='def')),
    (DEFAULT, NULL, '456', (SELECT id from users WHERE name='def')),
    (DEFAULT, '8-K', NULL, (SELECT id from users WHERE name='ghi'));
