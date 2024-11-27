package timeutils

import "time"

func CalculatePeriodBoundaries(date time.Time, offset ...int) (int64, int64) {
	mOffset := 1
	if len(offset) > 0 {
		mOffset = offset[0]
	}
	year, month, day := date.Date()
	location := date.Location()

	start := time.Date(year, month, day+mOffset, 0, 0, 0, 0, location).Unix()
	end := time.Date(year, month+1, day+mOffset+1, 23, 59, 59, 0, location).Unix()

	return start, end
}

func GetSalaryMonthRange(salaryDay ...int) (startDate time.Time, endDate time.Time) {
	day := 25

	if len(salaryDay) > 0 {
		day = salaryDay[0]
	}

	now := time.Now()

	// Determine if we are in the period after or before the 25th
	if now.Day() >= day {
		// Current period: 25th of current month to 24th of next month
		startDate = time.Date(now.Year(), now.Month(), day, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 1, -1)
	} else {
		// Previous period: 25th of previous month to 24th of current month
		startDate = time.Date(now.Year(), now.Month(), day, 0, 0, 0, 0, now.Location()).AddDate(0, -1, 0)
		endDate = startDate.AddDate(0, 1, -1)
	}

	// Set end date time to the end of the day
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	return startDate, endDate
}
