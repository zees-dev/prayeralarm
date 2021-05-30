package prayer

import (
	"fmt"
	"io"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/zees-dev/prayeralarm/aladhan"
)

type Prayer struct {
	Play bool          `json:"play"`
	Type aladhan.Adhan `json:"type"`
	Time time.Time     `json:"time"`
}

type Service struct {
	mutex   sync.RWMutex
	player  Player
	Prayers []Prayer
}

// NewService returns new adhan service that utilizes player to output adhan audio
func NewService(player Player) *Service {
	return &Service{player: player}
}

// GeneratePrayers extracts the monthly adhan timings from the calendar api response
func (svc *Service) GeneratePrayers(monthCalendar aladhan.MonthlyAdhanCalenderResponse) {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	svc.Prayers = make([]Prayer, 0)

	// Get all adhan timings after current time for remaining days of the month
	currentTime := time.Now()
	for _, timings := range monthCalendar.Data {
		for adhan, timeStr := range timings.Timings {
			fullTimeStr := fmt.Sprintf("%s %s", timings.Date.Readable, timeStr)
			adhanTime := getTime(fullTimeStr, timings.Meta.Timezone)
			if adhanTime.After(currentTime) {
				svc.Prayers = append(svc.Prayers, Prayer{Play: true, Type: adhan, Time: adhanTime})
			}
		}
	}

	// Sort upcoming adhans by time
	sort.Slice(svc.Prayers[:], func(i, j int) bool {
		return svc.Prayers[i].Time.Before(svc.Prayers[j].Time)
	})
}

// DisplayPrayerTimings renders upcoming calendar in ASCII table
// https://github.com/olekukonko/tablewriter#example-6----identical-cells-merging
func (svc *Service) DisplayPrayerTimings(writer io.Writer) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Date", "Adhan", "Time", "Play"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	for _, p := range svc.Prayers {
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
	table.Render()
}

// ExecuteAdhan plays adhan based on adhan timings if execution of the respective adhan is set to true
func (svc *Service) ExecutePrayers() {
	// Play the adhan at the correct times - from current time
	for _, p := range svc.Prayers {
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

// TurnOffAllAdhan sets adhan executions for all adhan timings of the month to be muted
func (svc *Service) TurnOffAllAdhan() {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	for i := range svc.Prayers {
		svc.Prayers[i].Play = false
	}
}

// TurnOnAllAdhan sets adhan executions for all adhan timings of the month to be played
func (svc *Service) TurnOnAllAdhan() {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	for i := range svc.Prayers {
		svc.Prayers[i].Play = true
	}
}

// ToggleAdhan toggles a single adhan timings execution
func (svc *Service) ToggleAdhan(index uint8) {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	svc.Prayers[index].Play = !svc.Prayers[index].Play
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
