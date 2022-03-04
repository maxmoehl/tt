default: install

install:
	go install -tags "json1" github.com/maxmoehl/tt/tt

build:
	go build -tags "json1" -o dist/tt github.com/maxmoehl/tt/tt
