#!/usr/bin/env bash
set -e

REPO="dinesh00509/gitease"
VERSION=${1:-latest}

echo " Installing GitEase..."

if [ "$VERSION" = "latest" ]; then
  VERSION=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep -Po '"tag_name": "\K.*?(?=")')
fi

OS=$(uname -s)
ARCH=$(uname -m)

case $OS in
  Linux)
    FILE="gitease_Linux_x86_64.tar.gz"
    ;;
  Darwin)
    FILE="gitease_Darwin_x86_64.tar.gz"
    ;;
  *)
    echo " Unsupported OS: $OS"
    exit 1
    ;;
esac

URL="https://github.com/$REPO/releases/download/$VERSION/$FILE"

echo " Downloading $URL ..."
curl -L -o /tmp/gitease.tar.gz "$URL"

echo "Extracting..."
tar -xzf /tmp/gitease.tar.gz -C /tmp

echo " Installing..."
sudo mv /tmp/gitease /usr/local/bin/gitease

echo " Installation complete!"
gitease --version

