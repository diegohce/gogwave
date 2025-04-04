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
)

func main() {

	message := "Hello, World!"

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
	err = gg.EncodeToWav(f, waveformMessage)
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

	messageFromWav, err := gogwave.DecodeFromWav(f)
	if err != nil {
		panic(err)
	}

	fmt.Println("message:            ", message)
	fmt.Println("messageFromWaveform:", string(messageFromWaveform))
	fmt.Println("messageFromWav:     ", string(messageFromWav))
}
