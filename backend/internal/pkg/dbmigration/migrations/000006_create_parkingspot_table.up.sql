CREATE EXTENSION IF NOT EXISTS cube;
CREATE EXTENSION IF NOT EXISTS earthdistance;

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
  PricePerHour DECIMAL NOT NULL,
  CONSTRAINT latlon_overlap_exclude EXCLUDE USING gist (earth_box(ll_to_earth(latitude, longitude), 3) with &&)
);

CREATE UNIQUE INDEX IF NOT EXISTS ParkingSpotUUIDIdx ON ParkingSpot(ParkingSpotUUID);
CREATE INDEX IF NOT EXISTS ParkingSpotCoordinateIdx ON ParkingSpot USING gist (ll_to_earth(latitude, longitude));

