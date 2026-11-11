-- Dev seed: 1 admin + 2 users
-- Password for all accounts: Password123!
-- bcrypt cost 10

INSERT INTO users (id, email, password_hash, role, is_active)
VALUES
    (
        '11111111-1111-1111-1111-111111111111',
        'admin@filetemp.local',
        '$2b$10$CqTNGHpppfXnzFkVRGD9aOiqfXymdNcB2zkjJ80uU9b8kScCIfyzi',
        'ADMIN',
        TRUE
    ),
    (
        '22222222-2222-2222-2222-222222222222',
        'alice@filetemp.local',
        '$2b$10$CqTNGHpppfXnzFkVRGD9aOiqfXymdNcB2zkjJ80uU9b8kScCIfyzi',
        'USER',
        TRUE
    ),
    (
        '33333333-3333-3333-3333-333333333333',
        'bob@filetemp.local',
        '$2b$10$CqTNGHpppfXnzFkVRGD9aOiqfXymdNcB2zkjJ80uU9b8kScCIfyzi',
        'USER',
        TRUE
    )
ON CONFLICT (email) DO NOTHING;
