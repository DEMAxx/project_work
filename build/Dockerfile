# Собираем в гошке
FROM golang:1.24 AS build

ENV BIN_FILE=/opt/previewer/previewer-app
ENV CODE_DIR=/go/src/
ENV CC=gcc

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

RUN apt-get update && \
    apt-get -qq install -y libvips-dev

# Собираем статический бинарник Go
ARG LDFLAGS
RUN CGO_ENABLED=1 PKG_CONFIG_PATH="/usr/lib/pkgconfig" go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/previewer/*

# На выходе тонкий образ
FROM build AS base

ENV CGO_ENABLED=1
ENV TARGETOS=linux
ENV TARGETARCH=amd64

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="previewer"
LABEL MAINTAINERS="otdupli@gmail.com"

ENV BIN_FILE=/opt/previewer/previewer-app
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE=/etc/previewer/config.toml

COPY ./configs/config.toml ${CONFIG_FILE}

RUN chmod +x ${BIN_FILE}

CMD ["/opt/previewer/previewer-app", "-config", "/etc/previewer/config.toml"]