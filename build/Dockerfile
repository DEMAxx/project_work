# Собираем в гошке
FROM golang:1.23 AS build

ENV BIN_FILE=/opt/previewer/previewer-app
ENV CODE_DIR=/go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/calendar/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="previewer"
LABEL MAINTAINERS="otdupli@gmail.com"

ENV BIN_FILE=/opt/previewer/previewer-app
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE=/etc/previewer/config.toml

COPY ./configs/config.toml ${CONFIG_FILE}

RUN chmod +x ${BIN_FILE}

CMD ["/opt/previewer/previewer-app", "-config", "/etc/previewer/config.toml"]