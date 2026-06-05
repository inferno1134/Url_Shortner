This is repo for url shortner in go lang


# URL Shortener

A simple URL shortening service built using Go.

## Features

- URL Shortening
- URL Redirection
- In-memory Storage
- Docker Support

## Run Locally

go run main.go

## Run Using Docker

docker build -t url-shortner .

docker run -p 8080:8080 url-shortner

## API

POST /shorten

{
  "url":"https://google.com"
}