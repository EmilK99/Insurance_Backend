FROM golang as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o flightappbackend

FROM alpine

WORKDIR /flightapp

COPY --from=builder /app/flightappbackend ./
COPY --from=builder /app/config/flight_app.toml ./config/flight_app.toml

CMD /flightapp/flightappbackend

