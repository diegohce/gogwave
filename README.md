[![Go Reference](https://pkg.go.dev/badge/github.com/diegohce/gogwave.svg)](https://pkg.go.dev/github.com/diegohce/gogwave) [![Go Report Card](https://goreportcard.com/badge/github.com/diegohce/gogwave)](https://goreportcard.com/report/github.com/diegohce/gogwave) 

# GoGWave

The gogwave package provides Go language bindings for the [ggwave](https://github.com/ggerganov/ggwave) C++ library.

**Note:** At the moment, gogwave supports linux systems only.

For further detail see the [GoDocs](https://pkg.go.dev/github.com/diegohce/gogwave)

## How to use

This example encodes a payload (`[]byte`) into a waveform, stores the waveform as a .wav file, and decodes the waveform back into a payload.

```Go
func main() {
	gg := New()
	defer gg.Close()

	waveform, err := gg.Encode([]byte("hola"), ProtocolAudibleNormal, 50)
	if err != nil {
		panic(err)
	}

	b, err := gg.Decode(waveform)
	if err != nil {
		panic(err)
	}

	payload := string(b)
    fmt.Println("payload:", payload)
}
```

## How to build & install ggwave

To install ggwave library follow the [ggwave building/installing instructions](https://github.com/ggerganov/ggwave?tab=readme-ov-file#building)

An alternative is using GoGWave's `Makefile`

Clone this repo:

```bash
git clone https://github.com/diegohce/gogwave
cd gogwave
```

Build (linux):

```bash
make
```

Install: 

```bash
sudo make install
```

after building ggwave, from GoGWave's root directory run:

```bash
go test .
```
tests should run fine.

**Note:** `ggwave` directory is safe to remove after building and installing.

## License

GoGWave is licensed under the Apache 2.0 license. Copyright (c) 2025 Diego Cena

## ggwave License

ggwave is licensed under the MIT license. Copyright (c) 2020 Georgi Gerganov
