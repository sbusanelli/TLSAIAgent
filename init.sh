
# Create certs directory if missing
mkdir -p certs

# Generate self-signed cert valid for 365 days
openssl req -x509 -newkey rsa:2048 -keyout certs/server.key -out certs/server.crt -days 365 -nodes \
  -subj "/C=US/ST=State/L=City/O=Org/CN=localhost"

# Verify
ls -la certs/
openssl x509 -in certs/server.crt -text -noout | head -20
