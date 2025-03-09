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

```bash
# Download the latest release for your platform
# Make executable and move to path
chmod +x webpanel_linux_amd64
sudo mv webpanel_linux_amd64 /usr/local/bin/webpanel
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
make build
```

## License

MIT
