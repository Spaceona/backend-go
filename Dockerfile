FROM golang:latest

WORKDIR /app

COPY  src/go.mod .
RUN go mod download

COPY ./src .
COPY ./migrations ./migrations
#TOOD remove me
#COPY .env .env
# DB now lives in ./data (mounted volume); no longer baked into the image
COPY ./firmware ./firmware
RUN go build -o /spacesona-go
RUN apt-get --allow-unauthenticated update
RUN apt install sqlite3 # todo remove later

EXPOSE 3001

CMD ["/spacesona-go"]
