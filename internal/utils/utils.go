package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// InstallDependencies menginstal semua dependensi yang diperlukan
func InstallDependencies() {
	fmt.Println("Menginstal dependensi yang diperlukan...")

	// Periksa OS
	osType := detectOS()

	// Instal Caddy
	installCaddy()

	// Instal PHP repository
	installPHPRepository(osType)

	// Tanya pengguna untuk menginstal PHP
	installPHPVersions()

	// Tanya pengguna untuk menginstal Composer
	installComposer()

	// Tanya pengguna untuk menginstal Node.js
	installNodeJS()

	// Setup direktori dan pengguna
	setupDirectories()

	// Setup modul default
	setupDefaultModules()

	fmt.Println("Instalasi dependensi selesai!")
}

// detectOS mendeteksi sistem operasi
func detectOS() string {
	// Periksa apakah ini Debian
	if _, err := os.Stat("/etc/debian_version"); !os.IsNotExist(err) {
		// Periksa apakah ini Ubuntu
		cmd := exec.Command("lsb_release", "-a")
		output, err := cmd.CombinedOutput()
		if err == nil && strings.Contains(string(output), "Ubuntu") {
			return "ubuntu"
		}
		return "debian"
	}

	// Default ke debian jika tidak dapat mendeteksi
	return "debian"
}

// installCaddy menginstal Caddy
func installCaddy() {
	fmt.Println("Menginstal Caddy...")

	// Tambahkan repositori Caddy
	cmd := exec.Command("sh", "-c", `
		apt-get update
		apt-get install -y debian-keyring debian-archive-keyring apt-transport-https
		curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
		curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | tee /etc/apt/sources.list.d/caddy-stable.list
		apt-get update
	`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Peringatan: Tidak dapat menambahkan repositori Caddy: %s\n", err)
	}

	// Instal Caddy
	cmd = exec.Command("apt-get", "install", "-y", "caddy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: Tidak dapat menginstal Caddy: %s\n", err)
		return
	}

	fmt.Println("Caddy berhasil diinstal")
}

// installPHPRepository menginstal repositori PHP
func installPHPRepository(osType string) {
	fmt.Println("Menginstal repositori PHP...")

	if osType == "ubuntu" {
		// Tambahkan PPA ondrej/php untuk Ubuntu
		cmd := exec.Command("sh", "-c", `
			apt-get update
			apt-get install -y software-properties-common
			add-apt-repository -y ppa:ondrej/php
			apt-get update
		`)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Peringatan: Tidak dapat menambahkan repositori PHP: %s\n", err)
		}
	} else {
		// Tambahkan repositori sury.org untuk Debian
		cmd := exec.Command("sh", "-c", `
			apt-get update
			apt-get install -y apt-transport-https lsb-release ca-certificates curl
			curl -sSLo /usr/share/keyrings/deb.sury.org-php.gpg https://packages.sury.org/php/apt.gpg
			sh -c 'echo "deb [signed-by=/usr/share/keyrings/deb.sury.org-php.gpg] https://packages.sury.org/php/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/php.list'
			apt-get update
		`)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Peringatan: Tidak dapat menambahkan repositori PHP: %s\n", err)
		}
	}

	fmt.Println("Repositori PHP berhasil diinstal")
}

// installPHPVersions menginstal versi PHP yang dipilih pengguna
func installPHPVersions() {
	fmt.Println("\nVersi PHP yang tersedia:")

	// Dapatkan daftar versi PHP yang tersedia
	cmd := exec.Command("apt", "list", "php*-fpm")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Peringatan: Tidak dapat mendapatkan daftar versi PHP: %s\n", err)
		return
	}

	// Parse output untuk mendapatkan versi PHP
	versions := parsePhpVersions(string(output))

	if len(versions) == 0 {
		fmt.Println("Tidak ada versi PHP yang tersedia")
		return
	}

	// Tampilkan versi yang tersedia
	for i, version := range versions {
		fmt.Printf("%d) PHP %s\n", i+1, version)
	}

	// Tanya pengguna untuk memilih versi
	fmt.Print("\nPilih versi PHP untuk diinstal (pisahkan dengan koma, kosongkan untuk melewati): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		fmt.Println("Melewati instalasi PHP")
		return
	}

	// Parse pilihan pengguna
	choices := strings.Split(input, ",")
	for _, choice := range choices {
		choice = strings.TrimSpace(choice)
		index, err := parseIndex(choice, len(versions))
		if err != nil {
			fmt.Printf("Peringatan: Pilihan tidak valid: %s\n", choice)
			continue
		}

		// Instal versi PHP yang dipilih
		version := versions[index]
		fmt.Printf("Menginstal PHP %s...\n", version)

		cmd := exec.Command("apt-get", "install", "-y",
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
			continue
		}

		// Buat modul Caddy untuk PHP
		createPhpModule(version)

		fmt.Printf("PHP %s berhasil diinstal\n", version)
	}
}

// installComposer menginstal Composer
func installComposer() {
	fmt.Print("\nApakah Anda ingin menginstal Composer? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if strings.ToLower(input) != "y" {
		fmt.Println("Melewati instalasi Composer")
		return
	}

	fmt.Println("Menginstal Composer...")

	cmd := exec.Command("sh", "-c", `
		php -r "copy('https://getcomposer.org/installer', 'composer-setup.php');"
		php composer-setup.php --install-dir=/usr/local/bin --filename=composer
		php -r "unlink('composer-setup.php');"
	`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: Tidak dapat menginstal Composer: %s\n", err)
		return
	}

	fmt.Println("Composer berhasil diinstal")
}

// installNodeJS menginstal Node.js
func installNodeJS() {
	fmt.Print("\nApakah Anda ingin menginstal Node.js? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if strings.ToLower(input) != "y" {
		fmt.Println("Melewati instalasi Node.js")
		return
	}

	fmt.Println("Menginstal Node.js...")

	cmd := exec.Command("sh", "-c", `
		curl -fsSL https://deb.nodesource.com/setup_lts.x | bash -
		apt-get install -y nodejs
	`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: Tidak dapat menginstal Node.js: %s\n", err)
		return
	}

	fmt.Println("Node.js berhasil diinstal")
}

// setupDirectories membuat direktori yang diperlukan
func setupDirectories() {
	fmt.Println("\nMembuat direktori yang diperlukan...")

	// Buat direktori Caddy
	if err := os.MkdirAll("/etc/caddy/sites.d", 0755); err != nil {
		fmt.Printf("Error: Tidak dapat membuat direktori /etc/caddy/sites.d: %s\n", err)
	}

	if err := os.MkdirAll("/etc/caddy/module.d", 0755); err != nil {
		fmt.Printf("Error: Tidak dapat membuat direktori /etc/caddy/module.d: %s\n", err)
	}

	// Buat direktori situs
	if err := os.MkdirAll("/apps/sites", 0755); err != nil {
		fmt.Printf("Error: Tidak dapat membuat direktori /apps/sites: %s\n", err)
	}

	// Buat direktori backup
	if err := os.MkdirAll("/backup/daily", 0755); err != nil {
		fmt.Printf("Error: Tidak dapat membuat direktori /backup/daily: %s\n", err)
	}

	if err := os.MkdirAll("/backup/weekly", 0755); err != nil {
		fmt.Printf("Error: Tidak dapat membuat direktori /backup/weekly: %s\n", err)
	}

	// Buat Caddyfile utama
	caddyfileContent := `{
	# Konfigurasi global
	admin off
	email admin@localhost
}

# Impor semua situs dari direktori sites.d
import sites.d/*.conf
`

	if err := os.WriteFile("/etc/caddy/Caddyfile", []byte(caddyfileContent), 0644); err != nil {
		fmt.Printf("Error: Tidak dapat menulis Caddyfile: %s\n", err)
	}

	fmt.Println("Direktori berhasil dibuat")
}

// setupDefaultModules membuat modul default
func setupDefaultModules() {
	fmt.Println("\nMembuat modul default...")

	modules := map[string]string{
		"spa": `(spa) {
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
}`,
		"security": `(security) {
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
}`,
		"ratelimit": `(ratelimit) {
	rate_limit {
		zone dynamic {
			key {remote_host}
			events 100
			window 10s
		}
	}
}`,
		"compression": `(compression) {
	encode gzip zstd
}`,
		"cache-headers": `(cache-headers) {
	@static {
		path *.css *.js *.png *.jpg *.jpeg *.gif *.ico *.svg *.woff *.woff2 *.ttf *.eot
	}
	header @static Cache-Control "public, max-age=31536000"
}`,
		"local-access": `(local-access) {
	@local {
		remote_ip 127.0.0.1 192.168.0.0/16 10.0.0.0/8 172.16.0.0/12
	}
	@notLocal {
		not remote_ip 127.0.0.1 192.168.0.0/16 10.0.0.0/8 172.16.0.0/12
	}
	handle @notLocal {
		respond "Akses ditolak" 403
	}
}`,
	}

	for name, content := range modules {
		modulePath := fmt.Sprintf("/etc/caddy/module.d/%s", name)
		if err := os.WriteFile(modulePath, []byte(content), 0644); err != nil {
			fmt.Printf("Error: Tidak dapat menulis modul %s: %s\n", name, err)
		}
	}

	fmt.Println("Modul default berhasil dibuat")
}

// createPhpModule membuat modul Caddy untuk PHP
func createPhpModule(version string) {
	moduleContent := fmt.Sprintf(`(php%s) {
	php_fastcgi unix//run/php/php%s-fpm.sock
}
`, version, version)

	modulePath := fmt.Sprintf("/etc/caddy/module.d/php%s", version)
	if err := os.WriteFile(modulePath, []byte(moduleContent), 0644); err != nil {
		fmt.Printf("Peringatan: Tidak dapat membuat modul PHP untuk Caddy: %s\n", err)
	}
}

// parsePhpVersions mengurai versi PHP dari output apt list
func parsePhpVersions(output string) []string {
	versions := []string{}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "php") && strings.Contains(line, "-fpm") {
			// Ekstrak versi dengan cara sederhana
			parts := strings.Split(line, "/")
			if len(parts) > 0 {
				pkgName := parts[0]
				if strings.HasPrefix(pkgName, "php") && strings.HasSuffix(pkgName, "-fpm") {
					version := pkgName[3 : len(pkgName)-4]
					if !contains(versions, version) {
						versions = append(versions, version)
					}
				}
			}
		}
	}

	return versions
}

// parseIndex mengurai indeks dari string
func parseIndex(s string, max int) (int, error) {
	var index int
	_, err := fmt.Sscanf(s, "%d", &index)
	if err != nil || index < 1 || index > max {
		return 0, fmt.Errorf("indeks tidak valid")
	}
	return index - 1, nil
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
