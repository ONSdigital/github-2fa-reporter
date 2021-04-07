FROM golang:alpine AS builder
RUN apk update && apk add --no-cache make
WORKDIR $GOPATH/src/github.com/ONSdigital/github-2fa-reporter/
COPY . .
RUN make

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/ONSdigital/github-2fa-reporter/build/linux-amd64/bin/github2fareporter /github2fareporter

ENTRYPOINT [ "/github2fareporter" ]