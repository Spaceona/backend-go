FROM golang:latest

WORKDIR /app

COPY  src/go.mod .
RUN go mod download

COPY ./src .
COPY ./migrations ./migrations
#TOOD remove me
COPY .env .env
COPY ./test.db ./test.db

RUN go build -o /spacesona-go

EXPOSE 3001

CMD ["/spacesona-go"]
