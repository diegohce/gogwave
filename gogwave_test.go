package gogwave_test

import (
	"os"
	"testing"

	"github.com/diegohce/gogwave"
)

func TestEncodeDecode(t *testing.T) {

	gg := gogwave.New()
	defer gg.Close()

	waveform, err := gg.Encode([]byte("hola"), gogwave.ProtocolAudibleNormal, 50)
	if err != nil {
		t.Fatal(err)
	}

	b, err := gg.Decode(waveform)
	if err != nil {
		t.Fatal(err)
	}

	payload := string(b)
	if payload != "hola" {
		t.Fatalf("got %s want hola", payload)
	}
}

func TestNewWithParams(t *testing.T) {

	p := gogwave.NewGGwaveParameters()

	gg := gogwave.NewWhithParams(p)
	defer gg.Close()

	waveform, err := gg.Encode([]byte("hola"), gogwave.ProtocolAudibleNormal, 50)
	if err != nil {
		t.Fatal(err)
	}

	b, err := gg.Decode(waveform)
	if err != nil {
		t.Fatal(err)
	}

	payload := string(b)
	if payload != "hola" {
		t.Fatalf("got %s want hola", payload)
	}
}

func TestRxDurationFrames(t *testing.T) {
	gg := gogwave.New()
	defer gg.Close()

	wafeform, err := gg.Encode([]byte("hola"), gogwave.ProtocolAudibleNormal, 50)
	if err != nil {
		t.Fatal(err)
	}

	_, err = gg.Decode(wafeform)
	if err != nil {
		t.Fatal(err)
	}

	df := gg.RxDurationFrames()
	if df != 52 {
		t.Errorf("got %d want 52", df)
	}
}

func TestTxToggleProtocol(t *testing.T) {
	gogwave.TxToggleProtocol(gogwave.ProtocolUltrasoundFast, false)
	gogwave.TxToggleProtocol(gogwave.ProtocolUltrasoundFast, true)
}

func TestRxToggleProtocol(t *testing.T) {
	gogwave.RxToggleProtocol(gogwave.ProtocolUltrasoundFast, false)
	gogwave.RxToggleProtocol(gogwave.ProtocolUltrasoundFast, true)
}

func TestRxProtocolSetStartFreq(t *testing.T) {
	gogwave.RxProtocolSetFreqStart(gogwave.ProtocolAudibleNormal, 1500)
}

func TestTxProtocolSetStartFreq(t *testing.T) {
	gogwave.TxProtocolSetFreqStart(gogwave.ProtocolAudibleNormal, 1500)
}

func TestSetLogFileNIL(t *testing.T) {
	gogwave.SetLogFile(nil)
}

func TestSetLogFileToStderr(t *testing.T) {
	gogwave.SetLogFile(os.Stderr)
}

func TestSetLogFileToFile(t *testing.T) {
	f, _ := os.OpenFile("/tmp/gogwave_test.log", os.O_CREATE|os.O_APPEND, 0644)
	gogwave.SetLogFile(f)
	f.Close()
}

func TestDecodeBufferSize(t *testing.T) {

	gg := gogwave.New()
	defer gg.Close()

	gg.SetDecodeBufferSize(8192)

	bs := gg.DecodeBufferSize()

	if bs != 8192 {
		t.Errorf("got %d want 8192", bs)
	}
}
