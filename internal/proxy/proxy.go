package proxy

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/yourusername/webpanel/pkg/caddy"
)

const (
	siteConfigDir = "/etc/caddy/sites.d"
)

// Add membuat situs proxy baru dengan domain dan target yang ditentukan
func Add(domain, target string) {
	// Validasi domain dan target
	if !isValidDomain(domain) {
		fmt.Printf("Error: Domain tidak valid: %s\n", domain)
		return
	}

	if !isValidTarget(target) {
		fmt.Printf("Error: Target tidak valid: %s\n", target)
		return
	}

	// Buat file konfigurasi Caddy
	configContent := fmt.Sprintf(`%s {
	reverse_proxy %s
}
`, domain, target)

	configPath := filepath.Join(siteConfigDir, "proxy."+domain+".conf")
	if err := ioutil.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		fmt.Printf("Error: Tidak dapat menulis file konfigurasi: %s\n", err)
		return
	}

	// Muat ulang Caddy
	if err := caddy.Reload(); err != nil {
		fmt.Printf("Peringatan: Tidak dapat memuat ulang Caddy: %s\n", err)
	}

	fmt.Printf("Situs proxy %s -> %s berhasil dibuat\n", domain, target)
}

// Remove menghapus situs proxy dengan domain yang ditentukan
func Remove(domain string) {
	// Validasi domain
	if !isValidDomain(domain) {
		fmt.Printf("Error: Domain tidak valid: %s\n", domain)
		return
	}

	// Konfirmasi penghapusan
	fmt.Printf("Anda yakin ingin menghapus situs proxy %s? (y/N): ", domain)
	var response string
	fmt.Scanln(&response)
	if strings.ToLower(response) != "y" {
		fmt.Println("Penghapusan dibatalkan")
		return
	}

	// Hapus file konfigurasi
	configPath := filepath.Join(siteConfigDir, "proxy."+domain+".conf")
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		fmt.Printf("Error: Tidak dapat menghapus file konfigurasi: %s\n", err)
		return
	}

	// Muat ulang Caddy
	if err := caddy.Reload(); err != nil {
		fmt.Printf("Peringatan: Tidak dapat memuat ulang Caddy: %s\n", err)
	}

	fmt.Printf("Situs proxy %s berhasil dihapus\n", domain)
}

// List menampilkan semua situs proxy yang dikonfigurasi
func List() {
	files, err := ioutil.ReadDir(siteConfigDir)
	if err != nil {
		fmt.Printf("Error: Tidak dapat membaca direktori konfigurasi: %s\n", err)
		return
	}

	proxyCount := 0
	fmt.Println("Situs proxy yang dikonfigurasi:")
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "proxy.") && strings.HasSuffix(file.Name(), ".conf") {
			domain := strings.TrimPrefix(file.Name(), "proxy.")
			domain = strings.TrimSuffix(domain, ".conf")

			// Baca file untuk mendapatkan target
			content, err := ioutil.ReadFile(filepath.Join(siteConfigDir, file.Name()))
			if err != nil {
				fmt.Printf("- %s -> [Error membaca target]\n", domain)
				continue
			}

			// Ekstrak target dari konfigurasi
			target := extractTarget(string(content))
			fmt.Printf("- %s -> %s\n", domain, target)
			proxyCount++
		}
	}

	if proxyCount == 0 {
		fmt.Println("Tidak ada situs proxy yang dikonfigurasi")
	}
}

// isValidDomain memeriksa apakah domain valid
func isValidDomain(domain string) bool {
	// Implementasi sederhana, bisa ditingkatkan dengan validasi regex yang lebih baik
	return len(domain) > 0 && !strings.Contains(domain, " ") && strings.Contains(domain, ".")
}

// isValidTarget memeriksa apakah target valid
func isValidTarget(target string) bool {
	// Implementasi sederhana, bisa ditingkatkan dengan validasi yang lebih baik
	return len(target) > 0 && !strings.Contains(target, " ")
}

// extractTarget mengekstrak target dari konfigurasi
func extractTarget(config string) string {
	lines := strings.Split(config, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "reverse_proxy ") {
			return strings.TrimPrefix(line, "reverse_proxy ")
		}
	}
	return "[Target tidak ditemukan]"
}
