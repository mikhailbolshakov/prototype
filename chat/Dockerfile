FROM golang:1.15.8-stretch

ENV GO111MODULE=on

WORKDIR /src

COPY ./proto ./proto
COPY ./kit ./kit
COPY ./chat/go.mod ./chat/go.sum ./chat/

WORKDIR /src/chat

RUN go mod download -x

COPY ./chat ./

WORKDIR /src/chat/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /main .

FROM alpine:3.12.3

WORKDIR /usr/local/bin/chat

RUN apk --no-cache add ca-certificates
COPY --from=0 /main ./main

ENTRYPOINT ["./main"]
