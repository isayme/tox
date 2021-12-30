FROM golang:1.15.3-alpine as builder
WORKDIR /app

ARG APP_NAME
ENV APP_NAME ${APP_NAME}
ARG APP_VERSION
ENV APP_VERSION ${APP_VERSION}

COPY . .
RUN mkdir -p ./dist  \
  && GO111MODULE=on GOPROXY=https://goproxy.io,direct go mod download \
  && go build -ldflags "-X github.com/isayme/go-toh2/util.Name=${APP_NAME} \
  -X github.com/isayme/go-toh2/util.Version=${APP_VERSION}" \
  -o ./dist/toh2 main.go

FROM alpine
WORKDIR /app

ARG APP_NAME
ENV APP_NAME ${APP_NAME}
ARG APP_VERSION
ENV APP_VERSION ${APP_VERSION}

# default config file
ENV CONF_FILE_PATH=/etc/toh2.yaml

COPY --from=builder /app/dist/toh2 /app/toh2

CMD ["/app/toh2", "-h"]
