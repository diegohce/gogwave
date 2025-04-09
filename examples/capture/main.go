package main

import (
	"os"
	"time"

	"github.com/diegohce/gogwave"
	"github.com/diegohce/gogwave/ext/gogaudio"
	_ "github.com/diegohce/gogwave/ext/gogaudio/gogcapture"
	_ "github.com/diegohce/gogwave/ext/gogaudio/wav"
)

func main() {

	cap, _ := gogaudio.NewCapture("raw", nil)

	cap.Capture()

	time.Sleep(15 * time.Second)

	cap.Close()

	b, _ := cap.CapturedBuffer()

	codec, _ := gogaudio.NewCodec("wav", nil)
	defer codec.Close()

	f, _ := os.OpenFile("out.wav", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	defer f.Close()

	codec.Encode(f, b, 48000, gogwave.GGWaveSampleFormatF32)
}
