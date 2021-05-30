package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/zees-dev/prayeralarm/aladhan"
	server "github.com/zees-dev/prayeralarm/http"
	"github.com/zees-dev/prayeralarm/prayer"
)

type cliFlags struct {
	city    string
	country string
	offset  string
	month   time.Month
	year    int
	port    uint
}

func runPrayerAlarm(flags cliFlags, adhanService *prayer.Service) {
	year, month := flags.year, flags.month
	for {
		monthCalendar := aladhan.GetMonthCalendar(flags.city, flags.country, flags.offset, month, year)

		adhanService.GeneratePrayers(monthCalendar)
		adhanService.DisplayPrayerTimings(os.Stdout)
		adhanService.ExecutePrayers()

		year, month, _ = time.Now().AddDate(0, 1, 0).Date()
	}
}

func main() {
	year, month, _ := time.Now().Date()

	cityPtr := flag.String("city", "Auckland", "city for which adhan timings are to be retrieved")
	countryPtr := flag.String("country", "NewZealand", "country for which adhan timings are to be retrieved")
	offsetPtr := flag.String("offsets", "0,0,0,0,0", "comma seperated string of adhan offsets (in mins) for the 5 daily adhans (fajr, dhuhr, asr, maghrib, isha)")
	yearPtr := flag.Int("year", year, "year of adhan playback")
	monthPtr := flag.Int("month", int(month), "month of adhan playback")
	portPtr := flag.Uint("port", 8000, "server port")

	flag.Parse()

	cliFlags := cliFlags{
		city:    *cityPtr,
		country: *countryPtr,
		offset:  *offsetPtr,
		year:    *yearPtr,
		month:   time.Month(*monthPtr),
		port:    *portPtr,
	}

	log.Printf(
		"Flags - city: %s, country: %s, offsets: %s, year: %d, month: %d, port: %d",
		cliFlags.city,
		cliFlags.country,
		cliFlags.offset,
		cliFlags.year,
		cliFlags.month,
		cliFlags.port,
	)

	player := prayer.NewStdOutPlayer()
	adhanService := prayer.NewService(player)

	go runPrayerAlarm(cliFlags, adhanService)

	handler := server.NewHandler(adhanService)
	server := &http.Server{
		Handler:      handler.Router,
		Addr:         fmt.Sprintf("127.0.0.1:%d", *portPtr),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Running prayeralarm server on port %d...", cliFlags.port)
	log.Fatal(server.ListenAndServe())
}
