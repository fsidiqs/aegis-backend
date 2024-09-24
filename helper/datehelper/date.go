package datehelper

import (
	"time"
)

func Date(value string) (time.Time, error) {
	layoutFormat := "2006-01-02"
	return time.Parse(layoutFormat, value)
}

func Truncate(t time.Time) time.Time {
	return t.Truncate(24 * time.Hour)
}

func GetDateSince(t time.Time) int {

	return int(time.Since(
		Truncate(
			t,
		),
	).Hours() / 24)
}

func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
