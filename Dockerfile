FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o youtube_service ./main.go

CMD ["./youtube_service"]