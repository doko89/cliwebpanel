package php

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	moduleDir = "/etc/caddy/module.d"
)

// List displays all available PHP versions
func List() {
	fmt.Println("Listing available PHP versions:")
	// TODO: Implementation
}

// ListInstalled displays all installed PHP versions
func ListInstalled() {
	fmt.Println("Listing installed PHP versions:")
	// TODO: Implementation
}

// Install installs a specific PHP version
func Install(version string) {
	fmt.Printf("Installing PHP version: %s\n", version)
	// TODO: Implementation
}

// Uninstall uninstalls a specific PHP version
func Uninstall(version string) {
	fmt.Printf("Uninstalling PHP version: %s\n", version)
	// TODO: Implementation
}

// ListModules displays available modules for a PHP version
func ListModules(version string) {
	fmt.Printf("Listing modules for PHP version: %s\n", version)
	// TODO: Implementation
}

// InstallModule installs a specific PHP module
func InstallModule(module string) {
	fmt.Printf("Installing PHP module: %s\n", module)
	// TODO: Implementation
}

// createPhpModule membuat modul Caddy untuk PHP
func createPhpModule(version string) {
	moduleContent := fmt.Sprintf(`(php%s) {
	php_fastcgi unix//run/php/php%s-fpm.sock
}
`, version, version)

	modulePath := filepath.Join(moduleDir, fmt.Sprintf("php%s", version))
	if err := ioutil.WriteFile(modulePath, []byte(moduleContent), 0644); err != nil {
		fmt.Printf("Peringatan: Tidak dapat membuat modul PHP untuk Caddy: %s\n", err)
	}
}

// removePhpModule menghapus modul Caddy untuk PHP
func removePhpModule(version string) {
	modulePath := filepath.Join(moduleDir, fmt.Sprintf("php%s", version))
	if err := os.Remove(modulePath); err != nil && !os.IsNotExist(err) {
		fmt.Printf("Peringatan: Tidak dapat menghapus modul PHP untuk Caddy: %s\n", err)
	}
}

// parsePhpVersions mengurai versi PHP dari output apt list
func parsePhpVersions(output string) []string {
	versions := []string{}
	re := regexp.MustCompile(`php(\d+\.\d+)-fpm`)

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			version := matches[1]
			if !contains(versions, version) {
				versions = append(versions, version)
			}
		}
	}

	return versions
}

// parseInstalledPhpVersions mengurai versi PHP terinstal dari output dpkg -l
func parseInstalledPhpVersions(output string) []string {
	versions := []string{}
	re := regexp.MustCompile(`php(\d+\.\d+)-fpm`)

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ii ") {
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				version := matches[1]
				if !contains(versions, version) {
					versions = append(versions, version)
				}
			}
		}
	}

	return versions
}

// parsePhpModules mengurai modul PHP dari output apt list
func parsePhpModules(output string, version string) []string {
	modules := []string{}
	re := regexp.MustCompile(fmt.Sprintf(`php%s-([a-z0-9]+)`, version))

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			module := matches[1]
			if module != "fpm" && module != "common" && module != "cli" {
				modules = append(modules, module)
			}
		}
	}

	return modules
}

// isPhpInstalled memeriksa apakah versi PHP terinstal
func isPhpInstalled(version string) bool {
	cmd := exec.Command("dpkg", "-l", fmt.Sprintf("php%s-fpm", version))
	err := cmd.Run()
	return err == nil
}

// isValidPhpVersion memeriksa apakah versi PHP valid
func isValidPhpVersion(version string) bool {
	re := regexp.MustCompile(`^\d+\.\d+$`)
	return re.MatchString(version)
}

// contains memeriksa apakah slice berisi nilai tertentu
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
