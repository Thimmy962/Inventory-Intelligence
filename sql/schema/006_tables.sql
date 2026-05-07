-- +goose up
CREATE TABLE IF NOT EXISTS staffs (
    id TEXT PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    username TEXT UNIQUE NOT NULL,
    pword TEXT UNIQUE NOT NULL,
    is_manager BOOLEAN NOT NULL,
    refresh_token TEXT UNIQUE
);

-- +goose down
DROP TABLE staffs;
