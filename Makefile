VERSION ?= latest

build-image: build
	docker build -t fkarakas/helloworld-grpc:$(VERSION) .

push-image: build-image
	docker push fkarakas/helloworld-grpc:$(VERSION)

build:
	go build -o hello-server ./server
	go build -o hello-client ./client

run-client: build
	./hello-client

run-server: build
	./hello-server
