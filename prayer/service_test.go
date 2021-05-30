package prayer

import (
	"testing"
	"time"
)

func TestGetTime(t *testing.T) {
	t.Run("test date time with timezone auckland", func(t *testing.T) {
		timeStr := "01 Jan 2021 04:14 (NZDT)"

		tz := "Pacific/Auckland"
		got := getTime(timeStr, tz)

		l, _ := time.LoadLocation(tz)
		want := time.Date(2021, 1, 1, 4, 14, 0, 0, l)

		if got.String() != want.String() {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("test date time with timezone new york", func(t *testing.T) {
		timeStr := "01 Jan 2021 04:14 (EST)"

		tz := "America/New_York"
		got := getTime(timeStr, tz)

		l, _ := time.LoadLocation(tz)
		want := time.Date(2021, 1, 1, 4, 14, 0, 0, l)

		if got.String() != want.String() {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("test date time with timezone perth", func(t *testing.T) {
		timeStr := "01 Jan 2021 04:14 (AWST)"

		lc := "Australia/Perth"
		got := getTime(timeStr, lc)

		l, _ := time.LoadLocation(lc)
		want := time.Date(2021, 1, 1, 4, 14, 0, 0, l)

		if got.String() != want.String() {
			t.Errorf("want %s, got %s", want, got)
		}
	})
}
