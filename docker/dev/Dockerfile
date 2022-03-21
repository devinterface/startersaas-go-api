FROM golang:1.17

WORKDIR /go/src/app/

COPY . .

RUN go mod download -x

RUN go get github.com/githubnemo/CompileDaemon

EXPOSE 3000

ENTRYPOINT CompileDaemon --build="go build main.go" --command="./main"