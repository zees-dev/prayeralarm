package aladhan

import "testing"

// TestGetMonthCalendar is an API integration test
func TestGetMonthCalendar(t *testing.T) {
	t.Run("gets calendar for Auckland NewZealand - successful API status with timezone", func(t *testing.T) {

		got := GetMonthCalendar("Auckland", "NewZealand", "0,0,0,0,0", 1, 1)

		want := MonthlyAdhanCalenderResponse{Code: 200, Status: "OK"}

		if got.Code != want.Code {
			t.Errorf("want %d, got %d", want.Code, got.Code)
		}

		if got.Status != want.Status {
			t.Errorf("want %s, got %s", want.Status, got.Status)
		}

		wantTz := "Pacific/Auckland"

		if got.Data[0].Meta.Timezone != wantTz {
			t.Errorf("want %s, got %s", wantTz, got.Data[0].Meta.Timezone)
		}
	})
}
