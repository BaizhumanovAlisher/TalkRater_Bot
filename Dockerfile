#Compiling program
FROM golang:1.21.9 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd/api/ -o /talk_rater_bot

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Create certificates for SSL/TLS connections
FROM alpine AS certificates-builder
WORKDIR /app

RUN apk add ca-certificates && update-ca-certificates

# Deploy the application binary into a lean image
FROM scratch AS build-release-stage

WORKDIR /

COPY --from=build-stage /talk_rater_bot /talk_rater_bot
COPY --from=certificates-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["./talk_rater_bot"]