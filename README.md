# WebPanel CLI

A command-line web panel tool for server administration with minimal sysadmin experience required.

## Features

- Site and proxy management
- Modular Caddy configuration system
- Backup system for sites and databases
- Database management
- PHP version management
- Software installation helpers

## Installation

### Quick Install (Debian/Ubuntu)

```bash
curl -fsSL https://raw.githubusercontent.com/doko89/webpanel/main/scripts/install.sh | sudo bash
```

### Manual Installation

```bash
# Download the latest release for your platform
wget https://github.com/doko89/webpanel/releases/download/v0.1.0/webpanel_linux_amd64.tar.gz
tar -xzf webpanel_linux_amd64.tar.gz
chmod +x webpanel_linux_amd64
sudo mv webpanel_linux_amd64 /usr/local/bin/webpanel

# Install dependencies
sudo webpanel install
```

## Usage

```bash
# Get help
webpanel help

# Site management
webpanel site add domain.com
webpanel site remove domain.com
webpanel site list

# PHP management
webpanel php list
webpanel php install 8.1

# Module management
webpanel module enable php81 domain.com
```

## Building from source

```bash
git clone https://github.com/doko89/webpanel.git
cd webpanel
make build
```

## License

MIT
