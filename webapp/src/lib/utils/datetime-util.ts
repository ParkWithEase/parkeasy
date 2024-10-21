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

export function getNextMonday(d: Date) {
    const next_week_monday = getMonday(d);
    next_week_monday.setDate(next_week_monday.getDate() + 7);
    return next_week_monday;
}

export function getPreviousMonday(d: Date) {
    const prev_week_monday = getMonday(d);
    prev_week_monday.setDate(prev_week_monday.getDate() - 7);
    return prev_week_monday;
}