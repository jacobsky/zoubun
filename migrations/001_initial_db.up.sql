-- used to store arbitrary configuration values used by the server.
CREATE TABLE config (
    config_key TEXT UNIQUE,
    config_value JSON
);

-- User table which is to track registrations
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    verification_code TEXT,
    verified BOOLEAN,
    is_admin BOOLEAN,
    creation_date DATETIME
);

-- stores the API keys
CREATE TABLE user_keys (
    userid SERIAL,
    apikey TEXT UNIQUE,
    creation_date DATETIME,
    deleted BOOLEAN,
    FOREIGN KEY (userid) REFERENCES users (id)
);

-- User game progress
CREATE TABLE counters (
    userid SERIAL,
    current_count BIGINT,
    FOREIGN KEY (userid) REFERENCES users (id)
);
