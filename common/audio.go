package common

import (
	"bytes"
	"io"

	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/assets"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const audioSampleRate = 44100

var AudioManager *audioManager

type Sound int

const (
	SoundClick Sound = iota
	SoundDrop
	SoundSnap
	SoundComplete
)

type audioManager struct {
	context     *audio.Context
	sfxVolume   float64
	musicVolume float64
	musicPlayer *audio.Player
	sfxPlayers  []*audio.Player
	soundData   map[Sound][]byte
}

func NewAudioManager() *audioManager {
	s := GetSettings()
	am := &audioManager{
		context:   audio.NewContext(audioSampleRate),
		sfxVolume: s.SFXVolume,
		soundData: make(map[Sound][]byte),
	}
	am.SetMusicVolume(s.MusicVolume)
	am.loadSounds()
	return am
}

func (am *audioManager) loadSounds() {
	type entry struct {
		sound Sound
		name  string
	}
	entries := []entry{
		{SoundClick, "click"},
		{SoundDrop, "drop"},
		{SoundSnap, "snap"},
		{SoundComplete, "complete"},
	}
	for _, e := range entries {
		data := loadAudioFile(e.name)
		if len(data) > 0 {
			am.soundData[e.sound] = data
		}
	}
}

func loadAudioFile(name string) []byte {
	data, err := assets.Assets.ReadFile("audio/" + name + ".ogg")
	if err == nil && len(data) > 0 {
		return data
	}
	return nil
}

func isOGG(data []byte) bool {
	return len(data) > 4 && data[0] == 'O' && data[1] == 'g' && data[2] == 'g' && data[3] == 'S'
}

func decodeAudio(data []byte) (io.ReadSeeker, error) {
	if isOGG(data) {
		return vorbis.DecodeF32(bytes.NewReader(data))
	}
	return wav.DecodeF32(bytes.NewReader(data))
}

func (am *audioManager) SetSFXVolume(v float64) {
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	am.sfxVolume = v
}

func (am *audioManager) SetMusicVolume(v float64) {
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	am.musicVolume = v / 15
	if am.musicPlayer != nil {
		am.musicPlayer.SetVolume(am.musicVolume)
	}
}

func (am *audioManager) SFXVolume() float64   { return am.sfxVolume }
func (am *audioManager) MusicVolume() float64 { return am.musicVolume }

func (am *audioManager) Play(sound Sound) {
	data, ok := am.soundData[sound]
	if !ok || len(data) == 0 || am.sfxVolume <= 0 {
		return
	}

	stream, err := decodeAudio(data)
	if err != nil {
		return
	}

	player, err := am.context.NewPlayerF32(stream)
	if err != nil {
		return
	}
	player.SetVolume(am.sfxVolume)
	player.Play()

	var active []*audio.Player
	for _, p := range am.sfxPlayers {
		if p.IsPlaying() {
			active = append(active, p)
		}
	}
	am.sfxPlayers = append(active, player)
}

func (am *audioManager) PlayClick()    { am.Play(SoundClick) }
func (am *audioManager) PlayDrop()     { am.Play(SoundDrop) }
func (am *audioManager) PlaySnap()     { am.Play(SoundSnap) }
func (am *audioManager) PlayComplete() { am.Play(SoundComplete) }

func (am *audioManager) StartMusic(data []byte) {
	if am.musicPlayer != nil {
		am.musicPlayer.Close()
		am.musicPlayer = nil
	}

	if len(data) == 0 {
		return
	}

	r := bytes.NewReader(data)
	var stream io.ReadSeeker
	var length int64
	if isOGG(data) {
		s, err := vorbis.DecodeF32(r)
		if err != nil {
			return
		}
		stream = s
		length = s.Length()
	} else {
		s, err := wav.DecodeF32(r)
		if err != nil {
			return
		}
		stream = s
		length = s.Length()
	}

	loop := audio.NewInfiniteLoopF32(stream, length)
	player, err := am.context.NewPlayerF32(loop)
	if err != nil {
		return
	}
	player.SetVolume(am.musicVolume)
	player.Play()
	am.musicPlayer = player
}

func (am *audioManager) StopMusic() {
	if am.musicPlayer != nil {
		am.musicPlayer.Close()
		am.musicPlayer = nil
	}
}

func (am *audioManager) StartMusicFromFile(name string) {
	data := loadAudioFile(name)
	if len(data) > 0 {
		am.StartMusic(data)
	}
}
