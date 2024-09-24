--======================================================--
-- DROPS
--======================================================--

DROP TABLE IF EXISTS AuctionBooking CASCADE;
DROP TABLE IF EXISTS AuctionTimeSlot CASCADE;
DROP TABLE IF EXISTS StandardBooking CASCADE;
DROP TABLE IF EXISTS TimeUnit CASCADE;
DROP TABLE IF EXISTS Listing CASCADE;
DROP TABLE IF EXISTS Booking CASCADE;
DROP TABLE IF EXISTS ParkingSpot CASCADE;
DROP TABLE IF EXISTS Car CASCADE;
DROP TABLE IF EXISTS ResetPasswordToken CASCADE;
DROP TABLE IF EXISTS ValidationToken CASCADE;
DROP TABLE IF EXISTS SessionToken CASCADE;
DROP TABLE IF EXISTS "User" CASCADE;




--======================================================--
-- CREATES
--======================================================--

-- User details
CREATE TABLE IF NOT EXISTS "User" (
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
  CONSTRAINT fk_sessiontoken_user FOREIGN KEY (UserId) REFERENCES "User"(UserId)
);

-- Reset password token details for each user
CREATE TABLE IF NOT EXISTS ResetPasswordToken (
  ResetPasswordToken TEXT,
  ExpirationTime TIMESTAMP,
  IsUsed BOOLEAN,
  UserId INTEGER,
  PRIMARY KEY (ResetPasswordToken, UserId),
  CONSTRAINT fk_resetpasswordtoken_user FOREIGN KEY (UserId) REFERENCES "User"(UserId)
);

-- Validation token details for each user
CREATE TABLE IF NOT EXISTS ValidationToken (
  ValidationToken TEXT,
  IsUsed BOOLEAN,
  UserId INTEGER,
  PRIMARY KEY (ValidationToken, UserId),
  CONSTRAINT fk_validationtoken_user FOREIGN KEY (UserId) REFERENCES "User"(UserId)
);

-- Car details owned by a user
CREATE TABLE IF NOT EXISTS Car (
  CarId SERIAL PRIMARY KEY,
  LicensePlate TEXT,
  Make TEXT,
  Model TEXT,
  Color TEXT,
  UserId INTEGER,
  CONSTRAINT fk_car_user FOREIGN KEY (UserId) REFERENCES "User"(UserId)
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
  CONSTRAINT fk_parkingspot_user FOREIGN KEY (UserId) REFERENCES "User"(UserId)
);

-- Master table of bookings with fields that are common to standard and auction bookings
CREATE TABLE IF NOT EXISTS Booking (
  BookingID SERIAL PRIMARY KEY,
  BuyerUserId INTEGER,
  PaidAmount REAL,
  CONSTRAINT fk_booking_user FOREIGN KEY (BuyerUserId) REFERENCES "User"(UserId)
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




--======================================================--
-- INSERTS
--======================================================--

-- Insert mock data for User
INSERT INTO "User" (PasswordHash, FirstName, LastName, Email, IsVerified)
VALUES 
('hashpassword1', 'John', 'Doe', 'john.doe@example.com', TRUE),
('hashpassword2', 'Jane', 'Smith', 'jane.smith@example.com', FALSE),
('hashpassword3', 'Alice', 'Johnson', 'alice.johnson@example.com', TRUE);

-- Insert mock data for Car
INSERT INTO Car (LicensePlate, Make, Model, Color, UserId)
VALUES
('ABC123', 'Toyota', 'Corolla', 'Blue', 1),
('XYZ789', 'Honda', 'Civic', 'Red', 2),
('LMN456', 'Tesla', 'Model 3', 'White', 3);

-- Insert mock data for ParkingSpot
INSERT INTO ParkingSpot (UserId, Address, Longitude, Latitude, HasShelter, HasPlugIn, HasChargingStation)
VALUES
(1, '123 Main St, New York, NY', -74.0060, 40.7128, TRUE, FALSE, TRUE),
(2, '456 Elm St, San Francisco, CA', -122.4194, 37.7749, FALSE, TRUE, FALSE),
(3, '789 Maple Ave, Chicago, IL', -87.6298, 41.8781, TRUE, TRUE, TRUE);

-- Insert mock data for Booking
INSERT INTO Booking (BuyerUserId, PaidAmount)
VALUES
(1, 20.0),
(2, 15.0),
(3, 25.0);

-- Insert mock data for Listing
INSERT INTO Listing (ParkingSpotId, PricePerHour, IsActive)
VALUES
(1, 5.0, TRUE),
(2, 4.0, TRUE),
(3, 6.0, FALSE);

-- Insert mock data for TimeUnit
INSERT INTO TimeUnit (UnitNum, Date, ListingId, BookingId)
VALUES
(1, '2024-09-25', 1, 1),
(2, '2024-09-25', 2, 2),
(3, '2024-09-25', 3, NULL);

-- Insert mock data for StandardBooking
INSERT INTO StandardBooking (ListingId, StartUnitNum, EndUnitNum, Date)
VALUES
(1, 1, 4, '2024-09-26'),
(2, 2, 6, '2024-09-26'),
(3, 3, 8, '2024-09-26');

-- Insert mock data for AuctionTimeSlot
INSERT INTO AuctionTimeSlot (ListingId, CurrentBidId, Date, StartUnitNum, EndUnitNum, BidAmount)
VALUES
(1, 1, '2024-09-27', 1, 4, 12.0),
(2, 2, '2024-09-27', 5, 8, 10.0),
(3, 3, '2024-09-27', 9, 12, 8.0);

-- Insert mock data for AuctionBooking
INSERT INTO AuctionBooking (AuctionTimeSlotId)
VALUES
(1),
(2),
(3);

-- Insert mock data for SessionToken
INSERT INTO SessionToken (SessionToken, Device, ExpirationTime, UserId)
VALUES
('token123', 'iPhone', '2024-09-30 12:00:00', 1),
('token456', 'Android', '2024-09-30 12:00:00', 2),
('token789', 'Web', '2024-09-30 12:00:00', 3);

-- Insert mock data for ResetPasswordToken
INSERT INTO ResetPasswordToken (ResetPasswordToken, ExpirationTime, IsUsed, UserId)
VALUES
('resetToken123', '2024-09-28 12:00:00', FALSE, 1),
('resetToken456', '2024-09-28 12:00:00', TRUE, 2),
('resetToken789', '2024-09-28 12:00:00', FALSE, 3);

-- Insert mock data for ValidationToken
INSERT INTO ValidationToken (ValidationToken, IsUsed, UserId)
VALUES
('validationToken123', FALSE, 1),
('validationToken456', TRUE, 2),
('validationToken789', FALSE, 3);
