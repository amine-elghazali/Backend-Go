# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the local code into the container at the working directory
COPY . .

# Download Go modules
RUN go mod download

# Build the Go app
RUN go build -o main .

# Expose the port that your application will run on
EXPOSE 3000

# Command to run the executable
CMD ["./main"]
