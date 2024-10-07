-- Auth details
CREATE TABLE IF NOT EXISTS Auth (
  AuthID SERIAL PRIMARY KEY NOT NULL,
  AuthUUID UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  Email TEXT NOT NULL,
  PasswordHash TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS AuthUUIDs ON Auth (AuthUUID);
