# Builder image
FROM golang:alpine AS build

ADD . /build
WORKDIR /build

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o /build/rfeed .

# Application image
FROM scratch AS rfeed

ADD . /rfeed
ADD . /rfeed/db

WORKDIR /rfeed

COPY --from=build /build/rfeed /rfeed
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/rfeed/rfeed"]
