FROM golang:latest

WORKDIR /app

RUN go get go.mongodb.org/mongo-driver/mongo
RUN go get github.com/streadway/amqp

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Command to run the executable
CMD ["./main"]