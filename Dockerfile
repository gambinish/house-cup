# DEVELOPMENT
# docker build -t yourusername/yourapp:dev --target=development .

FROM golang:1.20-alpine as development
WORKDIR /app

# add some necessary packages
RUN apk update && \
    apk add libc-dev && \
    apk add gcc && \
    apk add make

# prevent the re-installation of vendors at every change in the source code
COPY ./go.mod go.sum ./
RUN go mod download && go mod verify

# Install Compile Daemon for go. We'll use it to watch changes in go files
RUN go install github.com/githubnemo/CompileDaemon@latest

# Copy and build the app
COPY . .
COPY ./entrypoint.sh /entrypoint.sh

# wait-for-it requires bash, which alpine doesn't ship with by default. Use wait-for instead
ADD https://raw.githubusercontent.com/eficode/wait-for/v2.1.0/wait-for /usr/local/bin/wait-for
RUN chmod +rx /usr/local/bin/wait-for /entrypoint.sh

ENTRYPOINT [ "sh", "/entrypoint.sh" ]

# PRODUCTION
# docker build -t ngambino0192/house-cup-api:prod --target=production .
# docker run -p 9090:9090 ngambino0192/house-cup-api:prod

# Use a minimal base image for the production stage
FROM golang:1.20-alpine as production

# Set the working directory
WORKDIR /app

# Install necessary packages
RUN apk update && \
    apk add --no-cache libc-dev gcc make

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy a different environment file
COPY .env.prod /app/.env

# Copy the application source code
COPY . .

# Build the Go application
RUN go build -o app

# Expose the port on which the application will run
EXPOSE 9090

# Set the default command to run the application
CMD ["./app"]