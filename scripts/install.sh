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

# Buat direktori yang diperlukan
mkdir -p /etc/caddy/sites.d
mkdir -p /etc/caddy/module.d
mkdir -p /apps/sites
mkdir -p /backup/daily
mkdir -p /backup/weekly

# Unduh binary webpanel dari GitHub
ARCH=$(dpkg --print-architecture)
VERSION="0.1.0"  # Ganti dengan versi terbaru

case $ARCH in
    amd64)
        BINARY_URL="https://github.com/doko89/webpanel/releases/download/v${VERSION}/webpanel_linux_amd64.tar.gz"
        ;;
    i386)
        BINARY_URL="https://github.com/doko89/webpanel/releases/download/v${VERSION}/webpanel_linux_i386.tar.gz"
        ;;
    arm64)
        BINARY_URL="https://github.com/doko89/webpanel/releases/download/v${VERSION}/webpanel_linux_arm64.tar.gz"
        ;;
    armhf)
        BINARY_URL="https://github.com/doko89/webpanel/releases/download/v${VERSION}/webpanel_linux_armv7.tar.gz"
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
mv /tmp/webpanel/webpanel_linux_* /usr/local/bin/webpanel
chmod +x /usr/local/bin/webpanel

# Buat Caddyfile dasar
cat > /etc/caddy/Caddyfile << EOF
{
    # Konfigurasi global
    admin off
    email admin@localhost
}

# Impor semua situs dari direktori sites.d
import sites.d/*.conf
EOF

# Buat modul default
mkdir -p /etc/caddy/module.d
cat > /etc/caddy/module.d/spa << EOF
(spa) {
    @spa {
        not path *.php
        not path /api/*
        not path *.js
        not path *.css
        not path *.png
        not path *.jpg
        not path *.jpeg
        not path *.svg
        not path *.gif
        not path *.ico
        not path *.woff
        not path *.woff2
        not path *.ttf
        not path *.eot
        file {
            try_files {path} /index.html
        }
    }
    rewrite @spa /index.html
}
EOF

cat > /etc/caddy/module.d/security << EOF
(security) {
    header {
        # Keamanan dasar
        X-XSS-Protection "1; mode=block"
        X-Content-Type-Options "nosniff"
        X-Frame-Options "SAMEORIGIN"
        Referrer-Policy "strict-origin-when-cross-origin"
        
        # Hapus header yang tidak perlu
        -Server
        -X-Powered-By
    }
}
EOF

cat > /etc/caddy/module.d/ratelimit << EOF
(ratelimit) {
    rate_limit {
        zone dynamic {
            key {remote_host}
            events 100
            window 10s
        }
    }
}
EOF

cat > /etc/caddy/module.d/compression << EOF
(compression) {
    encode gzip zstd
}
EOF

cat > /etc/caddy/module.d/cache-headers << EOF
(cache-headers) {
    @static {
        path *.css *.js *.png *.jpg *.jpeg *.gif *.ico *.svg *.woff *.woff2 *.ttf *.eot
    }
    header @static Cache-Control "public, max-age=31536000"
}
EOF

cat > /etc/caddy/module.d/local-access << EOF
(local-access) {
    @local {
        remote_ip 127.0.0.1 192.168.0.0/16 10.0.0.0/8 172.16.0.0/12
    }
    @notLocal {
        not remote_ip 127.0.0.1 192.168.0.0/16 10.0.0.0/8 172.16.0.0/12
    }
    handle @notLocal {
        respond "Akses ditolak" 403
    }
}
EOF

# Bersihkan
rm -rf /tmp/webpanel /tmp/webpanel.tar.gz

# Tambahkan pengguna caddy jika belum ada
if ! id -u caddy > /dev/null 2>&1; then
    useradd -r -d /var/lib/caddy -s /usr/sbin/nologin caddy
fi

# Atur kepemilikan direktori
chown -R caddy:caddy /etc/caddy
chown -R caddy:caddy /apps/sites
chown -R caddy:caddy /backup

echo "Webpanel berhasil diinstal!"
echo "Jalankan 'webpanel install' untuk menginstal dependensi yang diperlukan"
echo "  - Caddy web server"
echo "  - PHP (versi yang Anda pilih)"
echo "  - Composer (opsional)"
echo "  - Node.js (opsional)"