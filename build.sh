#!/bin/bash

PATH_AMD64="./bin/amd64"
PATH_I386="./bin/i386"
PATH_AMD64_STATIC="./bin/amd64st"
INSTALL_FILE="install.sh"
CONF_FILE="conf.yaml"

build_amd64() {
  GOOS=linux go build -o "${PATH_AMD64}/srvmon" ./cmd/*.go
  cp "$INSTALL_FILE" "${PATH_AMD64}/${INSTALL_FILE}"
  cp "$CONF_FILE" "${PATH_AMD64}/${CONF_FILE}"
  echo -n "Binary build done...!"
}

build_i386() {
  GOOS=linux GOARCH=386 go build -o "${PATH_I386}/srvmon" ./cmd/*.go
  cp "$INSTALL_FILE" "${PATH_I386}/${INSTALL_FILE}"
  cp "$CONF_FILE" "${PATH_I386}/${CONF_FILE}"
  echo -n "Binary build done...!"
}

build_amd64_static() {
  CGO_ENABLED=0 GOOS=linux go build -o "${PATH_AMD64_STATIC}/srvmon" ./cmd/*.go
  cp "$INSTALL_FILE" "${PATH_AMD64_STATIC}/${INSTALL_FILE}"
  cp "$CONF_FILE" "${PATH_AMD64_STATIC}/${CONF_FILE}"
  echo -n "Binary build done...!"
}

echo "Bulding srvmon...!"

echo "Select target architecture: [1|2|3]"
echo "1: Linux (AMD64)"
echo "2: Linux (i386)"
echo "3: Linux (AMD64) Statically Linked GLIBC"

read -p "Build for:" sel

case "$sel" in
  1)
    build_amd64
    ;;
  2)
    build_i386
    ;;
  3)
    build_amd64_static
    ;;
  *)
    echo "not a correct selection, aborting...!"
    ;;
esac
