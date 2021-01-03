package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

// displayCalendar renders upcoming calendar in ASCII table
// https://github.com/olekukonko/tablewriter#example-6----identical-cells-merging
func displayCalendar(adhanSlice []AdhanTime) {
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
