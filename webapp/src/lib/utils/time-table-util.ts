import { DAY_IN_A_WEEK, TOTAL_SEGMENTS_NUMBER } from '$lib/constants';
import { TimeSlotStatus } from '$lib/enum/timeslot-status';
import type { components } from '$lib/sdk/schema';

type TimeUnit = components['schemas']['TimeUnit'];

export function toTimeUnits(map: Map<number, TimeSlotStatus[][]>): TimeUnit[] {
    const result = new Array<TimeUnit>();
    for (const [week, value] of map.entries()) {
        const weekDate = new Date(week);
        result.push(
            ...value.flatMap((infoday, segment) =>
                infoday.reduce((acc, status, day) => {
                    if (status !== TimeSlotStatus.AVAILABLE) {
                        return acc;
                    }
                    const startTime = new Date(weekDate);
                    startTime.setDate(weekDate.getDate() + day);
                    startTime.setHours(0, segment * 30, 0, 0);
                    const endTime = new Date(startTime);
                    endTime.setMinutes(endTime.getMinutes() + 30);
                    acc.push({
                        start_time: startTime.toISOString(),
                        end_time: endTime.toISOString()
                    });
                    return acc;
                }, new Array<TimeUnit>())
            )
        );
    }
    return result;
}

export function initializeAvailabilityTables(
    availabilityTables: Map<number, TimeSlotStatus[][]>,
    current_monday: Date,
    today: Date
) {
    for (
        let week = new Date(current_monday);
        week.getTime() < today.getTime();
        week.setDate(week.getDate() + 7)
    ) {
        const availability: TimeSlotStatus[][] = new Array(TOTAL_SEGMENTS_NUMBER);
        availabilityTables.set(week.getTime(), availability);

        for (
            let day = new Date(week);
            day.getTime() < today.getTime();
            day.setDate(day.getDate() + 1)
        ) {
            let total_segment_number = TOTAL_SEGMENTS_NUMBER;
            const midnightToday = new Date(today);
            midnightToday.setHours(0, 0, 0, 0);
            if (day.getTime() >= midnightToday.getTime()) {
                total_segment_number = today.getHours() * 2 + Math.ceil(today.getMinutes() / 30);
                console.log(total_segment_number);
            }
            for (let segment = 0; segment < total_segment_number; segment++) {
                availability[segment] ??= [];
                availability[segment][day.getDate() - week.getDate()] = TimeSlotStatus.PASTDUE;
            }
        }
    }
}

export function createEmptyTable(): TimeSlotStatus[][] {
    return Array.from({ length: TOTAL_SEGMENTS_NUMBER }, () =>
        Array(DAY_IN_A_WEEK).fill(TimeSlotStatus.NONE)
    );
}
