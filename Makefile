

.PHONY: all build install

all: build

ggwave: 
	git clone https://github.com/ggerganov/ggwave --recursive

build: ggwave
	mkdir -p ggwave/build && \
	cd ggwave/build && cmake .. -DGGWAVE_BUILD_EXAMPLES=OFF -DCMAKE_BUILD_TYPE=Release && \
	make

install: ggwave build
	cd ggwave/build && make install && ldconfig