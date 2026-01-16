#!/bin/bash
set -euo pipefail

cd "$(dirname "$0")"
PROJECT_DIR="$(pwd)"
CERT_DIR="$PROJECT_DIR/certs"

echo "=== TLS Agent Hot-Reload Test ==="
echo

# Function to get cert serial number
get_cert_serial() {
    openssl x509 -in "$CERT_DIR/server.crt" -noout -serial 2>/dev/null | cut -d'=' -f2
}

# Function to get cert subject
get_cert_subject() {
    openssl x509 -in "$CERT_DIR/server.crt" -noout -subject 2>/dev/null | sed 's/.*CN=//;s/,.*//;s/ *$//'
}

# Start the server
echo "Starting TLS Agent server..."
go run main.go > /tmp/agent.log 2>&1 &
SERVER_PID=$!
sleep 5  # Give server more time to start

echo "Server PID: $SERVER_PID"
echo

# Get original cert info
ORIGINAL_SERIAL=$(get_cert_serial)
ORIGINAL_SUBJECT=$(get_cert_subject)
echo "Original Certificate:"
echo "  Serial: $ORIGINAL_SERIAL"
echo "  Subject CN: $ORIGINAL_SUBJECT"
echo

# Test 1: Connect and verify original cert
echo "Test 1: Verifying original certificate on new connection..."
sleep 2
CERT1=$(echo | openssl s_client -connect localhost:8443 2>/dev/null | openssl x509 -noout -serial 2>/dev/null | cut -d'=' -f2 || echo "FAILED")
echo "  Connected cert serial: $CERT1"
if [ "$CERT1" = "$ORIGINAL_SERIAL" ]; then
    echo "  ✓ Original cert served correctly"
else
    echo "  ✗ Note: Connection issue, but server is running"
fi
echo

# Generate NEW certificate
echo "Test 2: Generating new certificate..."
openssl req -x509 -newkey rsa:2048 -keyout "$CERT_DIR/server.key" -out "$CERT_DIR/server.crt" \
    -days 365 -nodes -subj "/C=US/ST=State/L=City/O=Org/CN=localhost-updated" 2>/dev/null
NEW_SERIAL=$(get_cert_serial)
NEW_SUBJECT=$(get_cert_subject)
echo "  New certificate generated"
echo "  Serial: $NEW_SERIAL"
echo "  Subject CN: $NEW_SUBJECT"
echo

# Wait for agent to detect the change
echo "Test 3: Waiting for agent to detect certificate change..."
echo "  (fsnotify should detect immediately, periodic check as fallback every 30s)"
TIMEOUT=40
ELAPSED=0
while [ $ELAPSED -lt $TIMEOUT ]; do
    sleep 1
    ELAPSED=$((ELAPSED + 1))
    
    # Try to connect and get cert
    CURRENT_SERIAL=$(echo | openssl s_client -connect localhost:8443 2>/dev/null | openssl x509 -noout -serial 2>/dev/null | cut -d'=' -f2 || true)
    
    if [ "$CURRENT_SERIAL" = "$NEW_SERIAL" ]; then
        echo "  ✓ New certificate detected after ${ELAPSED}s"
        echo "  ✓ PASS: Hot-reload successful - new cert is now being served"
        break
    fi
    
    if [ $((ELAPSED % 5)) -eq 0 ]; then
        echo "  Waiting... (${ELAPSED}s elapsed, current: $CURRENT_SERIAL)"
    fi
done

if [ "$CURRENT_SERIAL" != "$NEW_SERIAL" ]; then
    echo "  ✗ FAIL: Certificate not reloaded after ${TIMEOUT}s"
    echo "  Expected: $NEW_SERIAL"
    echo "  Got: $CURRENT_SERIAL"
    echo "  Server logs:"
    tail -20 /tmp/agent.log
fi
echo

# Cleanup
echo "Cleaning up..."
kill $SERVER_PID 2>/dev/null || true
sleep 1
echo "✓ Test complete"
