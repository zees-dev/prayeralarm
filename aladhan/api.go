package aladhan

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// GetMonthCalendar calls adhan API and returns serialized `MonthlyAdhanCalenderResponse` object from JSON response
// API Endpoint: https://aladhan.com/prayer-times-api#GetCalendarByCitys
// API Adhan Timing Tuning: https://aladhan.com/calculation-methods
// Example request: `curl 'http://api.aladhan.com/v1/calendarByCity?city=Auckland&country=NewZealand&method=3&month=12&year=2020&tune=0,0,0,0,0,0,0,0'`
func GetMonthCalendar(city string, country string, offsets string, month time.Month, year int) MonthlyAdhanCalenderResponse {
	offsetSlice := strings.Split(offsets, ",")
	fajr, dhuhr, asr, maghrib, isha := offsetSlice[0], offsetSlice[1], offsetSlice[2], offsetSlice[3], offsetSlice[4]
	// Tune order: Imsak,Fajr,Sunrise,Dhuhr,Asr,Maghrib,Sunset,Isha,Midnight
	tuneListStr := fmt.Sprintf("0,%s,0,%s,%s,%s,0,%s", fajr, dhuhr, asr, maghrib, isha)
	url := fmt.Sprintf(
		"http://api.aladhan.com/v1/calendarByCity?city=%s&country=%s&method=3&month=%d&year=%d&tune=%s",
		city,
		country,
		month,
		year,
		tuneListStr,
	)

	log.Printf("Calling API: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("API request to URL %s failed", url)
	}
	defer resp.Body.Close()

	var monthlyCalendarResp MonthlyAdhanCalenderResponse
	if err := json.NewDecoder(resp.Body).Decode(&monthlyCalendarResp); err != nil {
		log.Fatalf("Failed to decode URL response, incorrect struct formatting and/or field type(s)")
	}

	// Remove non-main adhans
	for _, timings := range monthlyCalendarResp.Data {
		for adhan := range timings.Timings {
			switch adhan {
			case Fajr, Dhuhr, Asr, Maghrib, Isha:
			default:
				delete(timings.Timings, adhan)
			}
		}
	}

	return monthlyCalendarResp
}
