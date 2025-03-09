package backup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	sitesDir        = "/apps/sites"
	backupDailyDir  = "/backup/daily"
	backupWeeklyDir = "/backup/weekly"
	cronFile        = "/etc/cron.d/webpanel-backup"
)

// Enable mengaktifkan backup untuk domain tertentu
func Enable(backupType, domain string) {
	// Validasi tipe backup
	if backupType != "daily" && backupType != "weekly" {
		fmt.Printf("Error: Tipe backup tidak valid: %s (harus daily atau weekly)\n", backupType)
		return
	}

	// Validasi domain
	siteDir := filepath.Join(sitesDir, domain)
	if _, err := os.Stat(siteDir); os.IsNotExist(err) {
		fmt.Printf("Error: Domain tidak ditemukan: %s\n", domain)
		return
	}

	// Buat direktori backup jika belum ada
	var backupDir string
	if backupType == "daily" {
		backupDir = filepath.Join(backupDailyDir, domain)
	} else {
		backupDir = filepath.Join(backupWeeklyDir, domain)
	}

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		fmt.Printf("Error: Tidak dapat membuat direktori backup: %s\n", err)
		return
	}

	// Tambahkan ke cron
	if err := addToCron(backupType, domain); err != nil {
		fmt.Printf("Error: Tidak dapat menambahkan ke cron: %s\n", err)
		return
	}

	fmt.Printf("Backup %s untuk %s berhasil diaktifkan\n", backupType, domain)
}

// Disable menonaktifkan backup untuk domain tertentu
func Disable(backupType, domain string) {
	// Validasi tipe backup
	if backupType != "daily" && backupType != "weekly" {
		fmt.Printf("Error: Tipe backup tidak valid: %s (harus daily atau weekly)\n", backupType)
		return
	}

	// Hapus dari cron
	if err := removeFromCron(backupType, domain); err != nil {
		fmt.Printf("Error: Tidak dapat menghapus dari cron: %s\n", err)
		return
	}

	fmt.Printf("Backup %s untuk %s berhasil dinonaktifkan\n", backupType, domain)
}

// AddDBBackup menambahkan backup database
func AddDBBackup(dbName string) {
	// Validasi nama database
	if !isValidDBName(dbName) {
		fmt.Printf("Error: Nama database tidak valid: %s\n", dbName)
		return
	}

	// Tambahkan ke cron
	if err := addDBToCron(dbName); err != nil {
		fmt.Printf("Error: Tidak dapat menambahkan backup database ke cron: %s\n", err)
		return
	}

	fmt.Printf("Backup database untuk %s berhasil diaktifkan\n", dbName)
}

// addToCron menambahkan tugas backup ke cron
func addToCron(backupType, domain string) error {
	// Baca file cron yang ada
	var cronContent string
	if _, err := os.Stat(cronFile); !os.IsNotExist(err) {
		content, err := ioutil.ReadFile(cronFile)
		if err != nil {
			return err
		}
		cronContent = string(content)
	}

	// Buat perintah backup
	var cronLine string
	if backupType == "daily" {
		cronLine = fmt.Sprintf("0 2 * * * root rsync -a --delete %s/ %s/\n",
			filepath.Join(sitesDir, domain),
			filepath.Join(backupDailyDir, domain))
	} else {
		cronLine = fmt.Sprintf("0 3 * * 0 root rsync -a --delete %s/ %s/\n",
			filepath.Join(sitesDir, domain),
			filepath.Join(backupWeeklyDir, domain))
	}

	// Periksa apakah sudah ada
	if strings.Contains(cronContent, cronLine) {
		return nil // Sudah ada, tidak perlu menambahkan lagi
	}

	// Tambahkan baris baru
	cronContent += cronLine

	// Tulis kembali file cron
	return ioutil.WriteFile(cronFile, []byte(cronContent), 0644)
}

// removeFromCron menghapus tugas backup dari cron
func removeFromCron(backupType, domain string) error {
	// Baca file cron yang ada
	if _, err := os.Stat(cronFile); os.IsNotExist(err) {
		return nil // File tidak ada, tidak perlu menghapus
	}

	content, err := ioutil.ReadFile(cronFile)
	if err != nil {
		return err
	}
	cronContent := string(content)

	// Buat pola untuk mencari baris yang akan dihapus
	var searchPattern string
	if backupType == "daily" {
		searchPattern = fmt.Sprintf("rsync -a --delete %s/ %s/",
			filepath.Join(sitesDir, domain),
			filepath.Join(backupDailyDir, domain))
	} else {
		searchPattern = fmt.Sprintf("rsync -a --delete %s/ %s/",
			filepath.Join(sitesDir, domain),
			filepath.Join(backupWeeklyDir, domain))
	}

	// Hapus baris yang cocok
	lines := strings.Split(cronContent, "\n")
	newLines := []string{}
	for _, line := range lines {
		if !strings.Contains(line, searchPattern) {
			newLines = append(newLines, line)
		}
	}
	newContent := strings.Join(newLines, "\n")

	// Tulis kembali file cron
	return ioutil.WriteFile(cronFile, []byte(newContent), 0644)
}

// addDBToCron menambahkan tugas backup database ke cron
func addDBToCron(dbName string) error {
	// Baca file cron yang ada
	var cronContent string
	if _, err := os.Stat(cronFile); !os.IsNotExist(err) {
		content, err := ioutil.ReadFile(cronFile)
		if err != nil {
			return err
		}
		cronContent = string(content)
	}

	// Buat perintah backup
	backupDir := filepath.Join(backupDailyDir, "databases")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return err
	}

	cronLine := fmt.Sprintf("0 4 * * * root mysqldump -u root %s > %s/%s.sql\n",
		dbName, backupDir, dbName)

	// Periksa apakah sudah ada
	if strings.Contains(cronContent, cronLine) {
		return nil // Sudah ada, tidak perlu menambahkan lagi
	}

	// Tambahkan baris baru
	cronContent += cronLine

	// Tulis kembali file cron
	return ioutil.WriteFile(cronFile, []byte(cronContent), 0644)
}

// isValidDBName memeriksa apakah nama database valid
func isValidDBName(dbName string) bool {
	// Implementasi sederhana, bisa ditingkatkan dengan validasi yang lebih baik
	return len(dbName) > 0 && !strings.Contains(dbName, " ") && !strings.Contains(dbName, "/")
}
