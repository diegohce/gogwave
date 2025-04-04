package wav

import (
	"os"
	"testing"

	"github.com/diegohce/gogwave"
)

func TestToWav(t *testing.T) {
	gg := gogwave.New()
	defer gg.Close()

	codec, _ := newWavCodec(nil)

	waveform, err := gg.Encode([]byte("hola"), gogwave.ProtocolAudibleNormal, 50)
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove("out.wav")

	f, _ := os.OpenFile("out.wav", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)

	err = codec.Encode(f, waveform, int(gg.Params.SampleRateOut), gg.Params.SampleFormatOut)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	f, _ = os.OpenFile("out.wav", os.O_RDONLY, 0644)

	b, err := codec.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	payload := string(b)
	if payload != "hola" {
		t.Fatalf("got %s want hola", payload)
	}

}
