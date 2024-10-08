CREATE TABLE IF NOT EXISTS sessions (
	token TEXT PRIMARY KEY,
	data BYTEA NOT NULL,
	expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX SessionExpiryIdx ON sessions (expiry);