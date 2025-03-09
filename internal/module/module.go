package module

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/doko89/webpanel/pkg/caddy"
)

const (
	moduleDir     = "/etc/caddy/module.d"
	siteConfigDir = "/etc/caddy/sites.d"
)

// Enable enables a module for a specific domain
func Enable(module, domain string) {
	fmt.Printf("Enabling module %s for domain: %s\n", module, domain)
	// Validasi modul
	if !isModuleAvailable(module) {
		fmt.Printf("Error: Modul tidak tersedia: %s\n", module)
		return
	}

	// Validasi domain
	configPath := filepath.Join(siteConfigDir, domain+".conf")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("Error: Domain tidak ditemukan: %s\n", domain)
		return
	}

	// Baca konfigurasi situs
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error: Tidak dapat membaca file konfigurasi: %s\n", err)
		return
	}

	// Periksa apakah modul sudah diaktifkan
	if strings.Contains(string(content), "import "+module) {
		fmt.Printf("Modul %s sudah diaktifkan untuk %s\n", module, domain)
		return
	}

	// Tambahkan modul ke konfigurasi
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.Contains(line, domain+" {") {
			// Tambahkan modul setelah baris pembuka
			lines[i+1] = lines[i+1] + "\n\timport " + module
			break
		}
	}

	// Tulis kembali konfigurasi
	newContent := strings.Join(lines, "\n")
	if err := ioutil.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		fmt.Printf("Error: Tidak dapat menulis file konfigurasi: %s\n", err)
		return
	}

	// Muat ulang Caddy
	if err := caddy.Reload(); err != nil {
		fmt.Printf("Peringatan: Tidak dapat memuat ulang Caddy: %s\n", err)
	}

	fmt.Printf("Modul %s berhasil diaktifkan untuk %s\n", module, domain)
}

// Disable disables a module for a specific domain
func Disable(module, domain string) {
	fmt.Printf("Disabling module %s for domain: %s\n", module, domain)
	// Validasi domain
	configPath := filepath.Join(siteConfigDir, domain+".conf")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("Error: Domain tidak ditemukan: %s\n", domain)
		return
	}

	// Baca konfigurasi situs
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error: Tidak dapat membaca file konfigurasi: %s\n", err)
		return
	}

	// Periksa apakah modul diaktifkan
	if !strings.Contains(string(content), "import "+module) {
		fmt.Printf("Modul %s tidak diaktifkan untuk %s\n", module, domain)
		return
	}

	// Hapus modul dari konfigurasi
	newContent := strings.Replace(string(content), "\timport "+module, "", -1)
	newContent = strings.Replace(newContent, "\n\n", "\n", -1) // Bersihkan baris kosong ganda

	// Tulis kembali konfigurasi
	if err := ioutil.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		fmt.Printf("Error: Tidak dapat menulis file konfigurasi: %s\n", err)
		return
	}

	// Muat ulang Caddy
	if err := caddy.Reload(); err != nil {
		fmt.Printf("Peringatan: Tidak dapat memuat ulang Caddy: %s\n", err)
	}

	fmt.Printf("Modul %s berhasil dinonaktifkan untuk %s\n", module, domain)
}

// List displays all modules enabled for a domain
func List(domain string) {
	fmt.Printf("Listing modules for domain: %s\n", domain)
	// Validasi domain
	configPath := filepath.Join(siteConfigDir, domain+".conf")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("Error: Domain tidak ditemukan: %s\n", domain)
		return
	}

	// Baca konfigurasi situs
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error: Tidak dapat membaca file konfigurasi: %s\n", err)
		return
	}

	// Cari modul yang diaktifkan
	modules := []string{}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "import ") {
			moduleName := strings.TrimPrefix(line, "import ")
			modules = append(modules, moduleName)
		}
	}

	// Tampilkan modul
	if len(modules) == 0 {
		fmt.Printf("Tidak ada modul yang diaktifkan untuk %s\n", domain)
		return
	}

	fmt.Printf("Modul yang diaktifkan untuk %s:\n", domain)
	for _, module := range modules {
		fmt.Println("-", module)
	}
}

// ListAvailable displays all available modules
func ListAvailable() {
	fmt.Println("Listing all available modules:")
	files, err := ioutil.ReadDir(moduleDir)
	if err != nil {
		fmt.Printf("Error: Tidak dapat membaca direktori modul: %s\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("Tidak ada modul yang tersedia")
		return
	}

	fmt.Println("Modul yang tersedia:")
	for _, file := range files {
		if !file.IsDir() {
			moduleName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			fmt.Println("-", moduleName)
		}
	}
}

// isModuleAvailable memeriksa apakah modul tersedia
func isModuleAvailable(moduleName string) bool {
	modulePath := filepath.Join(moduleDir, moduleName)
	_, err := os.Stat(modulePath)
	return !os.IsNotExist(err)
}
