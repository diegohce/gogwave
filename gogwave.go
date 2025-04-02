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

type ggWaveParameters struct {
	payloadLength        int
	sampleRateInp        float32
	sampleRateOut        float32
	sampleRate           float32
	samplesPerFrame      int
	soundMarkerThreshold float32
	sampleFormatInp      GGWaveSampleFormatType
	sampleFormatOut      GGWaveSampleFormatType
	operatingMode        int
}

type GGWave struct {
	i                C.ggwave_Instance
	params           ggWaveParameters
	decodeBufferSize int
}

// New Create a new GGWave instance
//
//	Note: Make sure to deallocate the instance at the end by calling Close()
func New() *GGWave {
	c_gg_params := C.ggwave_getDefaultParameters()

	inst := C.ggwave_init(c_gg_params)

	params := ggWaveParameters{
		payloadLength:        int(c_gg_params.payloadLength),
		sampleRateInp:        float32(c_gg_params.sampleRateInp),
		sampleRateOut:        float32(c_gg_params.sampleRateOut),
		sampleRate:           float32(c_gg_params.sampleRate),
		samplesPerFrame:      int(c_gg_params.samplesPerFrame),
		soundMarkerThreshold: float32(c_gg_params.soundMarkerThreshold),
		sampleFormatInp:      GGWaveSampleFormatType(c_gg_params.sampleFormatInp),
		sampleFormatOut:      GGWaveSampleFormatType(c_gg_params.sampleFormatOut),
		operatingMode:        int(c_gg_params.operatingMode),
	}

	return &GGWave{i: inst, params: params, decodeBufferSize: 4096}
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
