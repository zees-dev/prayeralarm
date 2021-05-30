package prayer

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/zees-dev/prayeralarm/aladhan"
)

type service struct {
	mutex           sync.RWMutex
	timings         []aladhan.AdhanTime
	adhanExecutions []bool
	player          Player
}

// NewService returns new adhan service that utilizes player to output adhan audio
func NewService(player Player) *service {
	return &service{player: player}
}

// SetAdhanTimings sets the adhan timings and the respective executions to true for the month
func (svc *service) SetAdhanTimings(timings []aladhan.AdhanTime) {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	// set adhan executions to be true
	svc.adhanExecutions = make([]bool, len(timings))
	for i := range timings {
		svc.adhanExecutions[i] = true
	}

	svc.timings = timings
}

// DisplayCalendar renders upcoming calendar in ASCII table
// https://github.com/olekukonko/tablewriter#example-6----identical-cells-merging
func (svc *service) DisplayAdhanTimings(writer io.Writer) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Date", "Adhan", "Time"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	for _, adhan := range svc.timings {
		year, month, day := adhan.Time.Date()
		dateStr := fmt.Sprintf("%s %d-%s-%d", adhan.Time.Weekday(), day, month, year)
		table.Append([]string{dateStr, string(adhan.Type), adhan.Time.Format("03:04:05 PM")})
	}
	table.Render()
}

// ExecuteAdhan plays adhan based on adhan timings if execution of the respective adhan is set to true
func (svc *service) ExecuteAdhan() {
	// Play the adhan at the correct times - from current time
	for i, adhanTiming := range svc.timings {
		timeTillNextAdhan := time.Until(adhanTiming.Time)
		log.Printf("Waiting %s for %s adhan...", timeTillNextAdhan, adhanTiming.Type)

		time.Sleep(timeTillNextAdhan)

		// Only play adhan if its set to execute
		if svc.adhanExecutions[i] {
			log.Printf("Playing %s adhan at %s...", adhanTiming.Type, adhanTiming.Time)
			if err := svc.player.Play(adhanTiming.Type); err != nil {
				log.Fatalln(err)
			}
		} else {
			log.Printf("Skipping %s adhan at %s since execution is set to false", adhanTiming.Type, adhanTiming.Time)
		}
	}
}

// TurnOffAllAdhan sets adhan executions for all adhan timings of the month to be muted
func (svc *service) TurnOffAllAdhan() {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	for i := range svc.adhanExecutions {
		svc.adhanExecutions[i] = false
	}
}

// TurnOnAllAdhan sets adhan executions for all adhan timings of the month to be played
func (svc *service) TurnOnAllAdhan() {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	for i := range svc.adhanExecutions {
		svc.adhanExecutions[i] = true
	}
}

// ToggleAdhan toggles a single adhan timings execution
func (svc *service) ToggleAdhan(index uint8) {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	svc.adhanExecutions[index] = !svc.adhanExecutions[index]
}
