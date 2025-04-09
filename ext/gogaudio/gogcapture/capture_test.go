package gogcapture

import (
	"os"
	"testing"
	"time"

	"github.com/diegohce/gogwave"
	"github.com/diegohce/gogwave/ext/gogaudio"
	_ "github.com/diegohce/gogwave/ext/gogaudio/wav"
)

func TestCapture(t *testing.T) {
	cap, _ := gogaudio.NewCapture("raw", nil)
	cap.Capture()

	time.Sleep(5 * time.Second)

	cap.Close()

	capturedRaw, _ := cap.CapturedBuffer()

	codec, _ := gogaudio.NewCodec("wav", nil)
	defer codec.Close()

	f, err := os.OpenFile("out.wav", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	codec.Encode(f, capturedRaw, 48000, gogwave.GGWaveSampleFormatF32)

}
