package timeutils

import "time"

func CalculatePeriodBoundaries(date time.Time, initDay ...int) (int64, int64) {
	startDay := 25
	if len(initDay) > 0 {
		startDay = initDay[0]
	}
	year, month, day := date.Date()
	location := date.Location()

	if day < 25 {
		month = month - 1
		if month == 0 {
			month = 12
			year--
		}
	}

	start := time.Date(year, month, startDay, 0, 0, 0, 0, location).Unix()
	end := time.Date(year, month+1, startDay-1, 23, 59, 59, 0, location).Unix()

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
