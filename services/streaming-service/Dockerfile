FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o streaming-server ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/streaming-server .
EXPOSE 8083
CMD ["./streaming-server"]