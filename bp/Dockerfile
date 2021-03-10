FROM golang:1.15.8-stretch

ENV GO111MODULE=on

WORKDIR /src

COPY ./proto ./proto
COPY ./kit ./kit
COPY ./bp/go.mod ./bp/go.sum ./bp/

WORKDIR /src/bp

RUN go mod download -x

COPY ./bp ./

WORKDIR /src/bp/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /main .

FROM alpine:3.12.3

WORKDIR /usr/local/bin/bp

RUN apk --no-cache add ca-certificates
COPY --from=0 /main ./main
COPY ./bp/bpmn ./bpmn

ENTRYPOINT ["./main"]
