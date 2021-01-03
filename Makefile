all: clean test build

run:
	CGO_ENABLED=1 go run .

build:
	go mod download
	CGO_ENABLED=1 go build

clean:
	rm -rf ./prayeralarm

test:
	CGO_ENABLED=1 go test -race . -v

docker-up:
	docker build -t prayeralarm:alpine .
	docker run --rm -it --name prayeralarm \
		-e PULSE_SERVER=host.docker.internal \
		-v ~/.config/pulse:/root/.config/pulse \
		prayeralarm:alpine

docker-down:
	docker stop prayeralarm
