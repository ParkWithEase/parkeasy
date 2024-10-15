-- Each record corresponds to 30 mins of time
CREATE TABLE IF NOT EXISTS TimeUnit (
  UnitNum SMALLINT NOT NULL,
  Date DATE NOT NULL,
  ListingId BIGINT NOT NULL REFERENCES Listing(ListingId),
  BookingId BIGINT REFERENCES Booking(BookingID),
  PRIMARY KEY (UnitNum, Date, ListingId)
);