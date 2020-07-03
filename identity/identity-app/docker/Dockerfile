FROM golang:1.13-alpine as builder

WORKDIR /app

ENV GO111MODULE=on
ENV CGO_ENABLED=0

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -ldflags "-s -w" -o /tmp/app .

FROM alpine:latest

COPY ab_views ab_views
COPY static static
COPY users.example.json .
COPY --from=builder /tmp/app .

ENTRYPOINT ["./app"]
