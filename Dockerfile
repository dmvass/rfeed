# Builder image
FROM golang:1.12-alpine AS build
 
RUN apk add bash git gcc libc-dev

ADD . /build
WORKDIR /build

COPY . .

RUN go mod download
RUN go test -v ./...
RUN GOOS=linux GOARCH=amd64 go build -o /build/rfeed

# Application image
FROM alpine AS rfeed

ADD . /rfeed
WORKDIR /rfeed

COPY --from=build /build/rfeed /rfeed

RUN chown -R nobody:nobody /rfeed

USER nobody:nobody

ENTRYPOINT ["rfeed"]
