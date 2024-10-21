-- Each record corresponds to 30 mins of time
CREATE TABLE IF NOT EXISTS TimeUnit (
  StartTime TIMESTAMPTZ NOT NULL,
  EndTime TIMESTAMPTZ NOT NULL,
  ParkingSpotUUID UUID NOT NULL REFERENCES ParkingSpot(ParkingSpotUUID),
  BookingId BIGINT DEFAULT NULL REFERENCES Booking(BookingID),
  Status TEXT NOT NULL,
  PRIMARY KEY (StartTime, EndTime, ParkingSpotUUID)
);

