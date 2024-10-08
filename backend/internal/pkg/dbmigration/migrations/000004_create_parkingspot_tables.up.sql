-- Parking spot details
CREATE TABLE IF NOT EXISTS ParkingSpot (
  ParkingSpotId SERIAL PRIMARY KEY,
  ParkingSpotUUID UNIQUE NOT NULL,
  UserId INTEGER NOT NULL REFERENCES Users(UserID),
  Address TEXT NOT NULL,
  Longitude REAL NOT NULL,
  Latitude REAL NOT NULL,
  HasShelter BOOLEAN NOT NULL DEFAULT false,
  HasPlugIn BOOLEAN NOT NULL DEFAULT false,
  HasChargingStation BOOLEAN NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX IF NOT EXISTS ParkingSpotUUIDs ON ParkingSpot (ParkingSpotUUID);