#!/bin/bash

URL="http://localhost:8080"
MAX_ATTEMPTS=10
WAIT_TIME=1

attempt=1
while [ $attempt -le $MAX_ATTEMPTS ]; do
    echo "Attempt $attempt/$MAX_ATTEMPTS..."

    if curl --output /dev/null --silent --head --fail "$URL"; then
        echo "Server is up!"
        bunx vite
        exit 0
    else
        echo "Attempt $attempt failed. Retrying in $WAIT_TIME seconds..."
        sleep $WAIT_TIME
    fi

    attempt=$((attempt + 1))
done

echo "Failed to reach the server after $MAX_ATTEMPTS attempts."
exit 1
