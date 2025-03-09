package caddy

import (
	"fmt"
	"os/exec"
)

// Reload memuat ulang konfigurasi Caddy
func Reload() error {
	cmd := exec.Command("systemctl", "reload", "caddy")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tidak dapat memuat ulang Caddy: %w", err)
	}
	return nil
}

// Start memulai layanan Caddy
func Start() error {
	cmd := exec.Command("systemctl", "start", "caddy")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tidak dapat memulai Caddy: %w", err)
	}
	return nil
}

// Stop menghentikan layanan Caddy
func Stop() error {
	cmd := exec.Command("systemctl", "stop", "caddy")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tidak dapat menghentikan Caddy: %w", err)
	}
	return nil
}

// Status mendapatkan status layanan Caddy
func Status() (string, error) {
	cmd := exec.Command("systemctl", "status", "caddy")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("tidak dapat mendapatkan status Caddy: %w", err)
	}
	return string(output), nil
}

// Validate memvalidasi konfigurasi Caddy
func Validate() error {
	cmd := exec.Command("caddy", "validate", "--config", "/etc/caddy/Caddyfile")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("konfigurasi Caddy tidak valid: %w", err)
	}
	return nil
}
