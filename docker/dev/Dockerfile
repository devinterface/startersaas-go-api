FROM golang:1.21

WORKDIR /go/src/app/

COPY . .

RUN apt update

RUN go mod download -x

RUN go install github.com/cosmtrek/air@latest

EXPOSE 3000

CMD ["air", "-c", ".air.toml"]
