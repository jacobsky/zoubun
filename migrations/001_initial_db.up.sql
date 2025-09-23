-- used to store arbitrary configuration values used by the server.
CREATE TABLE config (
    config_key TEXT UNIQUE NOT NULL,
    -- In a more complex scenario, this might be better to set as a JSON type.
    -- For this project it will not get used beyond strings and intengers.
    -- Hence leaving as text for now.
    config_value TEXT NOT NULL
);

-- User table which is to track registrations
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    verified BOOLEAN DEFAULT FALSE,
    is_admin BOOLEAN DEFAULT FALSE,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Verification challenge tracker
CREATE TABLE user_challenge (
    userid SERIAL UNIQUE,
    challenge TEXT,
    deadline TIMESTAMP
);


-- stores the API keys
CREATE TABLE user_keys (
    userid SERIAL,
    apikey1 TEXT UNIQUE,
    apikey2 TEXT UNIQUE,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (userid) REFERENCES users (id)
);

-- User game progress
CREATE TABLE counters (
    userid SERIAL,
    current_count BIGINT DEFAULT 0,
    FOREIGN KEY (userid) REFERENCES users (id)
);


-- Add some basic data to the database
INSERT INTO config (config_key, config_value)
VALUES ('motd', '"Hello! Welcome to Zoubun!  こんにちは！増分へようこそ！"');

WITH created_user AS (INSERT INTO users (username, verified, is_admin) VALUES ('admin', true, true) RETURNING id)
INSERT INTO counters (userid) SELECT created_user.id FROM created_user;
