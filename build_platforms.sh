echo "Building binaries for Linux..."
GOOS=linux GOARCH=amd64 go build -o build/password-manager_linux-amd64
GOOS=linux GOARCH=arm64 go build -o build/password-manager_linux-arm64

echo "Building binaries for MacOS..."
GOOS=darwin GOARCH=amd64 go build -o build/password-manager_darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o build/password-manager_darwin-arm64

echo "Building binaries for Windows..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC="x86_64-w64-mingw32-gcc" go build -o build/password-manager_windows-amd64.exe

echo "Done"