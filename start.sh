#!/bin/bash
# Start WhatsApp service script

echo "Starting WhatsApp Service..."
cd "$(dirname "$0")"

# Check if node_modules exists
if [ ! -d "node_modules" ]; then
    echo "Installing dependencies..."
    npm install
fi

# Start the service
echo "Starting service on port ${WHATSAPP_SERVICE_PORT:-3002}..."
npm start






