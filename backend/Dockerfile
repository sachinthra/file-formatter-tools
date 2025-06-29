FROM golang:1.22 as builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 go build -o app ./cmd/main.go

FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    libjpeg-turbo8 libpng16-16 libwebp7 libtiff5 \
    libopencv-core-dev libopencv-imgproc-dev libopencv-imgcodecs-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/internal ./internal
CMD ["./app"]
EXPOSE 8081