FROM golang:1.26-alpine

WORKDIR /app

COPY go.mod ./
COPY . .

RUN go build -o url-shortner .

EXPOSE 8080

CMD ["./url-shortner"]