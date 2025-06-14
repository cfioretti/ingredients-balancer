FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apk add --no-cache protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

RUN mkdir -p pkg/infrastructure/grpc/proto/generated
RUN protoc --go_out=pkg/infrastructure/grpc/proto/generated --go-grpc_out=pkg/infrastructure/grpc/proto/generated pkg/infrastructure/grpc/proto/ingredients_balancer.proto

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ingredients-balancer cmd/main.go

FROM alpine:3.21

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/ingredients-balancer .

EXPOSE 50052

ENV PORT=50052

RUN adduser -D -g '' appuser
USER appuser

CMD ["./ingredients-balancer"]
