//Return the monday of the week corresponding to the input date
export function getMonday(d: Date) {
    const this_week_monday = new Date(d);
    this_week_monday.setSeconds(0);
    this_week_monday.setMinutes(0);
    this_week_monday.setHours(0);
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