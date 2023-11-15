FROM golang:1.21 AS builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY main.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o app

EXPOSE 9092

FROM scratch
COPY --from=builder /app/app /
CMD ["/app"]