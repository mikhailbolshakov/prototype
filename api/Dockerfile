FROM golang:1.15.8-stretch

ENV GO111MODULE=on

WORKDIR /src

COPY ./proto ./proto
COPY ./kit ./kit
COPY ./api/go.mod ./api/go.sum ./api/

WORKDIR /src/api

RUN go mod download -x

COPY ./api ./

WORKDIR /src/api/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /main .

FROM alpine:3.12.3

WORKDIR /usr/local/bin/api

RUN apk --no-cache add ca-certificates
COPY --from=0 /main ./main

ENTRYPOINT ["./main"]
