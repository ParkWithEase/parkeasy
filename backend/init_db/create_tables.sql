--======================================================--
-- CREATES
--======================================================--

-- User details
CREATE TABLE IF NOT EXISTS Users (
  UserId SERIAL PRIMARY KEY,
  PasswordHash TEXT,
  FirstName TEXT,
  LastName TEXT,
  Email TEXT,
  IsVerified BOOLEAN
);

-- Session token details for each user
CREATE TABLE IF NOT EXISTS SessionToken (
  SessionToken TEXT,
  Device TEXT,
  ExpirationTime TIMESTAMP,
  UserId INTEGER,
  PRIMARY KEY (SessionToken, UserId),
  CONSTRAINT fk_sessiontoken_user FOREIGN KEY (UserId) REFERENCES Users(UserId)
);

-- Reset password token details for each user
CREATE TABLE IF NOT EXISTS ResetPasswordToken (
  ResetPasswordToken TEXT,
  ExpirationTime TIMESTAMP,
  IsUsed BOOLEAN,
  UserId INTEGER,
  PRIMARY KEY (ResetPasswordToken, UserId),
  CONSTRAINT fk_resetpasswordtoken_user FOREIGN KEY (UserId) REFERENCES Users(UserId)
);

-- Validation token details for each user
CREATE TABLE IF NOT EXISTS ValidationToken (
  ValidationToken TEXT,
  IsUsed BOOLEAN,
  UserId INTEGER,
  PRIMARY KEY (ValidationToken, UserId),
  CONSTRAINT fk_validationtoken_user FOREIGN KEY (UserId) REFERENCES Users(UserId)
);

-- Car details owned by a user
CREATE TABLE IF NOT EXISTS Car (
  CarId SERIAL PRIMARY KEY,
  LicensePlate TEXT,
  Make TEXT,
  Model TEXT,
  Color TEXT,
  UserId INTEGER,
  CONSTRAINT fk_car_user FOREIGN KEY (UserId) REFERENCES Users(UserId)
);

-- Parking spot details owned by a user
CREATE TABLE IF NOT EXISTS ParkingSpot (
  ParkingSpotId SERIAL PRIMARY KEY,
  UserId INTEGER,
  Address TEXT,
  Longitude REAL,
  Latitude REAL,
  HasShelter BOOLEAN,
  HasPlugIn BOOLEAN,
  HasChargingStation BOOLEAN,
  CONSTRAINT fk_parkingspot_user FOREIGN KEY (UserId) REFERENCES Users(UserId)
);

-- Master table of bookings with fields that are common to standard and auction bookings
CREATE TABLE IF NOT EXISTS Booking (
  BookingID SERIAL PRIMARY KEY,
  BuyerUserId INTEGER,
  PaidAmount REAL,
  CONSTRAINT fk_booking_user FOREIGN KEY (BuyerUserId) REFERENCES Users(UserId)
);

-- Listing details produced by a user for a given parking spot
CREATE TABLE IF NOT EXISTS Listing (
  ListingId SERIAL PRIMARY KEY,
  ParkingSpotId INTEGER,
  PricePerHour REAL,
  IsActive BOOLEAN,
  CONSTRAINT fk_listing_parkingspot FOREIGN KEY (ParkingSpotId) REFERENCES ParkingSpot(ParkingSpotId)
);

-- Each record corresponds to 30 mins of time for a given listing
CREATE TABLE IF NOT EXISTS TimeUnit (
  UnitNum SMALLINT,
  Date DATE,
  ListingId INTEGER,
  BookingId INTEGER,
  PRIMARY KEY (UnitNum, Date, ListingId),
  CONSTRAINT fk_timeunit_listing FOREIGN KEY (ListingId) REFERENCES Listing(ListingId),
  CONSTRAINT fk_timeunit_booking FOREIGN KEY (BookingId) REFERENCES Booking(BookingID)
);

-- Booking details from a given listing made by a user
CREATE TABLE IF NOT EXISTS StandardBooking (
  StandardBookingId SERIAL PRIMARY KEY,
  ListingId INTEGER,
  StartUnitNum SMALLINT,
  EndUnitNum SMALLINT,
  Date DATE,
  CONSTRAINT fk_standardbooking_listing FOREIGN KEY (ListingId) REFERENCES Listing(ListingId),
  CONSTRAINT fk_standardbooking_booking FOREIGN KEY (StandardBookingId) REFERENCES Booking(BookingID)
);

-- A time window on a particular day of a given listing that is to be auctioned
CREATE TABLE IF NOT EXISTS AuctionTimeSlot (
  AuctionTimeSlotId SERIAL PRIMARY KEY,
  ListingId INTEGER,
  CurrentBidId INTEGER,
  Date DATE,
  StartUnitNum SMALLINT,
  EndUnitNum SMALLINT,
  BidAmount REAL,
  CONSTRAINT fk_auctiontimeslot_listing FOREIGN KEY (ListingId) REFERENCES Listing(ListingId)
);

-- Auction booking details from a given listing made by a user
CREATE TABLE IF NOT EXISTS AuctionBooking (
  AuctionBookingId SERIAL PRIMARY KEY,
  AuctionTimeSlotId INTEGER,
  CONSTRAINT fk_auctionbooking_auctiontimeslot FOREIGN KEY (AuctionTimeSlotId) REFERENCES AuctionTimeSlot(AuctionTimeSlotId),
  CONSTRAINT fk_auctionbooking_booking FOREIGN KEY (AuctionBookingId) REFERENCES Booking(BookingID)
);