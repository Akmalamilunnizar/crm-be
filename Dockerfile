FROM golang:1.22 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./cmd/myapp/main.go

FROM alpine:latest as production
WORKDIR /app
COPY --from=build /app/myapp .
COPY --from=build /app/.env .
CMD ["./myapp"]
EXPOSE 3001