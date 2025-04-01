FROM golang:1.24 AS builder
WORKDIR /app

COPY cmd cmd
COPY internal internal
COPY go.mod .

RUN CGO_ENABLED=0 go build -o /bin/writy cmd/main.go

FROM scratch
COPY --from=builder /bin/writy /bin/writy
ENTRYPOINT [ "/bin/writy" ]
