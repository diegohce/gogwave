package gogwave

/*
#include <stdio.h>
#include <stdlib.h>
#include <ggwave/ggwave.h>
*/
import "C"

import (
	"errors"
	"os"
	"unsafe"
)

// GGWaveParameters configuration options for ggwave and gogwave.
type GGWaveParameters struct {
	// PayloadLength payload length
	PayloadLength int

	// SampleRateInp capture sample rate
	SampleRateInp float32

	// SampleRateOut playback sample rate
	SampleRateOut float32

	// SampleRate the operating sample rate
	SampleRate float32

	// SamplesPerFrame number of samples per audio frame
	SamplesPerFrame int

	// SoundMarkerThreshold sound marker detection threshold
	SoundMarkerThreshold float32

	// SampleFormatInp format of the captured audio samples
	SampleFormatInp GGWaveSampleFormatType

	// SampleFormatOut format of the playback audio samples
	SampleFormatOut GGWaveSampleFormatType

	// OperatingMode operating mode
	OperatingMode int

	// DecodeBufferSize max decoding buffer size
	DecodeBufferSize int
}

func NewGGwaveParameters() GGWaveParameters {
	c_gg_params := C.ggwave_getDefaultParameters()

	params := GGWaveParameters{
		PayloadLength:        int(c_gg_params.payloadLength),
		SampleRateInp:        float32(c_gg_params.sampleRateInp),
		SampleRateOut:        float32(c_gg_params.sampleRateOut),
		SampleRate:           float32(c_gg_params.sampleRate),
		SamplesPerFrame:      int(c_gg_params.samplesPerFrame),
		SoundMarkerThreshold: float32(c_gg_params.soundMarkerThreshold),
		SampleFormatInp:      GGWaveSampleFormatType(c_gg_params.sampleFormatInp),
		SampleFormatOut:      GGWaveSampleFormatType(c_gg_params.sampleFormatOut),
		OperatingMode:        int(c_gg_params.operatingMode),
		DecodeBufferSize:     4096,
	}

	return params
}

type GGWave struct {
	i                C.ggwave_Instance
	Params           GGWaveParameters
	decodeBufferSize int
}

// New Creates a new GGWave instance using default paramaters values.
//
// Calling:
//
//	gg := gogwave.New()
//
// is the same as:
//
//	 p := gogwave.NewGGwaveParameters()
//	 gg := gogwave.NewWithParams(p)
//
//		Note: Make sure to deallocate the instance at the end by calling Close()
func New() *GGWave {
	c_gg_params := C.ggwave_getDefaultParameters()

	inst := C.ggwave_init(c_gg_params)

	params := GGWaveParameters{
		PayloadLength:        int(c_gg_params.payloadLength),
		SampleRateInp:        float32(c_gg_params.sampleRateInp),
		SampleRateOut:        float32(c_gg_params.sampleRateOut),
		SampleRate:           float32(c_gg_params.sampleRate),
		SamplesPerFrame:      int(c_gg_params.samplesPerFrame),
		SoundMarkerThreshold: float32(c_gg_params.soundMarkerThreshold),
		SampleFormatInp:      GGWaveSampleFormatType(c_gg_params.sampleFormatInp),
		SampleFormatOut:      GGWaveSampleFormatType(c_gg_params.sampleFormatOut),
		OperatingMode:        int(c_gg_params.operatingMode),
	}

	return &GGWave{i: inst, Params: params, decodeBufferSize: 4096}
}

// NewWithParams Create a new GGWave instance with specified parameters.
// Calling:
//
//	gg := gogwave.New()
//
// is the same as:
//
//	 p := gogwave.NewGGwaveParameters()
//	 gg := gogwave.NewWithParams(p)
//
//		Note: Make sure to deallocate the instance at the end by calling Close()
func NewWhithParams(params GGWaveParameters) *GGWave {

	c_gg_params := C.ggwave_Parameters{
		payloadLength:        C.int(params.PayloadLength),
		sampleRateInp:        C.float(params.SampleRateInp),
		sampleRateOut:        C.float(params.SampleRateOut),
		sampleRate:           C.float(params.SampleRate),
		samplesPerFrame:      C.int(params.SamplesPerFrame),
		soundMarkerThreshold: C.float(params.SoundMarkerThreshold),
		sampleFormatInp:      C.ggwave_SampleFormat(params.SampleFormatInp),
		sampleFormatOut:      C.ggwave_SampleFormat(params.SampleFormatOut),
		operatingMode:        C.int(params.OperatingMode),
	}

	inst := C.ggwave_init(c_gg_params)

	return &GGWave{i: inst, Params: params, decodeBufferSize: params.DecodeBufferSize}
}

// Close Free GGWave resources
func (gg *GGWave) Close() error {
	C.ggwave_free(gg.i)
	return nil
}

// DecodeBufferSize returns the maximum decode buffer size
func (gg *GGWave) DecodeBufferSize() int {
	return gg.decodeBufferSize
}

// SetDecodeBufferSize sets the maximum decode buffer size.
//
// Default size: 4096 bytes.
func (gg *GGWave) SetDecodeBufferSize(size int) {
	gg.decodeBufferSize = size
}

// Encode data into audio waveform
//
//	payload  - the data to encode
//	protocol - the GGWaveProtocolType to use for encoding
//	volume   - the volume of the generated waveform [0, 100]
//	           usually 25 is OK and you should not go over 50
//
// returns:
//
//	waveform - the generated audio waveform.
func (gg *GGWave) Encode(payload []byte, protocol GGWaveProtocolType, volume int) ([]byte, error) {

	payloadPtr := unsafe.SliceData(payload)

	waveSize := C.ggwave_encode(gg.i, unsafe.Pointer(payloadPtr), C.int(len(payload)), C.ggwave_ProtocolId(ProtocolAudibleNormal), C.int(25), nil, C.int(1))
	if waveSize < C.int(0) {
		return nil, errors.New("encoding error")
	}

	out := make([]byte, int(waveSize))

	outPtr := unsafe.SliceData(out)

	waveSize = C.ggwave_encode(gg.i, unsafe.Pointer(payloadPtr), C.int(len(payload)), C.ggwave_ProtocolId(ProtocolAudibleNormal), C.int(25), unsafe.Pointer(outPtr), C.int(0))
	if waveSize < C.int(0) {
		return nil, errors.New("encoding error")
	}

	return out, nil
}

// Decode an audio waveform into data
//
//	waveform - the audio waveform
//
// returns:
//
//	payloadBuffer - stores the decoded data on success
func (gg *GGWave) Decode(waveform []byte) ([]byte, error) {

	waveformPtr := unsafe.SliceData(waveform)

	payload := make([]byte, gg.decodeBufferSize)
	payloadPtr := unsafe.SliceData(payload)

	bcount := C.ggwave_decode(gg.i, unsafe.Pointer(waveformPtr), C.int(len(waveform)), unsafe.Pointer(payloadPtr))
	if int(bcount) <= 0 {
		return nil, errors.New("error decoding waveform")
	}

	return payload[:int(bcount)], nil
}

// RxDurationFrames returns recvDuration_frames value for a rx protocol
func (gg *GGWave) RxDurationFrames() int {
	df := C.ggwave_rxDurationFrames(gg.i)
	return int(df)
}

// SetLogFile Change file stream for internal ggwave logging. NULL - disable logging
//
//	Note: not thread-safe. Do not call while any GGWave instances are running
func SetLogFile(f *os.File) {
	if f == nil {
		C.ggwave_setLogFile(nil)
		return
	}

	fd := f.Fd()
	c_mode := C.CString("a")
	defer C.free(unsafe.Pointer(c_mode))

	fp := C.fdopen(C.int(fd), c_mode)
	C.ggwave_setLogFile(unsafe.Pointer(fp))
}

// RxToggleProtocol Toggle Rx protocols on and off
//
//	protocol - GGWaveProtocolType of the Tx protocol to modify
//	state    - false: disable || true: enable
//
// If an Rx protocol is enabled, newly constructued GGWave instances will attempt to decode
// received data using this protocol. By default, all protocols are enabled.
// Use this function to restrict the number of Rx protocols used in the decoding
// process. This helps to reduce the number of false positives and improves the transmission
// accuracy, especially when the Tx/Rx protocol is known in advance.
//
//	Note: this function does not affect the decoding process of instances that have
//	already been created.
func RxToggleProtocol(protocol GGWaveProtocolType, state bool) {
	s := C.int(0)
	if state {
		s = C.int(1)
	}
	C.ggwave_rxToggleProtocol(C.ggwave_ProtocolId(protocol), s)
}

// TxToggleProtocol Toggle Tx protocols on and off
//
//	protocol - GGWaveProtocolType of the Tx protocol to modify
//	state    - false: disable || true: enable
//
// If an Tx protocol is enabled, newly constructued GGWave instances will be able to transmit
// data using this protocol. By default, all protocols are enabled.
// Use this function to restrict the number of Tx protocols used for transmission.
// This can reduce the required memory by the GGWave instance.
//
//	Note: this function does not affect instances that have already been created.
func TxToggleProtocol(protocol GGWaveProtocolType, state bool) {
	s := C.int(0)
	if state {
		s = C.int(1)
	}
	C.ggwave_txToggleProtocol(C.ggwave_ProtocolId(protocol), s)
}

// RxProtocolSetFreqStart Set freqStart for an Rx protocol
func RxProtocolSetFreqStart(protocol GGWaveProtocolType, freqStart int) {
	C.ggwave_rxProtocolSetFreqStart(C.ggwave_ProtocolId(protocol), C.int(freqStart))
}

// TxProtocolSetFreqStart Set freqStart for a Tx protocol
func TxProtocolSetFreqStart(protocol GGWaveProtocolType, freqStart int) {
	C.ggwave_txProtocolSetFreqStart(C.ggwave_ProtocolId(protocol), C.int(freqStart))
}
