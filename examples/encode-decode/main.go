//go:build !test

// What it does:
//
// This example encodes a string into a ggwave waveform,
// saves the waveform in wav format (you can listen to it),
// and decodes the encoded waveform and wav file.
//
// How to run:
//
// go run main.go
//

package main

import (
	"fmt"
	"os"

	"github.com/diegohce/gogwave"
	"github.com/diegohce/gogwave/ext/gogaudio"
	_ "github.com/diegohce/gogwave/ext/gogaudio/wav"
)

func main() {

	message := "Hello, World!"

	wavcodec, err := gogaudio.NewCodec("wav", nil)
	if err != nil {
		panic(err)
	}
	defer wavcodec.Close()

	// Comment the next line for ggwave to write log output
	gogwave.SetLogFile(nil)

	gg := gogwave.New()
	defer gg.Close()

	// encode message as ggwave waveform
	waveformMessage, err := gg.Encode([]byte(message), gogwave.ProtocolAudibleNormal, 50)
	if err != nil {
		panic(err)
	}

	// save as .wav
	// you can listen to message.wav in your favorite music player
	// :-)
	f, err := os.OpenFile("message.wav", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	err = wavcodec.Encode(f, waveformMessage, int(gg.Params.SampleRateOut), gg.Params.SampleFormatOut)
	if err != nil {
		panic(err)
	}
	f.Close()

	// decode waveform
	messageFromWaveform, err := gg.Decode(waveformMessage)
	if err != nil {
		panic(err)
	}

	// decode .wav file
	f, err = os.Open("message.wav")
	if err != nil {
		panic(err)
	}

	messageFromWav, err := wavcodec.Decode(f)
	if err != nil {
		panic(err)
	}

	fmt.Println("message:            ", message)
	fmt.Println("messageFromWaveform:", string(messageFromWaveform))
	fmt.Println("messageFromWav:     ", string(messageFromWav))
}
