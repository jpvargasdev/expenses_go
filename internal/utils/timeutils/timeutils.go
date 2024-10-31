package timeutils 

import "time"

func CalculatePeriodBoundaries(date time.Time) (int64, int64) {
  year, month, day := date.Date()
  location := date.Location()

  if day < 25 {
    month = month - 1
    if month == 0 {
      month = 12
      year--
    }
  }

  start := time.Date(year, month, 25, 0, 0, 0, 0, location).Unix()
  end := time.Date(year, month+1, 24, 23, 59, 59, 0, location).Unix()

  return start, end
}
