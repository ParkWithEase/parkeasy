export enum TimeSlotStatus {
    NONE,
    AVAILABLE,
    BOOKED,
    AUCTIONED,
    PASTDUE
}

export const TimeSlotStatusConverter = {
    booked: TimeSlotStatus.BOOKED,
    available: TimeSlotStatus.AVAILABLE,
    none: TimeSlotStatus.NONE
};
