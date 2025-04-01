package gogwave

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

// DecodeFromWav reads from [r] and returns a [payload]
func DecodeFromWav(r io.ReadSeeker) ([]byte, error) {

	waveform, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	gg := New()
	defer gg.Close()

	return gg.Decode(waveform)
}

// EncodeToWav writes [waveform] to [w] in wav format
func (gg *GGWave) EncodeToWav(w io.WriteSeeker, waveform []byte) error {
	var bitDepth int

	switch gg.params.sampleFormatOut {
	case GGWaveSampleFormatU8, GGWaveSampleFormatI8:
		bitDepth = 8

	case GGWaveSampleFormatU16, GGWaveSampleFormatI16:
		bitDepth = 16

	default:
		bitDepth = 32
	}

	var value int32
	buf := bytes.NewBuffer(waveform)

	auBuf := audio.IntBuffer{
		Format:         audio.FormatMono48000,
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

	enc := wav.NewEncoder(w,
		int(gg.params.sampleRateOut), bitDepth, 1, 1)
	defer enc.Close()

	err := enc.Write(&auBuf)
	if err != nil {
		return err
	}

	return nil
}
