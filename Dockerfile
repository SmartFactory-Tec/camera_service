# Specifies a parent image
FROM golang:1.20.2-bullseye

LABEL org.opencontainers.image.source=https://github.com/Smartfactory-Tec/camera_service

RUN curl -fsSL \
        https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
        sh

# Creates an app directory to hold your app’s source code
WORKDIR /camera_service

COPY go.mod .
COPY go.sum .

# Installs Go dependencies
RUN go mod download

# Copies everything from your root directory into /app
COPY . .

# Builds your app with optional configuration
RUN go build -buildvcs=false -o ./camera_service github.com/SmartFactory-Tec/camera_service/cmd/camera_service

ENV CAMERA_SERVICE_CONFIG=/config

# Tells Docker which network port your container listens on
EXPOSE 3000

# Specifies the executable command that runs when the container starts
ENTRYPOINT [ "/camera_service/camera_service" ]