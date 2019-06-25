FROM golang:1.12-alpine as builder

RUN apk --no-cache add git bzr bind-tools && \
    apk --update add alpine-sdk && \
    rm -rf /var/cache/apk/*l

WORKDIR /go/grpc/helloworld

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o hello-server ./server
RUN CGO_ENABLED=0 GOOS=linux go build -o hello-client ./client

FROM alpine

RUN apk --no-cache add git bzr bind-tools && \
    apk --update add alpine-sdk && \
    rm -rf /var/cache/apk/*l

WORKDIR /bin/

COPY --from=builder /go/grpc/helloworld/hello-server server
COPY --from=builder /go/grpc/helloworld/hello-client client

#USER nobody

CMD /bin/${MODE}
