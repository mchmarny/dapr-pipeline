FROM golang:1.14.2 as builder

WORKDIR /src/
COPY . /src/

ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -tags netgo \
    -ldflags "-w -extldflags -static" \
    -mod vendor \
    -o ./service \
    .

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /src/service .
COPY --from=builder /src/resource ./resource/

ENTRYPOINT ["./service"]
