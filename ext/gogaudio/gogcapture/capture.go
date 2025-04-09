package gogcapture

import (
	"fmt"

	"github.com/gen2brain/malgo"

	"github.com/diegohce/gogwave/ext/gogaudio"
)

type RawCapture struct {
	ctx             *malgo.AllocatedContext
	device          *malgo.Device
	audioBuffer     []byte
	capturedSamples uint32
}

func newRawCapture(_ any) (gogaudio.AudioCapture, error) {

	return &RawCapture{}, nil
}

func (wc *RawCapture) Close() error {
	if wc.device != nil {
		wc.device.Uninit()
		wc.device = nil
	}
	if wc.ctx != nil {
		wc.ctx.Uninit()
		wc.ctx.Free()
		wc.ctx = nil
	}
	return nil
}

func (wc *RawCapture) CapturedBuffer() ([]byte, error) {
	return wc.audioBuffer, nil
}

func (wc *RawCapture) Capture() error {
	var err error
	wc.ctx, err = malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		fmt.Printf("LOG <%v>\n", message)
	})
	if err != nil {
		return err
	}

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS32
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = 48000
	deviceConfig.Alsa.NoMMap = 1

	wc.capturedSamples = 0
	wc.audioBuffer = make([]byte, 0)

	sizeInBytes := uint32(malgo.SampleSizeInBytes(deviceConfig.Capture.Format))
	onRecvFrames := func(pSample2, pSample []byte, framecount uint32) {

		sampleCount := framecount * deviceConfig.Capture.Channels * sizeInBytes

		newCapturedSampleCount := wc.capturedSamples + sampleCount

		wc.audioBuffer = append(wc.audioBuffer, pSample...)

		wc.capturedSamples = newCapturedSampleCount

	}

	captureCallbacks := malgo.DeviceCallbacks{
		Data: onRecvFrames,
	}
	wc.device, err = malgo.InitDevice(wc.ctx.Context, deviceConfig, captureCallbacks)
	if err != nil {
		return err
	}

	return wc.device.Start()
}

func init() {
	gogaudio.RegisterCapture("raw", newRawCapture)
}
