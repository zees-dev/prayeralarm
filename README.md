# Prayer Alarm

[![CI](https://github.com/zees-dev/prayeralarm/workflows/CI/badge.svg)](https://github.com/zees-dev/prayeralarm/actions?query=workflow%3ACI)
[![CD](https://github.com/zees-dev/prayeralarm/workflows/CD/badge.svg)](https://github.com/zees-dev/prayeralarm/actions?query=workflow%3ACD)
[![Release](https://github.com/zees-dev/prayeralarm/workflows/Release/badge.svg)](https://github.com/zees-dev/prayeralarm/releases)\
[![Go Report Card](https://goreportcard.com/badge/github.com/zees-dev/prayeralarm)](https://goreportcard.com/report/github.com/zees-dev/prayeralarm)

A cross-platform islamic prayer alarm.

The prayer alarm binary will run the adhan prayer call (audio) based on prayer timings retrieved from the [Adhan API](https://aladhan.com/prayer-times-api).

- Prayer timings are retrieved on a monthly basis (starting with current month by default).
- Prayer calls (adhan audio) are played at the respective prayer times.

The prayer alarm admin dashboard (Web UI) can be viewed at port `8080`.

## Prayer alarm configuration parameters

| Name      | Description                                                   | Value                 |
| --------- | ------------------------------------------------------------- | --------------------- |
| `city`    | City for which to retrieve prayer calendar                    | `"Auckland"`          |
| `country` | Country for which to retrieve prayer calendar                 | `"NewZealand"`        |
| `offsets` | Prayer call offsets to fine-tune prayer adhan timings (negative numbers are supported)          | `"0,0,0,0,0"` |
| `year`    | Year of prayer calendar                                       | `2021` (current year) |
| `month`   | Month of prayer calendar                                      | `6` (current month)   |
| `output`  | Output device to play adhan at prayer time; supported options are `stdout`, `native` and `omx`  | `omx`         |
| `port`    | Port to serve admin UI dashboard (web server)                 | `8080`                |

### Prayer time offsets

Offsetting prayer call times is also supported. Prayer call's can be offset by a specified number of minutes by providing an optional **offsets** flag when running the binary.  
For example, to respectively offset the _Maghrib_ and _Isha_ prayer calls to run 5 mins later and 3 mins earlier, the binary can be run with the following flag: **-offsets "0,0,0,5,-3"**  
By default, offsets for all prayer times are set to **0**; i.e. **0,0,0,0,0**.

## Development

### Pre-requisites

- Golang (developed on v1.15)
- Pre-requisites for sound/output dependencies
  - `output=native`
    - [Oto](https://github.com/hajimehoshi/oto) - based on OS
    - `Oto` requires `CGO_ENABLED=1`
  - `output=omx`
    - [Omxplayer](https://github.com/huceke/omxplayer) - A CLI application that can play audio files

### Steps

Use the [Makefile](./Makefile) to test, build and run the project; alternatively manual instructions are defined below.

#### Install dependencies

```sh
go mod download
```

#### Build binary

```sh
go build
```

Note: If using `native` `output`, then build with: `CGO_ENABLED=1 go build`

#### Run binary - examples

**In foreground:**

```sh
./prayeralarm
```

- Run with overrides - optional city, country and offset flags

```sh
./prayeralarm -city Auckland -country NewZealand -offsets "5,0,-5,-10,0"
```

**In background (as service) - with log file:**
  
```sh
nohup ./prayeralarm > adhan.log &
```

- Kill background process
  
```sh
kill $(ps -ef | grep prayeralarm| cut -f4 -d" " | head -1)
```

## Dependencies

- [Oto](https://github.com/hajimehoshi/oto) for sound playback (playing Adhan audio)
- [go-mp3](https://github.com/hajimehoshi/go-mp3) for cross-platform MP3 playback (playing Adhan audio)
- pulseaudio server on host (if running prayeralarm in container)
  - run pulseaudio server: `pulseaudio --load=module-native-protocol-tcp --exit-idle-time=-1 --daemon`
  - stop pulseaudio server: `pulseaudio --kill`

## Licence

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

---

## TODO

- [x] Integration testing
- [x] Tagged releases - using semantic versioning
- [x] Cross-platform release binaries
  - [ ] Raspberry Pi release binary
  - [x] Docker based setup - [example](https://gitlab.com/dev.786zshan/golang-project-bootstrapper)
- [x] CI pipeline
- [x] Multiplatform release binaries
- [x] CD pipeline
  - [ ] One workflow should release to Raspberry PI  

## Roadmap

- [x] A web-ui to view prayer calendar - for current day and month
- [x] Ability to toggle on-off prayer call
