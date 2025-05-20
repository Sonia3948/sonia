
#!/bin/bash

# Create go.mod if it doesn't exist
if [ ! -f go.mod ]; then
    echo "Initializing Go module..."
    go mod init elsaidaliya
fi

echo "Installing dependencies..."
go get -u github.com/gin-gonic/gin
go get -u github.com/gin-contrib/cors
go get -u go.mongodb.org/mongo-driver/mongo
go get -u github.com/joho/godotenv
go get -u github.com/google/uuid
go get -u golang.org/x/crypto/bcrypt

echo "Updating and tidying dependencies..."
go mod tidy

echo "Verifying dependencies..."
go mod verify

echo "Installation completed."
echo "You can now start the server with: go run main.go"
