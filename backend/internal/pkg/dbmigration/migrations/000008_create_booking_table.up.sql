-- Master table of bookings with fields that are common to standard and auction bookings
CREATE TABLE IF NOT EXISTS Booking (
  BookingID BIGSERIAL PRIMARY KEY,
  BuyerUserId BIGINT NOT NULL REFERENCES Users(UserId),
  PaidAmount REAL NOT NULL 
);