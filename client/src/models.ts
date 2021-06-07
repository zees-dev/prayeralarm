
type Adhan = "Fajr" | "Dhuhr" | "Asr" | "Maghrib" | "Isha"

export interface Timing {
    date: string
    prayers: PrayerCall[]
}
export interface PrayerCall {
    play: boolean
    time: string
    type: Adhan
    index: number
}
