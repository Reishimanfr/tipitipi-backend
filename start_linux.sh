#!/bin/bash

DIST_PATH=./src/backend/dist
URL=http://localhost:2333/heartbeat
MAX_ATTEMPTS=10
WAIT_TIME=1
attempt=1

# Create dist folder if it doesn't exist
if [ ! -d "$DIST_PATH" ]; then
    mkdir -p "$DIST_PATH"
fi

# Build server binary if not built
if [ ! -f "$DIST_PATH/server" ]; then
    # Cd into backend folder to build code
    cd ./src/backend || exit

    # Build server binary
    go build -o "./dist/server"
    cd ../../ || exit
fi

# Start server
"$DIST_PATH/server" &

# Wait for server to respond or exit
while [ $attempt -le $MAX_ATTEMPTS ]; do
    echo "Attempt $attempt/$MAX_ATTEMPTS..."

    if curl --output /dev/null --silent --head --fail "$URL"; then
        echo "Server is up!"
        # Start up vite if server responded
        bunx vite
        exit 0
    else
        echo "Attempt $attempt failed. Retrying in $WAIT_TIME seconds..."
        sleep $WAIT_TIME
    fi

    attempt=$((attempt + 1))
done

echo "Failed to reach the server after $MAX_ATTEMPTS attempts."
pkill server
exit 1
