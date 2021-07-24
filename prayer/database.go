package prayer

import "sync"

type PrayerDatabase interface {
	Timings() []DailyPrayerTimings
	SetTimings(prayerTimings []DailyPrayerTimings)
	SetPrayerTime(dailyPrayerIndex int, prayerIndex int, prayerTime Prayer)
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
