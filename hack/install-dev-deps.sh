#!/bin/sh -eu

install_protoc() {
    version="${1}"
    arch="${2}"

    PROTOC_ZIP="protoc-${version}-${arch}.zip"

    curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v$version/$PROTOC_ZIP
    sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
    rm -f $PROTOC_ZIP
}

install_golangci_lint() {
  version="${1}"
  curl --silent --fail --location \
    "https://install.goreleaser.com/github.com/golangci/golangci-lint.sh" \
    | sh -s -- -b $GOPATH/bin "${version}"
  unset version
}

install_protoc "3.6.1" "osx-x86_64"
install_golangci_lint "v1.15.0"