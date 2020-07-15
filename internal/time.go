package internal

import "time"

var (
	maxTime       = time.Unix(0, (1<<63)-1)
	utcLayout     = "2006-01-02 15:04:05"
	yearlyLayout  = "2006"
	monthlyLayout = "200601"
	dailyLayout   = "20060102"
	hourlyLayout  = "2006010215"
)

// TimestampByMaxTime returns max unix nano.
func TimestampByMaxTime() int64 {
	return maxTime.UnixNano()
}

// StringToTime converts string to time.
func StringToTime(s string) time.Time {
	t, err := time.Parse(utcLayout, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// TimeToString converts time to string.
func TimeToString(t time.Time) string {
	return t.Format(utcLayout)
}

// YearlyStringToTime converts string formatted yearly layout to time.
func YearlyStringToTime(s string) time.Time {
	t, err := time.Parse(yearlyLayout, s)
	if err != nil {
		return time.Time{}
	}
	return t
}


// TimeToYearlyStringFormat converts time to string formatted yearly layout.
func TimeToYearlyStringFormat(t time.Time) string {
	return t.Format(yearlyLayout)
}

// MonthlyStringToTime converts string formatted monthly layout to time.
func MonthlyStringToTime(s string) time.Time {
	t, err := time.Parse(monthlyLayout, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// TimeToMonthlyStringFormat converts time to string formatted monthly layout.
func TimeToMonthlyStringFormat(t time.Time) string {
	return t.Format(monthlyLayout)
}

// DailyStringToTime converts string formatted daily layout to time.
func DailyStringToTime(s string) time.Time {
	t, err := time.Parse(dailyLayout, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// TimeToDailyStringFormat converts time to string formatted daily layout
func TimeToDailyStringFormat(t time.Time) string {
	return t.Format(dailyLayout)
}


// HourlyStringToTime converts string formatted hourly layout to time.
func HourlyStringToTime(s string) time.Time {
	t, err := time.Parse(hourlyLayout, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// TimeToHourlyStringFormat converts time to string formatted hourly layout
func TimeToHourlyStringFormat(t time.Time) string {
	return t.Format(hourlyLayout)
}
