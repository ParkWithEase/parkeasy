-- Master table of bookings with fields that are common to standard and auction bookings
CREATE TABLE IF NOT EXISTS Booking (
  BookingID BIGSERIAL PRIMARY KEY,
  BuyerUserId BIGINT NOT NULL REFERENCES Users(UserId),
  BookingUUID UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  PaidAmount REAL NOT NULL 
);

CREATE UNIQUE INDEX IF NOT EXISTS BookingUUIDIdx ON Booking(BookingUUID);