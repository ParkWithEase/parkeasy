import { DAY_IN_A_WEEK, TOTAL_SEGMENTS_NUMBER } from "$lib/constants";
import { TimeSlotStatus } from "$lib/enum/timeslot-status";
import { getDateWithDayOffset } from "./datetime-util";


//Need to define type when intergrate with backend
export function extractRelevantTimeSlot(full_time_slot: [], startDate: Date, endDate: Date) {
    let relevant_time_slots = [];
    full_time_slot?.forEach((slot) => {
        if (slot.date >= startDate && slot.date < endDate) {
            relevant_time_slots = [...relevant_time_slots, slot];
        }
    });
    return relevant_time_slots;
}

export function constructAvailabilityTable(
    today: Date,
    weekStart: Date,
    original_slot_data,
    edit_records,
) {
    const availability_table : TimeSlotStatus[][] = Array.from({ length: TOTAL_SEGMENTS_NUMBER }, () =>
        Array(DAY_IN_A_WEEK).fill(0)
    ); 
    const current_seg = today.getUTCHours() * 2 + Math.ceil(today.getUTCMinutes() / 30);
    const cutoff_date = new Date(
        Date.UTC(today.getFullYear(), today.getUTCMonth(), today.getUTCDate())
    );
    for (let day = 0; day < DAY_IN_A_WEEK; day++) {
        const currentDate = new Date(weekStart);
        currentDate.setUTCDate(weekStart.getUTCDate() + day);
        currentDate.setUTCHours(0, 0, 0, 0);
        for (let seg = 0; seg < TOTAL_SEGMENTS_NUMBER; seg++) {
            if (
                currentDate.getTime() < cutoff_date.getTime() ||
                (currentDate.getTime() == cutoff_date.getTime() && seg < current_seg)
            ) {
                availability_table[seg][day] = TimeSlotStatus.PASTDUE;
            } else {
                availability_table[seg][day] = TimeSlotStatus.NONE;
            }
        }
    }

    original_slot_data.forEach((slot) => {
        const day = (slot.date.getUTCDay() || 7) - 1;
        availability_table[slot.segment][day] = slot.status;
    });

    edit_records.forEach((slot) => {
        const day = (slot.date.getUTCDay() || 7) - 1;
        availability_table[slot.segment][day] = slot.status;
    });
    return availability_table;
}

export function getWeekAvailabilityTable(currentTime: Date, weekStart: Date, time_slots_data, edit_records)
{
    const weekEnd = getDateWithDayOffset(weekStart, DAY_IN_A_WEEK);
    const relevantTimeData = extractRelevantTimeSlot(time_slots_data, weekStart, weekEnd);
    const relevateEditRecords = extractRelevantTimeSlot(edit_records, weekStart, weekEnd);
    return constructAvailabilityTable(currentTime, weekStart, relevantTimeData, relevateEditRecords);
}
