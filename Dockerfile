FROM golang:1.17 as builder

WORKDIR /app

# Copy the local code to the container image.
COPY . .

# Build the app.
RUN go build main.go

# Run the web service on container startup. Here use the local port specified above.
CMD ["./main"]