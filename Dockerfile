FROM golang:1.18-alpine as build
WORKDIR /go/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./bin/server ./cmd/server/main.go

FROM alpine:3.16
RUN adduser --disabled-password --gecos "" --uid 1000 default
USER default
WORKDIR /app
COPY --from=build --chown=default /go/app/bin ./
ENTRYPOINT [ "./server" ]
