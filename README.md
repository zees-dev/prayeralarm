# Prayer Alarm

A cross-platform single binary islamic prayer alarm.

The prayer alarm binary will run the adhan prayer call (audio) based on prayer timings retrieved from the [Adhan API](https://aladhan.com/prayer-times-api) for a specified location.

By default, the prayer call timings are set for:

- city = `auckland`
- country = `new zealand`

## How it works

Prayer timings are retrieved on a monthly basis (starting with current month).
The prayer calls (audio) is played at the exact prayer times - for specified location.

### Config overrides

#### Prayer time offsets

Offsetting prayer call times is also supported. Prayer call's can be offset by a specified number of minutes by providing an optional **offsets** flag when running the binary.  
For example, to respectively offset the _Maghrib_ and _Isha_ prayer calls to run 5 mins later and 3 mins earlier, the binary can be run with the following flag: **-offsets "0 0 0 5 -3"**  
By default, offsets for all prayer times are set to **0**.

## Development

### Pre-requisites

- Golang (developed on v1.13)
- Pre-requsites for sound player dependency [Oto](https://github.com/hajimehoshi/oto - based on OS

### Steps

#### Build binary

  ```sh
  go build
  ```

#### Run binary

- In foreground

  ```sh
  ./prayeralarm-go
  ```

- Run with overrides - optional city, country and offset flags

  ```sh
  ./prayeralarm-go -city auckland -country "new zealand" -offsets "5 0 -5 -10 0"
  ```

- In background with log file
  
  ```sh
  nohup ./prayeralarm-go > adhan.log &
  ```

  - Kill background process
  
    ```sh
    kill $(ps -ef | grep prayeralarm-go| cut -f4 -d" " | head -1)
    ```

## Dependencies

- [Oto](https://github.com/hajimehoshi/oto) for cross-platform MP3 playback (playing Adhan audio)
- [packr2](https://github.com/gobuffalo/packr/tree/master/v2) for converting adhan audio files (mp3) into bytes which have been committed to source control as a single [file](./internal/mp3/mp3-files.go); this allows us to distribute the generated binary without depending on external adhan audio files on the file system

## Licence

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

---

## R&D

- [ ] compare omx player sound with golang based sound
- [ ] view adhaan API
- [ ] create relevant JSON struct(s) - can use JSON converter plugin
- [ ] use jsonschema to verify API
- [ ] gracefully handle timezones
- [ ] use config file (.env) and flag options for custom timing overrides

---

## TODO

- [ ] Integration testing
- [ ] Tagged releases - using semantiv versioning
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
