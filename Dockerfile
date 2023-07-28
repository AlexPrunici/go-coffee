FROM golang:1.20-alpine

WORKDIR /app

COPY ./src ./src
COPY go.mod ./
COPY go.sum ./

RUN go mod download

RUN go build -o ./server ./src/

EXPOSE 8080

CMD ["./server"]
