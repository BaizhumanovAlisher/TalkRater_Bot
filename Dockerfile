#Compiling program
FROM golang:1.21.9 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd/main/ -o /talk_rater_bot

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM alpine AS build-release-stage

RUN apk add --no-cache tzdata # it is for time.Location
RUN apk add --no-cache postgresql-client

WORKDIR /

COPY --from=build-stage /talk_rater_bot /talk_rater_bot

CMD ["./talk_rater_bot"]