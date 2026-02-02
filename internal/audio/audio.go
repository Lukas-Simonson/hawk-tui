package audio

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var speakerInitialized = false

// InitSpeaker initializes the global speaker system once
// using the provided sample rate. This must be called before
// playing any audio to prevent startup lag or clipping.
func InitSpeaker(sampleRate beep.SampleRate) {
	if !speakerInitialized {
		speaker.Init(sampleRate, sampleRate.N(time.Second/10))
		speakerInitialized = true
		log.Println("Speaker initialized")
	}
}

// WarmUpSpeaker performs a short, silent playback to fully
// initialize the audio backend. This removes the initial
// delay that can cause short sounds to clip on first use.
func WarmUpSpeaker(sampleRate beep.SampleRate) {
	InitSpeaker(sampleRate)

	silence := beep.Callback(func() {})

	done := make(chan bool)
	speaker.Play(beep.Seq(
		silence,
		beep.Callback(func() {
			done <- true
		}),
	))

	<-done
}

// PlaySound plays an MP3 file located at the given relative path.
// It automatically initializes the speaker if needed and blocks
// until the sound has finished playing.
func PlaySound(path string) {

	// Resolve the absolute file path relative to the current working directory
	workingDirectory, _ := os.Getwd()
	fullPath := filepath.Join(workingDirectory, path)

	// Attempt to open the requested audio file
	audioFile, fileOpenError := os.Open(fullPath)

	if fileOpenError != nil {
		log.Println("failed to open audio file:", fileOpenError)
		return
	}
	defer audioFile.Close()

	// Decode the MP3 stream into an audio format usable by the speaker
	stream, format, decodeError := mp3.Decode(audioFile)
	if decodeError != nil {
		log.Println("failed to decode mp3:", decodeError)
		return
	}
	defer stream.Close()

	// Ensure the speaker is initialized for the correct sample rate
	InitSpeaker(format.SampleRate)

	playbackFinished := make(chan bool)

	speaker.Play(beep.Seq(
		stream,
		beep.Callback(func() {
			playbackFinished <- true
		}),
	))

	<-playbackFinished
}
