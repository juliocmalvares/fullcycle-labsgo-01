FROM golang:1.24.2-alpine

WORKDIR /app

COPY go.mod go.sum ./
COPY .env ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]