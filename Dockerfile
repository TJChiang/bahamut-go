FROM golang:1.21.5
RUN mkdir -p /app
WORKDIR /app
COPY . .
RUN go mod download && go build -o bahamut
RUN apt-get update && apt-get install -y chromium
RUN go run github.com/playwright-community/playwright-go/cmd/playwright@v0.4001.0 install chromium
ENTRYPOINT ["./bahamut"]
