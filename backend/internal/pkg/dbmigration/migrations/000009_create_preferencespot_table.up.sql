CREATE TABLE IF NOT EXISTS PreferenceSpots (
    UserID BIGINT NOT NULL REFERENCES Users(UserID),
    ParkingSpotID BIGINT NOT NULL REFERENCES ParkingSpot(ParkingSpotID),
    PRIMARY KEY (UserID, ParkingSpotID)
);

