package database

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Create membuat database baru dengan pengguna dan kata sandi yang ditentukan
func Create(dbName, dbUser, dbPassword string) {
	// Validasi input
	if !isValidName(dbName) || !isValidName(dbUser) {
		fmt.Println("Error: Nama database dan pengguna hanya boleh berisi huruf, angka, dan garis bawah")
		return
	}

	// Buat database
	createDBCmd := exec.Command("mysql", "-e", fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", dbName))
	if err := createDBCmd.Run(); err != nil {
		fmt.Printf("Error: Tidak dapat membuat database: %s\n", err)
		return
	}

	// Buat pengguna dan berikan hak akses
	createUserCmd := exec.Command("mysql", "-e", fmt.Sprintf(
		"CREATE USER IF NOT EXISTS '%s'@'localhost' IDENTIFIED BY '%s'; "+
			"GRANT ALL PRIVILEGES ON `%s`.* TO '%s'@'localhost'; "+
			"FLUSH PRIVILEGES;",
		dbUser, dbPassword, dbName, dbUser))
	if err := createUserCmd.Run(); err != nil {
		fmt.Printf("Error: Tidak dapat membuat pengguna database: %s\n", err)
		return
	}

	fmt.Printf("Database %s dan pengguna %s berhasil dibuat\n", dbName, dbUser)
}

// Delete menghapus database dan penggunanya
func Delete(dbName string) {
	// Validasi input
	if !isValidName(dbName) {
		fmt.Println("Error: Nama database hanya boleh berisi huruf, angka, dan garis bawah")
		return
	}

	// Konfirmasi penghapusan
	fmt.Printf("PERINGATAN: Anda akan menghapus database %s dan semua datanya.\n", dbName)
	fmt.Printf("Ketik nama database untuk mengkonfirmasi: ")
	reader := bufio.NewReader(os.Stdin)
	confirmation, _ := reader.ReadString('\n')
	confirmation = strings.TrimSpace(confirmation)

	if confirmation != dbName {
		fmt.Println("Penghapusan dibatalkan: Konfirmasi tidak cocok")
		return
	}

	// Dapatkan pengguna yang terkait dengan database
	getUsersCmd := exec.Command("mysql", "-N", "-e", fmt.Sprintf(
		"SELECT user FROM mysql.db WHERE db='%s' AND host='localhost';", dbName))
	output, err := getUsersCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: Tidak dapat mendapatkan pengguna database: %s\n", err)
		return
	}

	users := strings.Split(strings.TrimSpace(string(output)), "\n")

	// Hapus database
	dropDBCmd := exec.Command("mysql", "-e", fmt.Sprintf("DROP DATABASE IF EXISTS `%s`;", dbName))
	if err := dropDBCmd.Run(); err != nil {
		fmt.Printf("Error: Tidak dapat menghapus database: %s\n", err)
		return
	}

	// Hapus pengguna
	for _, user := range users {
		user = strings.TrimSpace(user)
		if user == "" {
			continue
		}

		dropUserCmd := exec.Command("mysql", "-e", fmt.Sprintf(
			"DROP USER IF EXISTS '%s'@'localhost';", user))
		if err := dropUserCmd.Run(); err != nil {
			fmt.Printf("Peringatan: Tidak dapat menghapus pengguna %s: %s\n", user, err)
		}
	}

	fmt.Printf("Database %s dan penggunanya berhasil dihapus\n", dbName)
}

// isValidName memeriksa apakah nama database atau pengguna valid
func isValidName(name string) bool {
	// Implementasi sederhana, bisa ditingkatkan dengan validasi regex yang lebih baik
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_') {
			return false
		}
	}
	return len(name) > 0
}
