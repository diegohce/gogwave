vpath libggwave.so ggwave/build/src

.PHONY: all build install

all: build

build: ggwave libggwave.so

ggwave: 
	git clone https://github.com/ggerganov/ggwave --recursive

libggwave.so:
	mkdir -p ggwave/build && \
	cd ggwave/build && cmake .. -DGGWAVE_BUILD_EXAMPLES=OFF -DCMAKE_BUILD_TYPE=Release && \
	make

install: ggwave libggwave.so
	cd ggwave/build && make install && ldconfig