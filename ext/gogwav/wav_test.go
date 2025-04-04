package gogwav_test

import (
	"os"
	"testing"

	"github.com/diegohce/gogwave"
	"github.com/diegohce/gogwave/ext/gogwav"
)

func TestToWav(t *testing.T) {
	gg := gogwave.New()
	defer gg.Close()

	waveform, err := gg.Encode([]byte("hola"), gogwave.ProtocolAudibleNormal, 50)
	if err != nil {
		t.Fatal(err)
	}

	f, _ := os.OpenFile("out.wav", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)

	err = gogwav.EncodeToWav(f, waveform, int(gg.Params.SampleRateOut), gg.Params.SampleFormatOut)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	f, _ = os.OpenFile("out.wav", os.O_RDONLY, 0644)

	b, err := gogwav.DecodeFromWav(f)
	if err != nil {
		t.Fatal(err)
	}

	payload := string(b)
	if payload != "hola" {
		t.Fatalf("got %s want hola", payload)
	}

}
