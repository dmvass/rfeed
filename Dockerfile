# Builder image
FROM golang:alpine AS build

RUN apk update && apk add bash git gcc libc-dev

ADD . /build
WORKDIR /build

COPY . .

RUN go mod download
RUN go test -v ./...
RUN GOOS=linux GOARCH=amd64 go build -o /build/rfeed

# Application image
FROM scratch AS rfeed

ADD . /rfeed
ADD . /rfeed/db

WORKDIR /rfeed

COPY --from=build /build/rfeed /rfeed
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/rfeed/rfeed"]
