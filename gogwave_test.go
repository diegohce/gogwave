package gogwave

import (
	"os"
	"testing"
)

func TestEncodeDecode(t *testing.T) {

	gg := New()
	defer gg.Close()

	waveform, err := gg.Encode([]byte("hola"), ProtocolAudibleNormal, 50)
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

	p := NewGGwaveParameters()

	gg := NewWhithParams(p)
	defer gg.Close()

	waveform, err := gg.Encode([]byte("hola"), ProtocolAudibleNormal, 50)
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

func TestToWav(t *testing.T) {
	gg := New()
	defer gg.Close()

	waveform, err := gg.Encode([]byte("hola"), ProtocolAudibleNormal, 50)
	if err != nil {
		t.Fatal(err)
	}

	f, _ := os.OpenFile("out.wav", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)

	err = gg.EncodeToWav(f, waveform)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	f, _ = os.OpenFile("out.wav", os.O_RDONLY, 0644)

	b, err := DecodeFromWav(f)
	if err != nil {
		t.Fatal(err)
	}

	payload := string(b)
	if payload != "hola" {
		t.Fatalf("got %s want hola", payload)
	}

}

func TestRxDurationFrames(t *testing.T) {
	gg := New()
	defer gg.Close()

	wafeform, err := gg.Encode([]byte("hola"), ProtocolAudibleNormal, 50)
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
	TxToggleProtocol(ProtocolUltrasoundFast, false)
	TxToggleProtocol(ProtocolUltrasoundFast, true)
}

func TestRxToggleProtocol(t *testing.T) {
	RxToggleProtocol(ProtocolUltrasoundFast, false)
	RxToggleProtocol(ProtocolUltrasoundFast, true)
}

func TestRxProtocolSetStartFreq(t *testing.T) {
	RxProtocolSetFreqStart(ProtocolAudibleNormal, 1500)
}

func TestTxProtocolSetStartFreq(t *testing.T) {
	TxProtocolSetFreqStart(ProtocolAudibleNormal, 1500)
}

func TestSetLogFileNIL(t *testing.T) {
	SetLogFile(nil)
}

func TestSetLogFileToStderr(t *testing.T) {
	SetLogFile(os.Stderr)
}

func TestSetLogFileToFile(t *testing.T) {
	f, _ := os.OpenFile("/tmp/gogwave_test.log", os.O_CREATE|os.O_APPEND, 0644)
	SetLogFile(f)
	f.Close()
}

func TestDecodeBufferSize(t *testing.T) {

	gg := New()
	defer gg.Close()

	gg.SetDecodeBufferSize(8192)

	bs := gg.DecodeBufferSize()

	if bs != 8192 {
		t.Errorf("got %d want 8192", bs)
	}
}
