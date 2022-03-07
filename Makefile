install:
	go install -tags "json1" github.com/maxmoehl/tt/tt

install-grpc:
	go install -tags "json1 grpc" github.com/maxmoehl/tt/tt

build:
	go build -tags "json1" -o dist/tt github.com/maxmoehl/tt/tt

build-grpc:
	go build -tags "json1 grpc" -o dist/tt github.com/maxmoehl/tt/tt
