# ca-certs stage
FROM alpine:latest AS certs
RUN apk --update add ca-certificates

# build stage
FROM golang:latest AS build
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o smart-device-relay .

# final stage
FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /app/smart-device-relay /smart-device-relay
ENTRYPOINT ["/smart-device-relay"]