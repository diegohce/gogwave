# gogwave

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

	f, _ := os.OpenFile("out.wav", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)

	err = gg.EncodeToWav(f, waveform)
	if err != nil {
		panic(err)
	}
	f.Close()

	f, _ = os.OpenFile("out.wav", os.O_RDONLY, 0644)

	b, err := DecodeFromWav(f)
	if err != nil {
		panic(err)
	}

	payload := string(b)
    fmt.Println("payload:", payload)
}
```

## How to install

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

after building ggwave, from gogwave's root directory run:

```bash
go test .
```
tests should run fine.

Note that `ggwave` directory is safe to remove now.

## License

Licensed under the Apache 2.0 license. Copyright (c) 2025 Diego Cena

## ggwave License

ggwave is licensed under the MIT license. Copyright (c) 2020 Georgi Gerganov
