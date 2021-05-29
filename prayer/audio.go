package prayer

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"github.com/zees-dev/prayeralarm/aladhan"
)

type Player interface {
	Play(adhan aladhan.Adhan) error
}

// Default stdout for testing purposes
// TODO replicate similar functionality to http defaultservemux
type stdOut struct{}

func NewStdOutPlayer() stdOut {
	return stdOut{}
}

func (so stdOut) Play(adhan aladhan.Adhan) error {
	fmt.Println(adhan)
	return nil
}

// omxplayer binary (must be present in OS)
type omxPlayer struct{}

func NewOmxPlayer() omxPlayer {
	return omxPlayer{}
}

// Play plays an audio via omxplayer cli tool
func (op omxPlayer) Play(adhan aladhan.Adhan) error {
	var filename string
	switch adhan {
	case "Fajr":
		filename = "mp3/adhan-fajr.mp3"
	default:
		filename = "mp3/adhan-turkish.mp3"
	}
	// fajrAdhan = 'omxplayer -o local --vol 1000 mp3/adhan-fajr.mp3 > /dev/null 2>&1'
	// otherAdhan = 'omxplayer -o local --vol 1000 mp3/adhan-turkish.mp3 > /dev/null 2>&1'
	commandStr := strings.Split(fmt.Sprintf("omxplayer -o local --vol 1000 %s > /dev/null 2>&1", filename), " ")

	log.Println(fmt.Sprintf("executing command: %v", commandStr))

	_, err := exec.Command(commandStr[0], commandStr[1:]...).Output()
	if err != nil {
		return err
	}

	return nil
}

// mp3Player is primarily used to implement interface output mp3 to audio output device
type mp3Player struct{}

func NewMp3Player() mp3Player {
	return mp3Player{}
}

// Play reads the mp3 adhan file ) and outputs mp3 to audio device using `oto` and `go-mp3`
func (mp mp3Player) Play(adhan aladhan.Adhan) error {
	var filename string
	switch adhan {
	case "Fajr":
		filename = "mp3/adhan-fajr.mp3"
	default:
		filename = "mp3/adhan-turkish.mp3"
	}

	adhanF, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer adhanF.Close()

	decoder, err := mp3.NewDecoder(adhanF)
	if err != nil {
		return err
	}

	c, err := oto.NewContext(decoder.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer c.Close()

	player := c.NewPlayer()
	defer player.Close()

	fmt.Printf("playing bytes: %d[bytes]\n", decoder.Length())
	if _, err := io.Copy(player, decoder); err != nil {
		return err
	}
	return nil
}
