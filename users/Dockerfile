FROM golang:1.15.8-stretch

ENV GO111MODULE=on

WORKDIR /src

COPY ./proto ./proto
COPY ./kit ./kit
COPY ./users/go.mod ./users/go.sum ./users/

WORKDIR /src/users

RUN go mod download -x

COPY ./users ./

WORKDIR /src/users/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /main .

FROM alpine:3.12.3

WORKDIR /usr/local/bin/users

RUN apk --no-cache add ca-certificates
COPY --from=0 /main ./main

ENTRYPOINT ["./main"]
