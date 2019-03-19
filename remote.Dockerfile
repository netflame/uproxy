FROM golang:1.12-alpine as builder
WORKDIR /build
COPY ./*go* /build/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o uproxy


FROM scratch
WORKDIR /uproxy
COPY --from=builder /build/uproxy .
COPY ./*toml .

EXPOSE 6001
ENTRYPOINT [ "./uproxy" ]