#!/bin/bash

set -e

# rm -rf cert && mkdir $_
# openssl genrsa -out cert/private.key 4096
# openssl req -new -x509 -sha256 -days 1825 -key cert/private.key -subj "/CN=192.168.31.227" -addext "subjectAltName = DNS:localhost" -out cert/public.crt

# rm -rf ../client/cert && mkdir $_
# cp ./cert/public.crt ../client/cert/

rm -rf cert && mkdir $_
openssl genrsa -out cert/private.pem 4096
openssl req -new -x509 -sha256 -days 1825 -key cert/private.pem -subj "/CN=192.168.31.227" -addext "subjectAltName = DNS:localhost" -out cert/public.pem
openssl pkcs12 -export -in cert/public.pem -inkey cert/private.pem -out cert/server.p12

rm -rf ../client/cmd/cert && mkdir $_
cp -r ./cert/ ../client/cmd/


# go run /usr/local/go/src/crypto/tls/generate_cert.go --host localhost
# openssl pkcs12 -export -in cert/cert.pem -inkey cert/key.pem -out cert/server.p12