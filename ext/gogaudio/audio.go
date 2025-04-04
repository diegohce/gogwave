package gogaudio

import (
	"errors"
	"io"

	"github.com/diegohce/gogwave"
)

type AudioCodec interface {
	Decode(r io.ReadSeeker) ([]byte, error)
	Encode(w io.WriteSeeker, waveform []byte, SampleRateOut int, SampleFormatOut gogwave.GGWaveSampleFormatType) error
	Close() error
}

type NewAudioCodecFunc func(config any) (AudioCodec, error)

var (
	codecs          = map[string]NewAudioCodecFunc{}
	ErrInvalidCodec = errors.New("invalid codec")
)

func Register(name string, fn NewAudioCodecFunc) {
	codecs[name] = fn
}

func NewCodec(name string, config any) (AudioCodec, error) {
	fn, exists := codecs[name]
	if !exists {
		return nil, ErrInvalidCodec
	}
	return fn(config)
}
