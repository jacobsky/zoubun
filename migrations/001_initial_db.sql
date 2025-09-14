CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    userid TEXT UNIQUE NOT NULL,
    verification_code TEXT,
    verified BOOLEAN,
    api_key TEXT UNIQUE
);
CREATE TABLE counters (
    userid INTEGER,
    current_count BIGINT,
    FOREIGN KEY (userid) REFERENCES users (id)
);

---- create above / drop below ----

DROP TABLE users;
DROP TABLE counters;
