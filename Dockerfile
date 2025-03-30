FROM golang:1.21

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o vote-bot ./cmd

CMD ["./vote-bot"]