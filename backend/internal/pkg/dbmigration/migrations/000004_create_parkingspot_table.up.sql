-- Parking spot details
CREATE TABLE IF NOT EXISTS ParkingSpot (
  ParkingSpotId SERIAL PRIMARY KEY,
  UserId INTEGER UNIQUE NOT NULL REFERENCES Users(UserID),
  ParkingSpotUUID UUID UNIQUE NOT NULL,
  PostalCode TEXT NOT NULL, 
  CountryCode TEXT NOT NULL,
  City TEXT NOT NULL,
  StreetAddress TEXT NOT NULL,
  Longitude REAL NOT NULL,
  Latitude REAL NOT NULL,
  HasShelter BOOLEAN NOT NULL DEFAULT false,
  HasPlugIn BOOLEAN NOT NULL DEFAULT false,
  HasChargingStation BOOLEAN NOT NULL DEFAULT false,
  IsPublic BOOLEAN NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX IF NOT EXISTS ParkingSpotUUIDs ON ParkingSpot(ParkingSpotUUID);