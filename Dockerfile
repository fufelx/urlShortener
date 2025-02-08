FROM golang:1.23

WORKDIR /urlShortener

COPY . .

RUN go mod tidy
RUN go build -o main /urlShortener/cmd/main/main.go

RUN apt-get update && apt-get install -y bash

EXPOSE 3030

CMD ["bash", "-c", "set -a && [ -f .env ] && source .env && set +a && ./main"]
