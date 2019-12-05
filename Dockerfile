FROM golang:latest

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o htsref ./cmd
EXPOSE 3000

CMD ["./htsref"]
