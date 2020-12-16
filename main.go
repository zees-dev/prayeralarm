package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/zees-dev/prayeralarm/models"
)

// getTime converts string input to time object
func getTime(timeStr string) time.Time {
	dateFormat := "02 Jan 2006 15:04 (NZDT)"
	t, err := time.ParseInLocation(dateFormat, timeStr, time.Now().Location())
	if err != nil {
		log.Fatalf("Incorrect date-time input: %s", timeStr)
	}
	return t
}

// printRemainingCalendar renders upcoming calendar in ASCII table
// https://github.com/olekukonko/tablewriter#example-6----identical-cells-merging
func printRemainingCalendar(adhanSlice []models.AdhanTime) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "Adhan", "Time"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	for _, adhan := range adhanSlice {
		year, month, day := adhan.Time.Date()
		dateStr := fmt.Sprintf("%s %d-%s-%d", adhan.Time.Weekday(), day, month, year)
		table.Append([]string{dateStr, string(adhan.Type), adhan.Time.Format("03:04:05 PM")})
	}
	table.Render()
}

// getMonthCalendar calls adhan API and returns serialized `MonthlyAdhanCalenderResponse` object from JSON response
// API Endpoint: https://aladhan.com/prayer-times-api#GetCalendarByCitys
// API Adhan Timing Tuning: https://aladhan.com/calculation-methods
func getMonthCalendar(city string, country string, offsets string, month time.Month, year int) *models.MonthlyAdhanCalenderResponse {
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

	var monthlyCalendarResp models.MonthlyAdhanCalenderResponse
	if err := json.NewDecoder(resp.Body).Decode(&monthlyCalendarResp); err != nil {
		log.Fatalf("Failed to decode URL response, incorrect struct formatting and/or field type(s)")
	}

	// Remove non-main adhans
	for _, timings := range monthlyCalendarResp.Data {
		for adhan := range timings.Timings {
			switch adhan {
			case models.Fajr, models.Dhuhr, models.Asr, models.Maghrib, models.Isha:
			default:
				delete(timings.Timings, adhan)
			}
		}
	}

	return &monthlyCalendarResp
}

func runMonthlyTask(city string, country string, offsets string, year int, month time.Month) {
	monthCalendar := getMonthCalendar(city, country, offsets, month, year)
	adhanTimings := []models.AdhanTime{}

	// Get all adhan timings after current time for remaining days of the month
	currentTime := time.Now()
	for _, timings := range monthCalendar.Data {
		for adhan, timeStr := range timings.Timings {
			fullTimeStr := fmt.Sprintf("%s %s", timings.Date.Readable, timeStr)
			adhanTime := getTime(fullTimeStr)
			if adhanTime.After(currentTime) {
				adhanTimings = append(adhanTimings, models.AdhanTime{adhan, adhanTime})
			}
		}
	}

	// Sort upcoming adhans by time
	sort.Slice(adhanTimings[:], func(i, j int) bool {
		return adhanTimings[i].Time.Before(adhanTimings[j].Time)
	})

	printRemainingCalendar(adhanTimings)

	// Play the adhan at the correct times - from current time
	for _, adhanTiming := range adhanTimings {
		timeTillNextAdhan := adhanTiming.Time.Sub(time.Now())
		log.Printf("Waiting %s for %s adhan...", timeTillNextAdhan, adhanTiming.Type)
		time.Sleep(timeTillNextAdhan)

		log.Printf("Running %s adhan at %s...", adhanTiming.Type, adhanTiming.Time)
		// TODO write to sound buffer
	}

	// Recursively run for next month
	year, month, _ = time.Now().AddDate(0, 1, 0).Date()
	runMonthlyTask(city, country, offsets, year, month)
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

	runMonthlyTask(*cityPtr, *countryPtr, *offsetPtr, *yearPtr, time.Month(*monthPtr))
}
