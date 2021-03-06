package prayer

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/zees-dev/prayeralarm/aladhan"
)

type Prayer struct {
	Play  bool          `json:"play"`
	Type  aladhan.Adhan `json:"type"`
	Time  time.Time     `json:"time"`
	Index uint8         `json:"index"`
}

type DailyPrayerTimings struct {
	Date    time.Time `json:"date"`
	Prayers []Prayer  `json:"prayers"`
}

type Service struct {
	mutex              sync.RWMutex
	player             Player
	DailyPrayerTimings []DailyPrayerTimings
}

// NewService returns new adhan service that utilizes player to output adhan audio
func NewService(player Player) *Service {
	return &Service{
		player:             player,
		DailyPrayerTimings: []DailyPrayerTimings{},
	}
}

// InitialisePrayeralarm will initialise the service to run monthly prayer calls.
// The service will populate the prayer adhan timings on a monthly basis, then loop
// through all the prayers of the month (incrementally) to play the adhan at the specified
// prayer time to the provided player.
func (svc *Service) InitialisePrayeralarm(year int, month time.Month, city, country, offset string) {
	log.Println("running prayeralarm service...")
	go func() {
		for {
			monthCalendar := aladhan.GetMonthCalendar(city, country, offset, month, year)

			svc.generatePrayers(monthCalendar)
			svc.DisplayPrayerTimings(os.Stdout)
			svc.executePrayers()

			year, month, _ = time.Now().AddDate(0, 1, 0).Date()
		}
	}()
}

// generatePrayers extracts the monthly adhan timings from the calendar api response
func (svc *Service) generatePrayers(monthCalendar aladhan.MonthlyAdhanCalenderResponse) {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	svc.DailyPrayerTimings = make([]DailyPrayerTimings, 0)

	// Get all adhan timings after current time for remaining days of the month
	currentTime := time.Now()
	prayerIndex := uint8(0)
	for _, timings := range monthCalendar.Data {
		dailyPrayers := make([]Prayer, 0)
		for adhan, timeStr := range timings.Timings {
			fullTimeStr := fmt.Sprintf("%s %s", timings.Date.Readable, timeStr)
			adhanTime := getTime(fullTimeStr, timings.Meta.Timezone)
			if adhanTime.After(currentTime) {
				prayer := Prayer{Play: true, Type: adhan, Time: adhanTime, Index: prayerIndex}
				dailyPrayers = append(dailyPrayers, prayer)
				prayerIndex++
			}
		}
		if len(dailyPrayers) > 0 {
			dateTime, err := getDateFromTimestamp(timings.Date.Timestamp)
			if err != nil {
				log.Fatalln(err)
			}
			svc.DailyPrayerTimings = append(svc.DailyPrayerTimings, DailyPrayerTimings{
				Date:    dateTime,
				Prayers: dailyPrayers,
			})
		}
	}

	for _, dpt := range svc.DailyPrayerTimings {
		// Sort upcoming daily adhans by time
		sort.Slice(dpt.Prayers[:], func(i, j int) bool {
			return dpt.Prayers[i].Time.Before(dpt.Prayers[j].Time)
		})
	}
}

// DisplayPrayerTimings renders upcoming calendar in ASCII table
// https://github.com/olekukonko/tablewriter#example-6----identical-cells-merging
func (svc *Service) DisplayPrayerTimings(writer io.Writer) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Date", "Adhan", "Time", "Play"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	for _, dpt := range svc.DailyPrayerTimings {
		for _, p := range dpt.Prayers {
			year, month, day := p.Time.Date()
			dateStr := fmt.Sprintf("%s %d-%s-%d", p.Time.Weekday(), day, month, year)
			var playStr string
			if p.Play {
				playStr = "Yes"
			} else {
				playStr = "No"
			}
			table.Append([]string{dateStr, string(p.Type), p.Time.Format("03:04:05 PM"), playStr})
		}
	}
	table.Render()
}

// executePrayers plays prayer adhan based on adhan timings if execution of the respective prayer is set to true
func (svc *Service) executePrayers() {
	// Play the adhan at the correct times - from current time
	for _, dpt := range svc.DailyPrayerTimings {
		for _, p := range dpt.Prayers {
			timeTillNextAdhan := time.Until(p.Time)
			log.Printf("Waiting %s for %s adhan...", timeTillNextAdhan, p.Type)

			time.Sleep(timeTillNextAdhan)

			// Only play adhan if its set to execute
			if p.Play {
				log.Printf("Playing %s adhan at %s...", p.Type, p.Time)
				if err := svc.player.Play(p.Type); err != nil {
					log.Fatalln(err)
				}
			} else {
				log.Printf("Skipping %s adhan at %s since execution is set to false", p.Type, p.Time)
			}
		}
	}
}

// TurnOffAllAdhan sets adhan executions for all adhan timings of the month to be muted
func (svc *Service) TurnOffAllAdhan() []DailyPrayerTimings {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	for _, dpt := range svc.DailyPrayerTimings {
		for i := range dpt.Prayers {
			dpt.Prayers[i].Play = false
		}
	}
	return svc.DailyPrayerTimings
}

// TurnOnAllAdhan sets adhan executions for all adhan timings of the month to be played
func (svc *Service) TurnOnAllAdhan() []DailyPrayerTimings {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	for _, dpt := range svc.DailyPrayerTimings {
		for i := range dpt.Prayers {
			dpt.Prayers[i].Play = true
		}
	}

	return svc.DailyPrayerTimings
}

// ToggleAdhan toggles a single adhan timings execution by matching its unix timestamp
func (svc *Service) ToggleAdhan(index uint8) (*Prayer, error) {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	for _, dpt := range svc.DailyPrayerTimings {
		for i := range dpt.Prayers {
			if dpt.Prayers[i].Index == index {
				dpt.Prayers[i].Play = !dpt.Prayers[i].Play
				return &dpt.Prayers[i], nil
			}
		}
	}

	return nil, fmt.Errorf("unable to find prayer with index; index=%d", index)
}

// getDateFromTimestamp retrieves time from a unix timestamp
func getDateFromTimestamp(timestamp string) (time.Time, error) {
	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Now(), err
	}
	return time.Unix(i, 0), nil
}

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
