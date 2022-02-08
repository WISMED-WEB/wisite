#!/bin/bash

set -e

######################
# Become a Certificate Authority
######################

OUT=./cert
rm -rf $OUT/ && mkdir $_

# Generate private key
openssl genrsa -des3 -out $OUT/myCA.key 2048
# Generate root certificate
openssl req -x509 -new -nodes -key $OUT/myCA.key -sha256 -days 825 -out $OUT/myCA.pem

######################
# Create CA-signed certs
######################

NAME=localhost # Use your own domain name
# Generate a private key
openssl genrsa -out $OUT/$NAME.key 2048
# Create a certificate-signing request
openssl req -new -key $OUT/$NAME.key -out $OUT/$NAME.csr
# Create a config file for the extensions
>$OUT/$NAME.ext cat <<-EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = $NAME # Be sure to include the domain name here because Common Name is not so commonly honoured by itself
DNS.2 = bar.$NAME # Optionally, add additional domains (I've added a subdomain here)
IP.1 = 192.168.0.13 # Optionally, add an IP address (if the connection which you have planned requires it)
EOF
# Create the signed certificate
openssl x509 -req -in $OUT/$NAME.csr -CA $OUT/myCA.pem -CAkey $OUT/myCA.key -CAcreateserial \
-out $OUT/$NAME.crt -days 3650 -sha256 -extfile $OUT/$NAME.ext