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
                    const start_time = startTime.toISOString();
                    let end_time = endTime.toISOString();

                    //Handle daylight saving time by ensuring that the endtime is of the same zone as the startTime
                    if (startTime.getTimezoneOffset() < endTime.getTimezoneOffset()) {
                        end_time = manualUTCConvert(endTime, -1);
                    } else if (startTime.getTimezoneOffset() > endTime.getTimezoneOffset()) {
                        end_time = manualUTCConvert(startTime, 1);
                    }
                    acc.push({
                        start_time: start_time,
                        end_time: end_time
                    });
                    return acc;
                }, new Array<TimeUnit>())
            )
        );
    }
    console.log(result);
    return result;
}
function manualUTCConvert(date: Date, offset: number) {
    return (
        date.getUTCFullYear() +
        '-' +
        `${date.getUTCMonth() + 1}`.padStart(2, '0') +
        '-' +
        `${date.getUTCDate()}`.padStart(2, '0') +
        'T' +
        `${date.getUTCHours() + offset}`.padStart(2, '0') +
        ':' +
        `${date.getUTCMinutes()}`.padStart(2, '0') +
        ':' +
        `${date.getUTCSeconds()}`.padStart(2, '0') +
        '.' +
        String((date.getUTCMilliseconds() / 1000).toFixed(3)).slice(2, 5) +
        'Z'
    );
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
            }
            for (let segment = 0; segment < total_segment_number; segment++) {
                availability[segment] ??= [];
                availability[segment][(day.getDay() || 7) - 1] = TimeSlotStatus.PASTDUE;
            }
        }
    }
}

export function createEmptyTable(): TimeSlotStatus[][] {
    return Array.from({ length: TOTAL_SEGMENTS_NUMBER }, () =>
        Array(DAY_IN_A_WEEK).fill(TimeSlotStatus.NONE)
    );
}
