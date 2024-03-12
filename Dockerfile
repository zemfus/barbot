FROM golang:1.21.1 as env

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download -x

FROM env as build

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bot cmd/main.go

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /build/bot /bin/bot
ENTRYPOINT ["/bin/bot", "conf.yaml"]