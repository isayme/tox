FROM golang:1.22.3-alpine as builder
WORKDIR /app

ARG APP_NAME
ENV APP_NAME ${APP_NAME}
ARG APP_VERSION
ENV APP_VERSION ${APP_VERSION}

COPY . .
RUN mkdir -p ./dist && GO111MODULE=on go mod download
RUN go build -ldflags "-X github.com/isayme/tox/util.Name=${APP_NAME} \
    -X github.com/isayme/tox/util.Version=${APP_VERSION}" \
    -o ./dist/tox main.go

FROM alpine
WORKDIR /app

ARG APP_NAME
ENV APP_NAME ${APP_NAME}
ARG APP_VERSION
ENV APP_VERSION ${APP_VERSION}

# default config file
ENV CONF_FILE_PATH=/etc/tox.yaml

COPY --from=builder /app/dist/tox /app/tox

CMD ["/app/tox", "-h"]
