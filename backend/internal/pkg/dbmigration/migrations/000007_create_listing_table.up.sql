-- Listing details
CREATE TABLE IF NOT EXISTS Listing (
  ListingId BIGSERIAL PRIMARY KEY,
  ParkingSpotId BIGINT NOT NULL REFERENCES ParkingSpot(ParkingSpotId),
  ListingUUID UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
  PricePerHour REAL,
  IsActive BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE UNIQUE INDEX IF NOT EXISTS ListingUUIDIdx ON Listing(ListingUUID);
