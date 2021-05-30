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

type Service struct {
	mutex           sync.RWMutex
	player          Player
	Timings         []aladhan.AdhanTime
	AdhanExecutions []bool
}

// NewService returns new adhan service that utilizes player to output adhan audio
func NewService(player Player) *Service {
	return &Service{player: player}
}

// SetAdhanTimings sets the adhan timings and the respective executions to true for the month
func (svc *Service) SetAdhanTimings(timings []aladhan.AdhanTime) {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	// set adhan executions to be true
	svc.AdhanExecutions = make([]bool, len(timings))
	for i := range timings {
		svc.AdhanExecutions[i] = true
	}

	svc.Timings = timings
}

// DisplayCalendar renders upcoming calendar in ASCII table
// https://github.com/olekukonko/tablewriter#example-6----identical-cells-merging
func (svc *Service) DisplayAdhanTimings(writer io.Writer) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Date", "Adhan", "Time"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	for _, adhan := range svc.Timings {
		year, month, day := adhan.Time.Date()
		dateStr := fmt.Sprintf("%s %d-%s-%d", adhan.Time.Weekday(), day, month, year)
		table.Append([]string{dateStr, string(adhan.Type), adhan.Time.Format("03:04:05 PM")})
	}
	table.Render()
}

// ExecuteAdhan plays adhan based on adhan timings if execution of the respective adhan is set to true
func (svc *Service) ExecuteAdhan() {
	// Play the adhan at the correct times - from current time
	for i, adhanTiming := range svc.Timings {
		timeTillNextAdhan := time.Until(adhanTiming.Time)
		log.Printf("Waiting %s for %s adhan...", timeTillNextAdhan, adhanTiming.Type)

		time.Sleep(timeTillNextAdhan)

		// Only play adhan if its set to execute
		if svc.AdhanExecutions[i] {
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
func (svc *Service) TurnOffAllAdhan() {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	for i := range svc.AdhanExecutions {
		svc.AdhanExecutions[i] = false
	}
}

// TurnOnAllAdhan sets adhan executions for all adhan timings of the month to be played
func (svc *Service) TurnOnAllAdhan() {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	for i := range svc.AdhanExecutions {
		svc.AdhanExecutions[i] = true
	}
}

// ToggleAdhan toggles a single adhan timings execution
func (svc *Service) ToggleAdhan(index uint8) {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	svc.AdhanExecutions[index] = !svc.AdhanExecutions[index]
}
