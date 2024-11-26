FROM golang:1.21 as build

WORKDIR /app

# Download Go modules
COPY go.mod go.sum main.go ./
RUN go mod download

COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /api
FROM scratch
COPY --from=build /api /api

CMD ["/api"]
