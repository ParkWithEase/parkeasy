-- Parking spot details
CREATE TABLE IF NOT EXISTS ParkingSpot (
  ParkingSpotId BIGSERIAL PRIMARY KEY,
  UserId BIGINT NOT NULL REFERENCES Users(UserId),
  ParkingSpotUUID UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  PostalCode TEXT NOT NULL, 
  CountryCode TEXT NOT NULL,
  City TEXT NOT NULL,
  State TEXT NOT NULL,
  StreetAddress TEXT NOT NULL,
  Longitude DECIMAL(8,5) NOT NULL,
  Latitude DECIMAL(8,5) NOT NULL,
  HasShelter BOOLEAN NOT NULL DEFAULT false,
  HasPlugIn BOOLEAN NOT NULL DEFAULT false,
  HasChargingStation BOOLEAN NOT NULL DEFAULT false,
  PricePerHour DECIMAL NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS ParkingSpotUUIDIdx ON ParkingSpot(ParkingSpotUUID);
CREATE UNIQUE INDEX IF NOT EXISTS ParkingSpotCoordinateIdx ON ParkingSpot(Longitude, Latitude);

CREATE EXTENSION IF NOT EXISTS earthdistance CASCADE;