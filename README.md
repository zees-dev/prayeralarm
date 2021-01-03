# Prayer Alarm

A cross-platform single binary islamic prayer alarm.

The prayer alarm binary will run the adhan prayer call (audio) based on prayer timings retrieved from the [Adhan API](https://aladhan.com/prayer-times-api) for a specified location.

By default, the prayer call timings are set for:

- city = `Auckland`
- country = `NewZealand`

## How it works

Prayer timings are retrieved on a monthly basis (starting with current month by defaultt).
The prayer call (audio) is played at the respective prayer time - for specified location.

### Config overrides

#### Prayer time offsets

Offsetting prayer call times is also supported. Prayer call's can be offset by a specified number of minutes by providing an optional **offsets** flag when running the binary.  
For example, to respectively offset the _Maghrib_ and _Isha_ prayer calls to run 5 mins later and 3 mins earlier, the binary can be run with the following flag: **-offsets "0,0,0,5,-3"**  
By default, offsets for all prayer times are set to **0**; i.e. **0,0,0,0,0**.

## Development

### Pre-requisites

- Golang (developed on v1.13)
- Pre-requsites for sound player dependency [Oto](https://github.com/hajimehoshi/oto) - based on OS

### Steps

#### Install dependencies

  ```sh
  go mod download
  ```

#### Build binary

  ```sh
  CGO_CPPFLAGS='-Wno-error -Wno-nullability-completeness -Wno-expansion-to-defined -Wbuiltin-requires-header' CGO_ENABLED=1 go build
  ```

#### Run binary

- In foreground

  ```sh
  ./prayeralarm
  ```

- Run with overrides - optional city, country and offset flags

  ```sh
  ./prayeralarm -city auckland -country NewZealand -offsets "5,0,-5,-10,0"
  ```

- In background with log file
  
  ```sh
  nohup ./prayeralarm > adhan.log &
  ```

  - Kill background process
  
    ```sh
    kill $(ps -ef | grep prayeralarm| cut -f4 -d" " | head -1)
    ```

## Dependencies

- [Oto](https://github.com/hajimehoshi/oto) for cross-platform MP3 playback (playing Adhan audio)
- [go-mp3](https://github.com/hajimehoshi/go-mp3) for cross-platform MP3 playback (playing Adhan audio)

## Licence

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

---

## TODO

- [ ] Integration testing
- [ ] Tagged releases - using semantic versioning
- [ ] Cross-platform release binaries
  - [ ] Raspberry Pi release binary
  - [ ] Docker based setup - [example](https://gitlab.com/dev.786zshan/golang-project-bootstrapper)
- [ ] CI pipeline
- [ ] CD pipeline
  - [ ] One workflow should release to my Raspberry PI  

## Roadmap

- [ ] A web-ui to view prayer calendar - for current day and month
- [ ] Ability to toggle on-off prayer call

## File conversions

`.mp3` files converted to `.wav` files using [ffmpeg]

<!-- Imsak,Fajr,Sunrise,Dhuhr,Asr,Maghrib,Sunset,Isha,Midnight -->
curl 'http://api.aladhan.com/v1/calendarByCity?city=Auckland&country=NewZealand&method=3&month=12&year=2020&tune=0,0,0,0,0,0,0,0'
