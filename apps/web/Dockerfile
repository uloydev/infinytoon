FROM golang:1.23 AS builder

WORKDIR /app
COPY . .
RUN go mod download
WORKDIR /app/apps/web
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /tmp/app


FROM scratch AS runner
# RUN apk --no-cache add ca-certificates curl htop tzdata
COPY --from=builder /tmp/app /
EXPOSE 8080
ENTRYPOINT ["/app"]