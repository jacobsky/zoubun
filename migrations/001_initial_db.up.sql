-- used to store arbitrary configuration values used by the server.
CREATE TABLE config (
    config_key TEXT UNIQUE,
    config_value JSON
);

-- User table which is to track registrations
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    verified BOOLEAN,
    is_admin BOOLEAN,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Verification challenge tracker
CREATE TABLE user_challenge (
    userid SERIAL UNIQUE,
    challenge TEXT,
    deadline DATETIME
);


-- stores the API keys
CREATE TABLE user_keys (
    userid SERIAL,
    apikey TEXT UNIQUE,
    creation_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted BOOLEAN,
    FOREIGN KEY (userid) REFERENCES users (id)
);

-- User game progress
CREATE TABLE counters (
    userid SERIAL,
    current_count BIGINT,
    FOREIGN KEY (userid) REFERENCES users (id)
);


-- Add some basic data to the database
INSERT INTO config (config_key, config_value) VALUES ("motd", "Hello! Welcome to Zoubun!  こんにちは！増分へようこそ！");
INSERT INTO users (username, verified, is_admin) VALUES ("admin", true, true);
INSERT INTO counters (userid, current_count) VALUES (0, 0)
