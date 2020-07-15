FROM golang:1.13-alpine as builder

WORKDIR /app

RUN apk add --no-cache git
RUN go get github.com/derekparker/delve/cmd/dlv
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait

ENV GO111MODULE=on
ENV CGO_ENABLED=0

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -gcflags "all=-N -l" -o /tmp/app .

FROM alpine:latest

# Allow delve to run on Alpine based containers.
RUN apk add --no-cache libc6-compat
COPY --from=builder /go/bin/dlv /
# Allows waiting on other services
COPY --from=builder /wait /wait
# Main application
COPY --from=builder /tmp/app .

RUN chmod +x /wait
CMD /wait && if [ $DEBUG == "1" ]; then /dlv --listen=:40000 --headless=true --api-version=2 exec /app; else ./app; fi