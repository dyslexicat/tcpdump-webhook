FROM golang:rc-alpine AS build 
# need to turn off CGO since otherwise there might be dynamic links
ENV CGO_ENABLED 0

RUN apk add git openssl

WORKDIR /usr/local/go/src/tcpdump-web
ADD . .
RUN go build ./cmd/tcpdump-web

FROM scratch
WORKDIR /app
COPY --from=build /usr/local/go/src/tcpdump-web/tcpdump-web .
COPY --from=build /usr/local/go/src/tcpdump-web/ssl ssl

EXPOSE 8443

CMD ["/app/tcpdump-web"]
