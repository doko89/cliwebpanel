package caddy

import (
	"fmt"
	"os/exec"
)

// Reload triggers a Caddy configuration reload
func Reload() error {
	cmd := exec.Command("systemctl", "reload", "caddy")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error reloading Caddy: %v - %s", err, string(output))
	}
	return nil
}

// Restart restarts the Caddy service
func Restart() error {
	cmd := exec.Command("systemctl", "restart", "caddy")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error restarting Caddy: %v - %s", err, string(output))
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

// Status returns the status of the Caddy service
func Status() (string, error) {
	cmd := exec.Command("systemctl", "status", "caddy")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Don't return an error here as systemctl status returns non-zero
		// exit codes for stopped/failed services which is expected behavior
		return string(output), nil
	}
	return string(output), nil
}

// ValidateConfig validates the Caddy configuration without applying it
func ValidateConfig() error {
	cmd := exec.Command("caddy", "validate", "--config", "/etc/caddy/Caddyfile")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("invalid Caddy configuration: %v - %s", err, string(output))
	}
	return nil
}
