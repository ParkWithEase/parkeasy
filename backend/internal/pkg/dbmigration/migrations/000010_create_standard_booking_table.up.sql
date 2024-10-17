-- Standard Booking details
CREATE TABLE IF NOT EXISTS StandardBooking (
  StandardBookingId BIGSERIAL PRIMARY KEY,
  StandardBookingUUID UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  BookingID BIGINT NOT NULL REFERENCES Booking(BookingID),
  ListingId BIGINT NOT NULL, -- Add back FK
  StartUnitNum SMALLINT NOT NULL,
  EndUnitNum SMALLINT NOT NULL,
  Date DATE NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS StandardBookingUUIDIdx ON StandardBooking(StandardBookingUUID);