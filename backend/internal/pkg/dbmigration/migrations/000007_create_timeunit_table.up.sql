CREATE EXTENSION IF NOT EXISTS btree_gist;

-- Each record corresponds to 30 mins of time
CREATE TABLE IF NOT EXISTS TimeUnit (
  TimeRange TSTZRANGE NOT NULL,
  ParkingSpotId BIGINT NOT NULL REFERENCES ParkingSpot(ParkingSpotId),
  PRIMARY KEY (TimeRange, ParkingSpotId),
  EXCLUDE USING GIST (TimeRange WITH &&, ParkingSpotId WITH =)
);
