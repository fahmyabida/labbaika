// YYYY-MM-DD from LOCAL date parts — not toISOString(), which converts to UTC
// and can shift the date back a day in timezones ahead of UTC (e.g. WIB).
export const toLocalDateString = (d: Date) =>
    `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;

export const currentMonthRange = () => {
    const now = new Date();
    const start = new Date(now.getFullYear(), now.getMonth(), 1);
    const end = new Date(now.getFullYear(), now.getMonth() + 1, 0); // day 0 of next month = last day of this month
    return { start: toLocalDateString(start), end: toLocalDateString(end) };
};
