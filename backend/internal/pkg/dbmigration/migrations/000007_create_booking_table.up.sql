-- Master table of bookings with fields that are common to standard and auction bookings
CREATE TABLE IF NOT EXISTS Booking (
  BookingId BIGSERIAL PRIMARY KEY,
  BookingUUID UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  UserId BIGINT NOT NULL REFERENCES Users(UserId),
  ParkingSpotId BIGINT NOT NULL REFERENCES ParkingSpot(ParkingSpotId),
  PaidAmount DECIMAL NOT NULL 
);

CREATE UNIQUE INDEX IF NOT EXISTS BookingUUIDIdx ON Booking(BookingUUID);