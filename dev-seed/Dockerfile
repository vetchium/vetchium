FROM golang:1.24.1-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o seed-app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/seed-app .
RUN chmod +x /root/seed-app

ENV POSTGRES_URI=""
ENV SERVER_URL=""
ENV MAIL_PIT_URL=""

CMD ["/root/seed-app"]
