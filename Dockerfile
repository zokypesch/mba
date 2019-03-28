
FROM golang:1.7

# Set go bin which doesn't appear to be set already.
ENV GOBIN /go/bin

# build directories
RUN mkdir /app
RUN mkdir /go/src/app
ADD . /go/src/app
WORKDIR /go/src/app

# Go dep!
# RUN go get -u github.com/golang/dep/...
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure

EXPOSE 80

# Build my app
RUN go build -o /app/main .
CMD ["/app/main"]