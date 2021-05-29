package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/zees-dev/prayeralarm/aladhan"
	"github.com/zees-dev/prayeralarm/prayer"
)

func main() {
	year, month, _ := time.Now().Date()

	cityPtr := flag.String("city", "Auckland", "city for which adhan timings are to be retrieved")
	countryPtr := flag.String("country", "NewZealand", "country for which adhan timings are to be retrieved")
	offsetPtr := flag.String("offsets", "0,0,0,0,0", "comma seperated string of adhan offsets (in mins) for the 5 daily adhans (fajr, dhuhr, asr, maghrib, isha)")
	yearPtr := flag.Int("year", year, "year of adhan playback")
	monthPtr := flag.Int("month", int(month), "month of adhan playback")
	flag.Parse()

	year, month = *yearPtr, time.Month(*monthPtr)
	log.Printf("Flags - City: %s, Country: %s, Offsets: %s, Year: %d, Month: %d", *cityPtr, *countryPtr, *offsetPtr, year, month)

	player := prayer.NewOmxPlayer()
	adhanService := prayer.NewService(player)

	for {
		monthCalendar := aladhan.GetMonthCalendar(*cityPtr, *countryPtr, *offsetPtr, month, year)
		adhanTimings := aladhan.ExtractAdhanTimings(monthCalendar)

		adhanService.SetAdhanTimings(adhanTimings)
		adhanService.DisplayAdhanTimings(os.Stdout)
		adhanService.ExecuteAdhan()

		year, month, _ = time.Now().AddDate(0, 1, 0).Date()
	}
}
