package http

import (
	"fmt"

	"github.com/zees-dev/prayeralarm/aladhan"
)

type AdhanExecution struct {
	Status bool              `json:"status"`
	Adhan  aladhan.AdhanTime `json:"adhan"`
}

type AdhanExecutionMap map[int]AdhanExecution

func GenerateAdhanExecutionMap(executions []bool, timings []aladhan.AdhanTime) (AdhanExecutionMap, error) {
	if len(executions) != len(timings) {
		return nil, fmt.Errorf("adhan executions and timings length mismatch; len(executions):%d, len(timimngs):%d", len(executions), len(timings))
	}
	adhanExecutionMap := make(AdhanExecutionMap)
	for i, at := range timings {
		adhanExecutionMap[i] = AdhanExecution{
			Status: executions[i],
			Adhan:  at,
		}
	}

	return adhanExecutionMap, nil
}
