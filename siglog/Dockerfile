FROM golang:1.23
WORKDIR /usr/src/goapp
COPY . .
RUN go mod download && go mod verify
RUN go build -o app .
CMD ["/usr/src/goapp/app"]
