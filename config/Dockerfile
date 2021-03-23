FROM golang:1.15.8-stretch

ENV GO111MODULE=on

WORKDIR /src

COPY ./proto ./proto
COPY ./kit ./kit
COPY ./config/go.mod ./config/go.sum ./config/

WORKDIR /src/config

RUN go mod download -x

COPY ./config ./

WORKDIR /src/config/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /main .

FROM alpine:3.12.3

WORKDIR /usr/local/bin/config

RUN apk --no-cache add ca-certificates
COPY --from=0 /main ./main
COPY ./config/config.docker.yml ./config.yml

ENTRYPOINT ["./main"]
