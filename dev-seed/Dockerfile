FROM golang:1.24.1-alpine AS builder

# Set working directory to match the expected project structure
WORKDIR /workspace

# Copy the entire workspace to preserve relative paths
COPY . .

# Build the application
WORKDIR /workspace/dev-seed
RUN go build -o seed-app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /workspace/dev-seed/seed-app .
RUN chmod +x /root/seed-app

ENV POSTGRES_URI=""
ENV SERVER_URL=""
ENV MAIL_PIT_URL=""

CMD ["/root/seed-app"]
