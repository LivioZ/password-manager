echo "Building binaries for Linux..."
GOOS=linux GOARCH=amd64 go build -o build/password-manager_linux-amd64
GOOS=linux GOARCH=arm64 go build -o build/password-manager_linux-arm64

echo "Building binaries for MacOS..."
GOOS=darwin GOARCH=amd64 go build -o build/password-manager_darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o build/password-manager_darwin-arm64

echo "Done"