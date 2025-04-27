package service

import "time"

func IsValidDate(dateStr string) bool {
	layout := "2006-01-02"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		return false
	}
	now := time.Now()
	minDate := now.AddDate(0, 0, -90)
	return !date.Before(minDate) && !date.After(now)
}

// DaysBetween calculates the number of days between two dates.
func DaysBetween(start, end time.Time) int {
	if start.After(end) {
		start, end = end, start
	}
	diff := end.Sub(start)
	days := int(diff.Hours() / 24)

	return days
}
