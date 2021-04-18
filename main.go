package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

// getTime converts string input with tz location to time object
// https://yourbasic.org/golang/format-parse-string-time-date-example/
func getTime(timeStr string, location string) time.Time {
	dateFormat := "02 Jan 2006 15:04 (MST)"

	tl, err := time.LoadLocation(location)
	if err != nil {
		log.Fatalf(`Incorrect location input: "%s"`, location)
	}

	t, err := time.ParseInLocation(dateFormat, timeStr, tl)
	if err != nil {
		log.Fatalf(`Incorrect date-time input: "%s"`, timeStr)
	}
	return t
}

// getMonthCalendar calls adhan API and returns serialized `MonthlyAdhanCalenderResponse` object from JSON response
// API Endpoint: https://aladhan.com/prayer-times-api#GetCalendarByCitys
// API Adhan Timing Tuning: https://aladhan.com/calculation-methods
// Example request: `curl 'http://api.aladhan.com/v1/calendarByCity?city=Auckland&country=NewZealand&method=3&month=12&year=2020&tune=0,0,0,0,0,0,0,0'`
func getMonthCalendar(city string, country string, offsets string, month time.Month, year int) MonthlyAdhanCalenderResponse {
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

func executeMonthlyCalendar(monthCalendar MonthlyAdhanCalenderResponse, w io.Writer) {
	adhanTimings := []AdhanTime{}

	// Get all adhan timings after current time for remaining days of the month
	currentTime := time.Now()
	for _, timings := range monthCalendar.Data {
		for adhan, timeStr := range timings.Timings {
			fullTimeStr := fmt.Sprintf("%s %s", timings.Date.Readable, timeStr)
			adhanTime := getTime(fullTimeStr, timings.Meta.Timezone)
			if adhanTime.After(currentTime) {
				adhanTimings = append(adhanTimings, AdhanTime{Type: adhan, Time: adhanTime})
			}
		}
	}

	// Sort upcoming adhans by time
	sort.Slice(adhanTimings[:], func(i, j int) bool {
		return adhanTimings[i].Time.Before(adhanTimings[j].Time)
	})

	// Display calendar in console
	displayCalendar(adhanTimings)

	// Play the adhan at the correct times - from current time
	for _, adhanTiming := range adhanTimings {
		timeTillNextAdhan := adhanTiming.Time.Sub(time.Now())
		log.Printf("Waiting %s for %s adhan...", timeTillNextAdhan, adhanTiming.Type)
		if _, err := w.Write([]byte(adhanTiming.Type)); err != nil {
			log.Fatalln(err)
		}

		time.Sleep(timeTillNextAdhan)

		log.Printf("Running %s adhan at %s...", adhanTiming.Type, adhanTiming.Time)

		if _, err := w.Write([]byte(adhanTiming.Type)); err != nil {
			log.Fatalln(err)
		}
	}
}

func main() {
	year, month, _ := time.Now().Date()

	cityPtr := flag.String("city", "Auckland", "city for which adhan timings are to be retrieved")
	countryPtr := flag.String("country", "NewZealand", "country for which adhan timings are to be retrieved")
	offsetPtr := flag.String("offsets", "0,0,0,0,0", "comma seperated string of adhan offsets (in mins) for the 5 daily adhans (fajr, dhuhr, asr, maghrib, isha)")
	yearPtr := flag.Int("year", year, "year of adhan playback")
	monthPtr := flag.Int("month", int(month), "month of adhan playback")
	flag.Parse()

	log.Printf("Flags - City: %s, Country: %s, Offsets: %s, Year: %d, Month: %d", *cityPtr, *countryPtr, *offsetPtr, *yearPtr, *monthPtr)
	for {
		year, month = *yearPtr, time.Month(*monthPtr)

		monthCalendar := getMonthCalendar(*cityPtr, *countryPtr, *offsetPtr, month, year)
		executeMonthlyCalendar(monthCalendar, omxPlayer{})

		year, month, _ = time.Now().AddDate(0, 1, 0).Date()
	}
}
