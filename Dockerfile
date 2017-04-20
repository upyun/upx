FROM golang:alpine
COPY . /go/src/upx
WORKDIR /go/src/upx
RUN go get -d -v && go install -v
CMD ["upx"]
