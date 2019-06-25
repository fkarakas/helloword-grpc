VERSION ?= latest

build-image:
	docker build -t helloworld:$(VERSION) .

build:
	go build -o hello-server ./server
	go build -o hello-client ./client

run-client: build
	./hello-client

run-server: build
	./hello-server