//Return the monday of the week corresponding to the input date
export function getMonday(d: Date) {
    const this_week_monday = new Date(d);
    this_week_monday.setUTCSeconds(0);
    this_week_monday.setUTCMinutes(0);
    this_week_monday.setUTCHours(0);
    const day = this_week_monday.getUTCDay() || 7;
    const diff = this_week_monday.getUTCDate() - day + 1;
    this_week_monday.setUTCDate(diff);
    return this_week_monday;
}

//Return a date that is <offset> number days from the input date
export function getDateWithDayOffset(d: Date, offset: number) {
    const result = new Date(d);
    result.setUTCDate(result.getUTCDate() + offset);
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
