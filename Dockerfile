FROM golang:1.21-alpine AS builder

RUN apk add --no-cache \
    git \
    make \
    gcc \
    musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -a -installsuffix cgo \
    -o /bin/app ./cmd/main.go

FROM scratch AS runtime

COPY --from=builder /bin/app /bin/app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV APP_ENV=production \
    APP_PORT=8080

ENTRYPOINT ["/bin/app"]

EXPOSE 8080