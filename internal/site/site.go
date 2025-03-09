package site

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/yourusername/webpanel/pkg/caddy"
)

const (
	sitesDir      = "/apps/sites"
	siteConfigDir = "/etc/caddy/sites.d"
)

// Add membuat situs baru dengan domain yang ditentukan
func Add(domain string) {
	// Validasi domain
	if !isValidDomain(domain) {
		fmt.Printf("Error: Domain tidak valid: %s\n", domain)
		return
	}

	// Buat direktori situs
	siteDir := filepath.Join(sitesDir, domain)
	if err := os.MkdirAll(siteDir, 0755); err != nil {
		fmt.Printf("Error: Tidak dapat membuat direktori situs: %s\n", err)
		return
	}

	// Buat file konfigurasi Caddy
	configContent := fmt.Sprintf(`%s {
	root * %s
	file_server
}
`, domain, siteDir)

	configPath := filepath.Join(siteConfigDir, domain+".conf")
	if err := ioutil.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		fmt.Printf("Error: Tidak dapat menulis file konfigurasi: %s\n", err)
		return
	}

	// Atur kepemilikan direktori
	if err := os.Chown(siteDir, getCaddyUID(), getCaddyGID()); err != nil {
		fmt.Printf("Peringatan: Tidak dapat mengubah kepemilikan direktori: %s\n", err)
	}

	// Muat ulang Caddy
	if err := caddy.Reload(); err != nil {
		fmt.Printf("Peringatan: Tidak dapat memuat ulang Caddy: %s\n", err)
	}

	fmt.Printf("Situs %s berhasil dibuat\n", domain)
}

// Remove menghapus situs dengan domain yang ditentukan
func Remove(domain string) {
	// Validasi domain
	if !isValidDomain(domain) {
		fmt.Printf("Error: Domain tidak valid: %s\n", domain)
		return
	}

	// Konfirmasi penghapusan
	fmt.Printf("Anda yakin ingin menghapus situs %s? (y/N): ", domain)
	var response string
	fmt.Scanln(&response)
	if strings.ToLower(response) != "y" {
		fmt.Println("Penghapusan dibatalkan")
		return
	}

	// Hapus file konfigurasi
	configPath := filepath.Join(siteConfigDir, domain+".conf")
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		fmt.Printf("Error: Tidak dapat menghapus file konfigurasi: %s\n", err)
		return
	}

	// Muat ulang Caddy
	if err := caddy.Reload(); err != nil {
		fmt.Printf("Peringatan: Tidak dapat memuat ulang Caddy: %s\n", err)
	}

	fmt.Printf("Situs %s berhasil dihapus\n", domain)
	fmt.Printf("Catatan: Direktori situs di %s/%s tidak dihapus untuk keamanan data\n", sitesDir, domain)
}

// List menampilkan semua situs yang dikonfigurasi
func List() {
	files, err := ioutil.ReadDir(siteConfigDir)
	if err != nil {
		fmt.Printf("Error: Tidak dapat membaca direktori konfigurasi: %s\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("Tidak ada situs yang dikonfigurasi")
		return
	}

	fmt.Println("Situs yang dikonfigurasi:")
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".conf") {
			domain := strings.TrimSuffix(file.Name(), ".conf")
			fmt.Println("-", domain)
		}
	}
}

// isValidDomain memeriksa apakah domain valid
func isValidDomain(domain string) bool {
	// Implementasi sederhana, bisa ditingkatkan dengan validasi regex yang lebih baik
	return len(domain) > 0 && !strings.Contains(domain, " ") && strings.Contains(domain, ".")
}

// getCaddyUID mendapatkan UID pengguna caddy
func getCaddyUID() int {
	// Implementasi untuk mendapatkan UID pengguna caddy
	return 0 // Sementara mengembalikan 0 (root)
}

// getCaddyGID mendapatkan GID grup caddy
func getCaddyGID() int {
	// Implementasi untuk mendapatkan GID grup caddy
	return 0 // Sementara mengembalikan 0 (root)
}
