FROM golang:1.18-rc-alpine AS build-env
WORKDIR /src

RUN apk add --no-cache
COPY . .
RUN go build -o dispatcher ./cmd/dispatcher/*.go
RUN go build -o server ./cmd/*.go

FROM build-env
WORKDIR /app
COPY --from=build-env /src/dispatcher dispatcher
COPY --from=build-env /src/server server
# ENTRYPOINT ./loadbalancer