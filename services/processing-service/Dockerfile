FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o processing-server ./cmd/server
RUN go build -o processing-worker ./cmd/worker

FROM alpine:3.19
RUN apk add --no-cache ffmpeg
WORKDIR /app
COPY --from=build /app/processing-server .
COPY --from=build /app/processing-worker .
EXPOSE 9093
CMD ["./processing-server"]