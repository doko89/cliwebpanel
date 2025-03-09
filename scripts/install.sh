#!/bin/bash

# Script instalasi untuk webpanel

# Periksa apakah berjalan sebagai root
if [ "$(id -u)" -ne 0 ]; then
    echo "Error: Script harus dijalankan sebagai root"
    exit 1
fi

# Fungsi untuk mendeteksi OS
detect_os() {
    if [ -f /etc/debian_version ]; then
        if grep -q "Ubuntu" /etc/issue; then
            echo "ubuntu"
        else
            echo "debian"
        fi
    else
        echo "unsupported"
    fi
}

# Dapatkan OS
OS=$(detect_os)

if [ "$OS" = "unsupported" ]; then
    echo "Error: Sistem operasi tidak didukung. Webpanel hanya mendukung Debian dan Ubuntu."
    exit 1
fi

echo "Menginstal webpanel pada $OS..."

# Instal dependensi dasar
apt-get update
apt-get install -y curl wget gnupg2 ca-certificates lsb-release apt-transport-https

# Unduh binary webpanel dari GitHub
ARCH=$(dpkg --print-architecture)
VERSION="0.1.0"  # Ganti dengan versi terbaru

case $ARCH in
    amd64)
        BINARY_URL="https://github.com/yourusername/webpanel/releases/download/v${VERSION}/webpanel_${VERSION}_linux_amd64.tar.gz"
        ;;
    i386)
        BINARY_URL="https://github.com/yourusername/webpanel/releases/download/v${VERSION}/webpanel_${VERSION}_linux_i386.tar.gz"
        ;;
    arm64)
        BINARY_URL="https://github.com/yourusername/webpanel/releases/download/v${VERSION}/webpanel_${VERSION}_linux_arm64.tar.gz"
        ;;
    armhf)
        BINARY_URL="https://github.com/yourusername/webpanel/releases/download/v${VERSION}/webpanel_${VERSION}_linux_armv7.tar.gz"
        ;;
    *)
        echo "Error: Arsitektur $ARCH tidak didukung"
        exit 1
        ;;
esac

echo "Mengunduh webpanel untuk $ARCH..."
wget -O /tmp/webpanel.tar.gz $BINARY_URL

if [ $? -ne 0 ]; then
    echo "Error: Tidak dapat mengunduh webpanel"
    exit 1
fi

# Ekstrak binary
mkdir -p /tmp/webpanel
tar -xzf /tmp/webpanel.tar.gz -C /tmp/webpanel

# Pindahkan binary ke /usr/local/bin
mv /tmp/webpanel/webpanel /usr/local/bin/
chmod +x /usr/local/bin/webpanel

# Bersihkan
rm -rf /tmp/webpanel /tmp/webpanel.tar.gz

echo "Webpanel berhasil diinstal!"
echo "Jalankan 'webpanel install' untuk menginstal dependensi yang diperlukan" 