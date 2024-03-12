FROM golang:1.21.1 as env

WORKDIR /build

COPY go.mod .
COPY go.sum .
COPY deployments/local_config.yaml .

RUN go mod download -x

FROM env as build

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bot cmd/main.go

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /build/bot /bin/bot
COPY --from=build /build/local_config.yaml /bin/
ENTRYPOINT ["/bin/bot", "/bin/local_config.yaml"]