FROM golang:1.20-alpine as build-stage

RUN apk --no-cache add ca-certificates

WORKDIR /go/src/btc-price

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /btc-price ./cmd/price_check

#
# final build stage
#
FROM scratch

# Copy ca-certs for app web access
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-stage /btc-price /btc-price

# app uses port 8080
EXPOSE 8080

ENTRYPOINT ["/btc-price"]