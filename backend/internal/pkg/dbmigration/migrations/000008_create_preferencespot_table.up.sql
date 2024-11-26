CREATE TABLE IF NOT EXISTS PreferenceSpot (
    PreferenceSpotID BIGSERIAL PRIMARY KEY,
    UserID BIGINT NOT NULL REFERENCES Users(UserId),
    ParkingSpotId BIGINT NOT NULL REFERENCES ParkingSpot(ParkingSpotId)
);

CREATE UNIQUE INDEX IF NOT EXISTS PreferenceSpotIdx ON PreferenceSpot(UserID, ParkingSpotId);
