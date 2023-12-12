FROM golang:latest as build

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o youtube_service ./main.go

FROM ubuntu:latest

WORKDIR /

COPY --from=build app/leaderboard leaderboard
COPY --from=build app/config.json /config.json

CMD ["./youtube_service"]

