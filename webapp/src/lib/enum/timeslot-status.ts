export enum TimeSlotStatus {
    NONE,
    AVAILABLE,
    BOOKED,
    AUCTIONED,
    PASTDUE,
    BOOK_INTENT,
    EXPIRED_BOOK,
    ACTIVE_BOOK
}

export const TimeSlotStatusConverter = {
    booked: TimeSlotStatus.BOOKED,
    available: TimeSlotStatus.AVAILABLE,
    none: TimeSlotStatus.NONE
};
