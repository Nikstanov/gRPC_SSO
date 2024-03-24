FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN go build -o sso_service_grpc ./cmd/sso_service_grpc/main.go

FROM alpine

WORKDIR /app

COPY config/local.yaml .

COPY --from=builder /build/sso_service_grpc /app/sso_service_grpc

CMD ["./sso_service_grpc", "--config=./local.yaml"]