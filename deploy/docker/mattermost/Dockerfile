FROM golang:1.15.8-stretch as server-build

ENV GO111MODULE=on

# build server
WORKDIR /src

COPY ./mattermost-server/go.mod ./mattermost-server/go.sum ./mattermost-server/

WORKDIR /src/mattermost-server

RUN go mod download -x

COPY ./mattermost-server ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /mattermost ./cmd/mattermost

# build plugin

FROM golang:1.15.8-stretch as plugin-build

ENV GO111MODULE=on

WORKDIR /src

COPY ./mattermost-plugin/go.mod ./mattermost-plugin/go.sum ./mattermost-plugin/

WORKDIR /src/mattermost-plugin

RUN go mod download -x

COPY ./mattermost-plugin ./

RUN make dist

FROM alpine:3.12.3

ENV PATH="/mattermost/bin:${PATH}"

#ARG PUID=2000
#ARG PGID=2000

WORKDIR /mattermost

RUN apk --no-cache add \
    ca-certificates \
	curl \
	jq \
	libc6-compat \
	libffi-dev \
	libcap \
	linux-headers \
	mailcap \
	netcat-openbsd \
	xmlsec-dev \
	tzdata \
    && mkdir -p plugins config client/plugins bin fonts i18n templates src/plugins src/config

## copy prepackaged config and plugins to "src" dir, so it works for both docker and kuber
## as kuber removes all the content from the mounted folder we need to copy it after pvc created

COPY --from=server-build /src/mattermost-server/fonts/ ./fonts/
COPY --from=server-build /src/mattermost-server/i18n/ ./i18n/
COPY --from=server-build /src/mattermost-server/templates/ ./templates/
COPY --from=server-build /mattermost ./bin/mattermost
COPY --from=plugin-build /src/mattermost-plugin/dist/ ./src/plugins/
COPY ./config.json ./src/config/config.json
COPY ./mattermost-webapp/dist/ ./client/
COPY ./entrypoint.sh ./

#RUN addgroup -g ${PGID} mattermost \
#    && adduser -D -u ${PUID} -G mattermost -h /mattermost -D mattermost \
#    && chown -R mattermost:mattermost /mattermost /mattermost/config /mattermost/plugins /mattermost/client/plugins \
#    && setcap cap_net_bind_service=+ep /mattermost/bin/mattermost

RUN setcap cap_net_bind_service=+ep /mattermost/bin/mattermost

#USER mattermost

ENTRYPOINT ["./entrypoint.sh"]
CMD ["mattermost"]

EXPOSE 8065
