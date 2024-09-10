FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod go.mod
# COPY go.sum go.sum

RUN go mod download

COPY . .

RUN go build -o bin/app ./cmd/app

FROM alpine

COPY --from=builder /app/build /build
COPY --from=builder /app/bin/app /app

EXPOSE 8080

CMD [ "/app" ]
