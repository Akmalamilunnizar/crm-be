FROM golang:1.22 as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./cmd/myapp/main.go

FROM alpine:latest as production
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /app/myapp .
COPY --from=build /app/uploads ./uploads
COPY --from=build /app/public ./public
CMD ["./myapp"]
EXPOSE 3001