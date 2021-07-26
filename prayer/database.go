package prayer

import (
	"fmt"
	"sync"
	"time"
)

type PrayerDatabase interface {
	Timings() []DailyPrayerTimings
	SetTimings(prayerTimings []DailyPrayerTimings)
	SetPrayerTime(dailyPrayerIndex int, prayerIndex int, prayerTime Prayer)
	GetPrayerByTime(prayerTime time.Time) (Prayer, error)
}

type database struct {
	mutex         sync.RWMutex
	prayerTimings []DailyPrayerTimings
}

func NewPrayerDatabase() *database {
	return &database{prayerTimings: []DailyPrayerTimings{}}
}

func (db *database) SetTimings(prayerTimings []DailyPrayerTimings) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.prayerTimings = prayerTimings
}

func (db *database) Timings() []DailyPrayerTimings {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	return db.prayerTimings
}

func (db *database) SetPrayerTime(dailyPrayerIndex int, prayerIndex int, prayerTime Prayer) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.prayerTimings[dailyPrayerIndex].Prayers[prayerIndex] = prayerTime
}

func (db *database) GetPrayerByTime(prayerTime time.Time) (Prayer, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	for _, dpt := range db.prayerTimings {
		for _, prayer := range dpt.Prayers {
			if prayer.Time.Equal(prayerTime) {
				return prayer, nil
			}
		}
	}
	return Prayer{}, fmt.Errorf("unable to find prayer by time; t=%s", prayerTime)
}
