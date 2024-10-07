# Use the official Golang image
FROM golang:1.23

# Install dependencies for Chrome and other tools
RUN apt-get update && apt-get install -y \
    wget \
    curl \
    gnupg \
    ca-certificates \
    unzip \
    libnss3 \
    libxss1 \
    libasound2 \
    fonts-liberation \
    libappindicator3-1 \
    libatk-bridge2.0-0 \
    libgtk-3-0 \
    xdg-utils \
    --no-install-recommends

# Add Google Chrome's official GPG key and repo
RUN wget -q -O - https://dl.google.com/linux/linux_signing_key.pub | apt-key add - && \
    sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" > /etc/apt/sources.list.d/google-chrome.list' && \
    apt-get update && apt-get install -y google-chrome-stable

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Install Go dependencies
RUN go mod download

# Copy the source code
COPY cmd/scraper/ .

# Build the Go binary
RUN go build -o app .

# Set display port for headless Chrome
ENV DISPLAY=:99

# Command to run the binary
CMD ["./app"]
