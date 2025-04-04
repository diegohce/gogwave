package gogaudio

import (
	"io"
	"testing"

	"github.com/diegohce/gogwave"
)

func TestInvalidCodec(t *testing.T) {
	_, err := NewCodec("invalid", nil)
	if err == nil {
		t.Fatalf("got nil error want error")
	}
}

type noopCodec struct{}

func (c *noopCodec) Decode(r io.ReadSeeker) ([]byte, error) {
	return nil, nil
}

func (c *noopCodec) Encode(w io.WriteSeeker, waveform []byte, SampleRateOut int, SampleFormatOut gogwave.GGWaveSampleFormatType) error {
	return nil
}

func (c *noopCodec) Close() error { return nil }

func TestRegisterCodec(t *testing.T) {

	newNoopCodec := func(_ any) (AudioCodec, error) {
		return &noopCodec{}, nil
	}
	Register("noop", newNoopCodec)
}

func TestValidCodec(t *testing.T) {

	newNoopCodec := func(_ any) (AudioCodec, error) {
		return &noopCodec{}, nil
	}
	Register("noop", newNoopCodec)

	c, err := NewCodec("noop", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
}
