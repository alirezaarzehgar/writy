FROM golang:1.24 AS builder
WORKDIR /app

COPY cmd cmd
COPY internal internal
COPY libwrity libwrity
COPY cache cache
COPY Makefile VERSION go.mod go.sum ./

RUN go mod download && \
    CGO_ENABLED=0 go build -o /bin/writy cmd/main.go

FROM scratch
COPY --from=builder /bin/writy /bin/writy
ENTRYPOINT [ "/bin/writy" ]
