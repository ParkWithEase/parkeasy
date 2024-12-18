//Return the monday of the week corresponding to the input date
export function getMonday(d: Date) {
    const this_week_monday = new Date(d);
    this_week_monday.setSeconds(0);
    this_week_monday.setMinutes(0);
    this_week_monday.setHours(0);
    const day = this_week_monday.getDay() || 7;
    const diff = this_week_monday.getDate() - day + 1;
    this_week_monday.setDate(diff);
    return this_week_monday;
}

//Return a date that is <offset> number days from the input date
export function getDateWithDayOffset(d: Date, offset: number) {
    const result = new Date(d);
    result.setDate(result.getDate() + offset);
    return result;
}

export function convertDateToUTC(currentDate: Date) {
    const utcDate = new Date(
        Date.UTC(
            currentDate.getUTCFullYear(),
            currentDate.getUTCMonth(),
            currentDate.getUTCDate(),
            currentDate.getUTCHours(),
            currentDate.getUTCMinutes(),
            currentDate.getUTCSeconds()
        )
    );
    return utcDate;
}
