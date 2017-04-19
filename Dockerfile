FROM golang:alpine
COPY . /go/src/app
WORKDIR /go/src/app
RUN go-wrapper download && go-wrapper install
ENTRYPOINT ["go-wrapper", "run"]
