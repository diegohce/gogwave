package gogaudio

import (
	"errors"
	"io"

	"github.com/diegohce/gogwave"
)

type AudioCodec interface {
	Decode(r io.Reader) ([]byte, error)
	Encode(w io.WriteSeeker, waveform []byte, SampleRateOut int, SampleFormatOut gogwave.GGWaveSampleFormatType) error
	Close() error
}

type AudioCapture interface {
	Capture() error
	Close() error
	CapturedBuffer() ([]byte, error)
}

type NewAudioCodecFunc func(config any) (AudioCodec, error)
type NewAudioCaptureFunc func(config any) (AudioCapture, error)

var (
	codecs          = map[string]NewAudioCodecFunc{}
	ErrInvalidCodec = errors.New("invalid codec")

	capturers         = map[string]NewAudioCaptureFunc{}
	ErrInvalidCapture = errors.New("invalid capture")
)

func RegisterCodec(name string, fn NewAudioCodecFunc) {
	codecs[name] = fn
}

func RegisterCapture(name string, fn NewAudioCaptureFunc) {
	capturers[name] = fn
}

func NewCodec(name string, config any) (AudioCodec, error) {
	fn, exists := codecs[name]
	if !exists {
		return nil, ErrInvalidCodec
	}
	return fn(config)
}

func NewCapture(name string, config any) (AudioCapture, error) {
	fn, exists := capturers[name]
	if !exists {
		return nil, ErrInvalidCapture
	}
	return fn(config)
}
