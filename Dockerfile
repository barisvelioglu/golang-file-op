FROM golang:1.17 as builder

ENV HOME /app
ENV CGO_ENABLED 0
ENV GOOS linux

RUN GOCACHE=OFF

RUN go env -w GOPRIVATE=gitlab.com/thetasphere-maestro
RUN git config --global url."https://golang:jhY1xrivPLNGJTbkJCtJ@gitlab.com".insteadOf "https://gitlab.com"

# Move to working directory /build
WORKDIR /app

# Copy and download dependency using go mod
COPY go.mod .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

FROM alpine:latest

#RUN apk --no-cache add ca-certificates

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Command to run when starting the container
CMD ["/dist/main"]
