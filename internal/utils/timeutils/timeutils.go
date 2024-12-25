package timeutils

import (
	"strconv"
	"time"
)

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

func GetSalaryMonthRange(days ...string) (startDate time.Time, endDate time.Time) {
	var startDay, endDay int

	if days[0] != "" {
		startDay, _ = strconv.Atoi(days[0])
	} else {
		startDay = 25
	}

	if days[1] != "" {
		endDay, _ = strconv.Atoi(days[1])
	} else {
		endDay = 24
	}

	now := time.Now()

	if now.Day() >= startDay {
		// Current period: salaryDay of this month to endMonthDay of the next month
		startDate = time.Date(now.Year(), now.Month(), startDay, 0, 0, 0, 0, now.Location())
		endDate = time.Date(now.Year(), now.Month()+1, endDay, 23, 59, 59, 999999999, now.Location())
	} else {
		// Previous period: salaryDay of the previous month to endMonthDay of this month
		startDate = time.Date(now.Year(), now.Month()-1, startDay, 0, 0, 0, 0, now.Location())
		endDate = time.Date(now.Year(), now.Month(), endDay, 23, 59, 59, 999999999, now.Location())
	}

	return startDate, endDate
}
