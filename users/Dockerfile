# Build the application from source
FROM golang:1.22.5 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY .env .env

RUN CGO_ENABLED=0 GOOS=linux go build -o /users

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /users /users
COPY --from=build-stage /app/.env .env

EXPOSE 8000

ENTRYPOINT ["/users"]

## RUN THIS CMD >> docker build -t users:multistage -f Dockerfile .