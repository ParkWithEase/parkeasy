--======================================================--
-- INSERTS
--======================================================--

-- Insert mock data for User
INSERT INTO Users (PasswordHash, FirstName, LastName, Email, IsVerified)
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
