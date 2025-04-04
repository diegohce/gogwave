package wav

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"

	"github.com/diegohce/gogwave"
	"github.com/diegohce/gogwave/ext/gogaudio"
)

type WavCodec struct{}

func newWavCodec(_ any) (gogaudio.AudioCodec, error) {
	return &WavCodec{}, nil
}

func (c *WavCodec) Decode(r io.ReadSeeker) ([]byte, error) {
	return DecodeFromWav(r)
}

func (c *WavCodec) Encode(w io.WriteSeeker, waveform []byte, SampleRateOut int, SampleFormatOut gogwave.GGWaveSampleFormatType) error {
	return EncodeToWav(w, waveform, SampleRateOut, SampleFormatOut)
}

func (c *WavCodec) Close() error {
	return nil
}

// DecodeFromWav reads from [r] and returns a [payload]
func DecodeFromWav(r io.ReadSeeker) ([]byte, error) {

	waveform, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	gg := gogwave.New()
	defer gg.Close()

	return gg.Decode(waveform)
}

// EncodeToWav writes [waveform] to [w] in wav format
func EncodeToWav(w io.WriteSeeker, waveform []byte, SampleRateOut int, SampleFormatOut gogwave.GGWaveSampleFormatType) error {
	var bitDepth int

	switch SampleFormatOut {
	case gogwave.GGWaveSampleFormatU8, gogwave.GGWaveSampleFormatI8:
		bitDepth = 8

	case gogwave.GGWaveSampleFormatU16, gogwave.GGWaveSampleFormatI16:
		bitDepth = 16

	default:
		bitDepth = 32
	}

	var value int32
	buf := bytes.NewBuffer(waveform)

	auBuf := audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: 1,
			SampleRate:  SampleRateOut,
		},
		SourceBitDepth: bitDepth,
	}

	for {
		value = 0
		err := binary.Read(buf, binary.LittleEndian, &value)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		if errors.Is(err, io.EOF) {
			break
		}
		auBuf.Data = append(auBuf.Data, int(value))
	}

	enc := wav.NewEncoder(w, SampleRateOut, bitDepth, 1, 1)
	defer enc.Close()

	err := enc.Write(&auBuf)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	gogaudio.Register("wav", newWavCodec)
}
