package php

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/doko89/webpanel/pkg/caddy"
)

const (
	moduleDir = "/etc/caddy/module.d"
)

// List menampilkan semua versi PHP yang tersedia
func List() {
	// Implementasi tergantung pada OS
	// Contoh untuk Debian/Ubuntu
	cmd := exec.Command("apt", "list", "php*-fpm")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: Tidak dapat mendapatkan daftar versi PHP: %s\n", err)
		return
	}

	// Parse output untuk mendapatkan versi PHP
	versions := parsePhpVersions(string(output))

	if len(versions) == 0 {
		fmt.Println("Tidak ada versi PHP yang tersedia")
		return
	}

	fmt.Println("Versi PHP yang tersedia:")
	for _, version := range versions {
		installed := isPhpInstalled(version)
		if installed {
			fmt.Printf("- php%s (terinstal)\n", version)
		} else {
			fmt.Printf("- php%s\n", version)
		}
	}
}

// ListInstalled menampilkan versi PHP yang terinstal
func ListInstalled() {
	// Implementasi tergantung pada OS
	// Contoh untuk Debian/Ubuntu
	cmd := exec.Command("dpkg", "-l", "php*-fpm")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: Tidak dapat mendapatkan daftar versi PHP terinstal: %s\n", err)
		return
	}

	// Parse output untuk mendapatkan versi PHP terinstal
	versions := parseInstalledPhpVersions(string(output))

	if len(versions) == 0 {
		fmt.Println("Tidak ada versi PHP yang terinstal")
		return
	}

	fmt.Println("Versi PHP yang terinstal:")
	for _, version := range versions {
		fmt.Printf("- php%s\n", version)
	}
}

// Install menginstal versi PHP tertentu
func Install(version string) {
	// Validasi versi
	if !isValidPhpVersion(version) {
		fmt.Printf("Error: Versi PHP tidak valid: %s\n", version)
		return
	}

	// Periksa apakah sudah terinstal
	if isPhpInstalled(version) {
		fmt.Printf("PHP %s sudah terinstal\n", version)
		return
	}

	// Instal PHP
	fmt.Printf("Menginstal PHP %s...\n", version)
	cmd := exec.Command("apt", "install", "-y",
		fmt.Sprintf("php%s-fpm", version),
		fmt.Sprintf("php%s-common", version),
		fmt.Sprintf("php%s-cli", version),
		fmt.Sprintf("php%s-mysql", version),
		fmt.Sprintf("php%s-curl", version),
		fmt.Sprintf("php%s-gd", version),
		fmt.Sprintf("php%s-mbstring", version),
		fmt.Sprintf("php%s-xml", version),
		fmt.Sprintf("php%s-zip", version))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: Tidak dapat menginstal PHP %s: %s\n", version, err)
		return
	}

	// Buat modul Caddy untuk PHP
	createPhpModule(version)

	// Muat ulang Caddy
	if err := caddy.Reload(); err != nil {
		fmt.Printf("Peringatan: Tidak dapat memuat ulang Caddy: %s\n", err)
	}

	fmt.Printf("PHP %s berhasil diinstal\n", version)
}

// Uninstall menghapus instalasi versi PHP tertentu
func Uninstall(version string) {
	// Validasi versi
	if !isValidPhpVersion(version) {
		fmt.Printf("Error: Versi PHP tidak valid: %s\n", version)
		return
	}

	// Periksa apakah terinstal
	if !isPhpInstalled(version) {
		fmt.Printf("PHP %s tidak terinstal\n", version)
		return
	}

	// Konfirmasi penghapusan
	fmt.Printf("Anda yakin ingin menghapus PHP %s? (y/N): ", version)
	var response string
	fmt.Scanln(&response)
	if strings.ToLower(response) != "y" {
		fmt.Println("Penghapusan dibatalkan")
		return
	}

	// Hapus PHP
	fmt.Printf("Menghapus PHP %s...\n", version)
	cmd := exec.Command("apt", "remove", "-y", fmt.Sprintf("php%s*", version))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: Tidak dapat menghapus PHP %s: %s\n", version, err)
		return
	}

	// Hapus modul Caddy untuk PHP
	removePhpModule(version)

	// Muat ulang Caddy
	if err := caddy.Reload(); err != nil {
		fmt.Printf("Peringatan: Tidak dapat memuat ulang Caddy: %s\n", err)
	}

	fmt.Printf("PHP %s berhasil dihapus\n", version)
}

// ListModules menampilkan modul PHP yang tersedia
func ListModules(version string) {
	// Validasi versi
	if !isValidPhpVersion(version) {
		fmt.Printf("Error: Versi PHP tidak valid: %s\n", version)
		return
	}

	// Periksa apakah terinstal
	if !isPhpInstalled(version) {
		fmt.Printf("PHP %s tidak terinstal\n", version)
		return
	}

	// Dapatkan daftar modul
	cmd := exec.Command("apt", "list", fmt.Sprintf("php%s-*", version))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: Tidak dapat mendapatkan daftar modul PHP: %s\n", err)
		return
	}

	// Parse output
	modules := parsePhpModules(string(output), version)

	if len(modules) == 0 {
		fmt.Printf("Tidak ada modul PHP %s yang tersedia\n", version)
		return
	}

	fmt.Printf("Modul PHP %s yang tersedia:\n", version)
	for _, module := range modules {
		fmt.Println("-", module)
	}
}

// InstallModule menginstal modul PHP
func InstallModule(moduleSpec string) {
	// Validasi spesifikasi modul
	parts := strings.Split(moduleSpec, "-")
	if len(parts) != 2 {
		fmt.Printf("Error: Spesifikasi modul tidak valid: %s (gunakan format: 8.1-gd)\n", moduleSpec)
		return
	}

	version := parts[0]
	module := parts[1]

	// Validasi versi
	if !isValidPhpVersion(version) {
		fmt.Printf("Error: Versi PHP tidak valid: %s\n", version)
		return
	}

	// Periksa apakah PHP terinstal
	if !isPhpInstalled(version) {
		fmt.Printf("Error: PHP %s tidak terinstal\n", version)
		return
	}

	// Instal modul
	fmt.Printf("Menginstal modul PHP %s-%s...\n", version, module)
	cmd := exec.Command("apt", "install", "-y", fmt.Sprintf("php%s-%s", version, module))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: Tidak dapat menginstal modul: %s\n", err)
		return
	}

	fmt.Printf("Modul PHP %s-%s berhasil diinstal\n", version, module)
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
