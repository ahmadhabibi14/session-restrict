package helper

import (
	"errors"
	"fmt"
	"time"
)

// GetCurrentDateTime returns the current date and time in a specified format.
func GetCurrentDateTime(format string) string {
	return time.Now().Format(format)
}

// GetCurrentDate returns the current date in YYYY-MM-DD format.
func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

// GetCurrentTime returns the current time in HH:mm:ss format.
func GetCurrentTime() string {
	return time.Now().Format("15:04:05")
}

// FormatDateTime formats a given time.Time to a specified format.
func FormatDateTime(t time.Time, format string) string {
	return t.Format(format)
}

// ParseDate parses a date string in a specific format and returns a time.Time object.
func ParseDate(dateStr, format string) (time.Time, error) {
	parsedTime, err := time.Parse(format, dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}

// DaysBetween calculates the number of days between two dates.
func DaysBetween(startDate, endDate string, format string) (int, error) {
	start, err1 := time.Parse(format, startDate)
	end, err2 := time.Parse(format, endDate)

	if err1 != nil || err2 != nil {
		return 0, errors.New("invalid date format")
	}

	duration := end.Sub(start)
	return int(duration.Hours() / 24), nil
}

// AddDaysToDate adds a specific number of days to a given date.
func AddDaysToDate(dateStr string, days int, format string) (string, error) {
	date, err := time.Parse(format, dateStr)
	if err != nil {
		return "", err
	}
	newDate := date.AddDate(0, 0, days)
	return newDate.Format(format), nil
}

const (
	DateFormatYYYYMM = `2006-01`
)

// IsValidDate checks if a date string is valid based on a specific format.
func IsValidDate(dateStr, format string) bool {
	_, err := time.Parse(format, dateStr)
	return err == nil
}

// IsDateGreater checks if a base date is greater than a target date.
func IsDateGreater(baseDate, targetDate, format string) bool {
	if IsValidDate(baseDate, format) && IsValidDate(targetDate, format) {
		return baseDate > targetDate
	}

	return false
}

// GetWeekStartAndEnd returns the start and end of the current week.
func GetWeekStartAndEnd() (time.Time, time.Time) {
	now := time.Now()
	weekday := int(now.Weekday())
	startOfWeek := now.AddDate(0, 0, -weekday)
	endOfWeek := startOfWeek.AddDate(0, 0, 6)
	return startOfWeek, endOfWeek
}

// GetMonthStartAndEnd returns the start and end of the current month.
func GetMonthStartAndEnd() (time.Time, time.Time) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)
	return startOfMonth, endOfMonth
}

func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d)/float64(time.Millisecond))
	}
	return fmt.Sprintf("%.2fs", float64(d)/float64(time.Second))
}
