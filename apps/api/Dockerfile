FROM golang:1.22 AS builder

WORKDIR /src

COPY go.mod ./
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api .

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/api /app/api

ENV PORT=8080

EXPOSE 8080

ENTRYPOINT ["/app/api"]
