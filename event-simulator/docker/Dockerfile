FROM golang:1.23-alpine

WORKDIR /app

COPY /src/ .

# Build the executable
RUN CGO_ENABLED=0 GOOS=linux go build -o /simulator

# Run the executable
CMD ["/simulator"]