package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type Sound struct {
	game                   *Game
	sounds                 []*wav.Stream
	audioPlayers           []*audio.Player
	sfxVolume, musicVolume float64
}

func NewSound(game *Game) *Sound {
	s := &Sound{}
	s.game = game
	if s.game.audioContext == nil {
		s.game.audioContext = audio.NewContext(48000)
	}
	s.initSound()
	return s
}

func (c *Sound) initSound() {
	c.sounds = []*wav.Stream{}
	c.audioPlayers = []*audio.Player{}

	c.addSFXByResourceWav(BoomWav)    //0
	c.addSFXByResourceWav(Exp1Wav)    //1
	c.addSFXByResourceWav(Exp2Wav)    //2
	c.addSFXByResourceWav(HissWav)    //3
	c.addSFXByResourceWav(Launch1Wav) //4
	c.addSFXByResourceWav(Launch2Wav) //5
	c.addSFXByResourceWav(Tone1Wav)   //6

}

func (c *Sound) initSoundF() {
	c.sounds = []*wav.Stream{}
	c.audioPlayers = []*audio.Player{}

	c.addSFXByFilename("boom.wav")    //0
	c.addSFXByFilename("exp1.wav")    //1
	c.addSFXByFilename("exp2.wav")    //2
	c.addSFXByFilename("hiss.wav")    //3
	c.addSFXByFilename("launch1.wav") //4
	c.addSFXByFilename("launch2.wav") //5
	c.addSFXByFilename("tone1.wav")   //6

}

func (c *Sound) addSFXByFilename(sfxID string) {
	stream := c.getStreamFromFilename(sfxID)
	c.sounds = append(c.sounds, stream)
	player := c.streamToAudioContexts(stream)
	c.audioPlayers = append(c.audioPlayers, player)
}

func (c *Sound) addSFXByResourceWav(resID []byte) {
	wavStream, err := wav.DecodeF32(bytes.NewReader(resID))
	if err != nil {
		log.Fatal(err)
	}
	//stream := c.getStreamFromFilename(sfxID)
	c.sounds = append(c.sounds, wavStream)
	player := c.streamToAudioContexts(wavStream)
	c.audioPlayers = append(c.audioPlayers, player)
}

func (c *Sound) PlaySFX(sfxID int) error {
	player := c.audioPlayers[sfxID]
	if err := player.Rewind(); err != nil {
		return err
	}
	player.SetVolume(c.sfxVolume)
	player.Play()
	return nil

}

func (c *Sound) SetSFXVolume(iVolume int) {
	fVolume := 0.1 * float64(iVolume)
	c.sfxVolume = Clamp(0.0, 1.0, fVolume)

}

func (c *Sound) SetMusicVolume(iVolume int) {
	fVolume := 0.1 * float64(iVolume)
	c.musicVolume = Clamp(0.0, 1.0, fVolume)

}

func (c *Sound) StopSFX(sfxID int) error {
	player := c.audioPlayers[sfxID]
	if err := player.Rewind(); err != nil {
		return err
	}
	player.Close()
	return nil

}

func (c *Sound) playAudioPlayer(player *audio.Player) error {
	if err := player.Rewind(); err != nil {
		return err
	}
	player.Play()
	return nil

}

func (c *Sound) streamToAudioContexts(stream *wav.Stream) *audio.Player {
	var player *audio.Player
	player, err := c.game.audioContext.NewPlayerF32(stream)
	if err != nil {
		log.Fatal(err)
	}
	return player
}

func (c *Sound) getStreamFromFilename(filename string) *wav.Stream {

	path := filepath.Join(c.game.soundSubdir, filename)
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	stat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return nil

	}
	stream, err := wav.DecodeF32(bytes.NewReader(bs))
	if err != nil {
		log.Fatal(err)
	}
	return stream
}
