package main

import (
	"flag"
	"log"
	"time"

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
	output  string
}

func main() {
	year, month, _ := time.Now().Date()

	cityPtr := flag.String("city", "Auckland", "city for which adhan timings are to be retrieved")
	countryPtr := flag.String("country", "NewZealand", "country for which adhan timings are to be retrieved")
	offsetPtr := flag.String("offsets", "0,0,0,0,0", "comma seperated string of adhan offsets (in mins) for the 5 daily adhans (fajr, dhuhr, asr, maghrib, isha)")
	yearPtr := flag.Int("year", year, "year of adhan playback")
	monthPtr := flag.Int("month", int(month), "month of adhan playback")
	outputPtr := flag.String("output", string(prayer.OMX), "output device; supported options are `stdout`, `native` and `omx`")
	portPtr := flag.Uint("port", 8080, "server port")

	flag.Parse()

	cliFlags := cliFlags{
		city:    *cityPtr,
		country: *countryPtr,
		offset:  *offsetPtr,
		year:    *yearPtr,
		month:   time.Month(*monthPtr),
		output:  *outputPtr,
		port:    *portPtr,
	}

	log.Printf(
		"Flags - city: %s, country: %s, offsets: %s, year: %d, month: %d, output: %s, port: %d",
		cliFlags.city,
		cliFlags.country,
		cliFlags.offset,
		cliFlags.year,
		cliFlags.month,
		cliFlags.output,
		cliFlags.port,
	)

	player, err := prayer.GetPlayer(prayer.Output(cliFlags.output))
	if err != nil {
		log.Fatalln(err)
	}

	adhanService := prayer.NewService(player)
	adhanService.InitialisePrayeralarm(cliFlags.year, cliFlags.month, cliFlags.city, cliFlags.country, cliFlags.offset)

	server := server.NewServer(adhanService)
	server.Run(cliFlags.port)
}
