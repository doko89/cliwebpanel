package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/doko89/webpanel/internal/backup"
	"github.com/doko89/webpanel/internal/database"
	"github.com/doko89/webpanel/internal/module"
	"github.com/doko89/webpanel/internal/php"
	"github.com/doko89/webpanel/internal/proxy"
	"github.com/doko89/webpanel/internal/site"
	"github.com/doko89/webpanel/internal/utils"
)

func main() {
	fmt.Println("WebPanel CLI - Server Administration Tool")

	// Periksa apakah berjalan sebagai root
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error: Tidak dapat menentukan pengguna saat ini:", err)
		os.Exit(1)
	}

	if currentUser.Uid != "0" {
		fmt.Println("Error: Webpanel harus dijalankan sebagai root")
		os.Exit(1)
	}

	// Periksa argumen
	if len(os.Args) < 2 {
		displayHelp()
		os.Exit(1)
	}

	// Proses perintah
	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "site":
		handleSiteCommand(args)
	case "proxy":
		handleProxyCommand(args)
	case "module":
		handleModuleCommand(args)
	case "backup":
		handleBackupCommand(args)
	case "db", "database":
		handleDatabaseCommand(args)
	case "php":
		handlePHPCommand(args)
	case "install":
		handleInstallCommand(args)
	case "help":
		displayHelp()
	default:
		fmt.Printf("Error: Perintah tidak dikenal: %s\n", command)
		displayHelp()
		os.Exit(1)
	}
}

func handleSiteCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Subperintah site diperlukan")
		printSiteHelp()
		os.Exit(1)
	}

	subcommand := args[0]
	switch subcommand {
	case "add":
		if len(args) < 2 {
			fmt.Println("Error: Domain diperlukan")
			printSiteHelp()
			os.Exit(1)
		}
		site.Add(args[1])
	case "remove":
		if len(args) < 2 {
			fmt.Println("Error: Domain diperlukan")
			printSiteHelp()
			os.Exit(1)
		}
		site.Remove(args[1])
	case "list":
		site.List()
	default:
		fmt.Printf("Error: Subperintah site tidak dikenal: %s\n", subcommand)
		printSiteHelp()
		os.Exit(1)
	}
}

func handleProxyCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Subperintah proxy diperlukan")
		printProxyHelp()
		os.Exit(1)
	}

	subcommand := args[0]
	switch subcommand {
	case "add":
		if len(args) < 3 {
			fmt.Println("Error: Domain dan target diperlukan")
			printProxyHelp()
			os.Exit(1)
		}
		proxy.Add(args[1], args[2])
	case "remove":
		if len(args) < 2 {
			fmt.Println("Error: Domain diperlukan")
			printProxyHelp()
			os.Exit(1)
		}
		proxy.Remove(args[1])
	case "list":
		proxy.List()
	default:
		fmt.Printf("Error: Subperintah proxy tidak dikenal: %s\n", subcommand)
		printProxyHelp()
		os.Exit(1)
	}
}

func handleModuleCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Subperintah module diperlukan")
		printModuleHelp()
		os.Exit(1)
	}

	subcommand := args[0]
	switch subcommand {
	case "enable":
		if len(args) < 3 {
			fmt.Println("Error: Nama modul dan domain diperlukan")
			printModuleHelp()
			os.Exit(1)
		}
		module.Enable(args[1], args[2])
	case "disable":
		if len(args) < 3 {
			fmt.Println("Error: Nama modul dan domain diperlukan")
			printModuleHelp()
			os.Exit(1)
		}
		module.Disable(args[1], args[2])
	case "list":
		if len(args) < 2 {
			fmt.Println("Error: Domain diperlukan")
			printModuleHelp()
			os.Exit(1)
		}
		module.List(args[1])
	case "list-available":
		module.ListAvailable()
	default:
		fmt.Printf("Error: Subperintah module tidak dikenal: %s\n", subcommand)
		printModuleHelp()
		os.Exit(1)
	}
}

func handleBackupCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Subperintah backup diperlukan")
		printBackupHelp()
		os.Exit(1)
	}

	subcommand := args[0]
	switch subcommand {
	case "enable":
		if len(args) < 3 {
			fmt.Println("Error: Jenis backup dan domain diperlukan")
			printBackupHelp()
			os.Exit(1)
		}
		backup.Enable(args[1], args[2])
	case "disable":
		if len(args) < 3 {
			fmt.Println("Error: Jenis backup dan domain diperlukan")
			printBackupHelp()
			os.Exit(1)
		}
		backup.Disable(args[1], args[2])
	case "dbbackup":
		if len(args) < 3 && args[1] == "add" {
			fmt.Println("Error: Nama database diperlukan")
			printBackupHelp()
			os.Exit(1)
		}
		if args[1] == "add" {
			backup.AddDBBackup(args[2])
		} else {
			fmt.Printf("Error: Subperintah dbbackup tidak dikenal: %s\n", args[1])
			printBackupHelp()
			os.Exit(1)
		}
	default:
		fmt.Printf("Error: Subperintah backup tidak dikenal: %s\n", subcommand)
		printBackupHelp()
		os.Exit(1)
	}
}

func handleDatabaseCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Subperintah database diperlukan")
		printDatabaseHelp()
		os.Exit(1)
	}

	subcommand := args[0]
	switch subcommand {
	case "create":
		if len(args) < 4 {
			fmt.Println("Error: Nama database, pengguna, dan kata sandi diperlukan")
			printDatabaseHelp()
			os.Exit(1)
		}
		database.Create(args[1], args[2], args[3])
	case "delete":
		if len(args) < 2 {
			fmt.Println("Error: Nama database diperlukan")
			printDatabaseHelp()
			os.Exit(1)
		}
		database.Delete(args[1])
	default:
		fmt.Printf("Error: Subperintah database tidak dikenal: %s\n", subcommand)
		printDatabaseHelp()
		os.Exit(1)
	}
}

func handlePHPCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Subperintah php diperlukan")
		printPHPHelp()
		os.Exit(1)
	}

	subcommand := args[0]
	switch subcommand {
	case "list":
		php.List()
	case "installed":
		php.ListInstalled()
	case "install":
		if len(args) < 2 {
			fmt.Println("Error: Versi PHP diperlukan")
			printPHPHelp()
			os.Exit(1)
		}
		php.Install(args[1])
	case "uninstall":
		if len(args) < 2 {
			fmt.Println("Error: Versi PHP diperlukan")
			printPHPHelp()
			os.Exit(1)
		}
		php.Uninstall(args[1])
	case "module":
		if len(args) < 2 {
			fmt.Println("Error: Subperintah module diperlukan")
			printPHPModuleHelp()
			os.Exit(1)
		}

		moduleSubcommand := args[1]
		switch moduleSubcommand {
		case "list":
			if len(args) < 3 {
				fmt.Println("Error: Versi PHP diperlukan")
				printPHPModuleHelp()
				os.Exit(1)
			}
			php.ListModules(args[2])
		case "install":
			if len(args) < 3 {
				fmt.Println("Error: Nama modul PHP diperlukan")
				printPHPModuleHelp()
				os.Exit(1)
			}
			php.InstallModule(args[2])
		default:
			fmt.Printf("Error: Subperintah php module tidak dikenal: %s\n", moduleSubcommand)
			printPHPModuleHelp()
			os.Exit(1)
		}
	default:
		fmt.Printf("Error: Subperintah php tidak dikenal: %s\n", subcommand)
		printPHPHelp()
		os.Exit(1)
	}
}

func handleInstallCommand(args []string) {
	// Implementasi instalasi
	utils.InstallDependencies()
}

// Fungsi bantuan untuk mencetak dokumentasi
func displayHelp() {
	fmt.Println("Usage: webpanel [command] [options]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  site       Manage websites")
	fmt.Println("  proxy      Manage proxy configurations")
	fmt.Println("  module     Manage Caddy modules")
	fmt.Println("  backup     Manage backup configurations")
	fmt.Println("  db         Manage databases")
	fmt.Println("  php        Manage PHP installations")
	fmt.Println("  help       Display help information")
	fmt.Println("")
	fmt.Println("Run 'webpanel help [command]' for more information on a command.")
}

func printSiteHelp() {
	fmt.Println("Penggunaan: webpanel site <subperintah> [argumen...]")
	fmt.Println("\nSubperintah yang tersedia:")
	fmt.Println("  add <domain>      Menambahkan situs baru")
	fmt.Println("  remove <domain>   Menghapus situs")
	fmt.Println("  list              Menampilkan daftar situs")
}

func printProxyHelp() {
	fmt.Println("Penggunaan: webpanel proxy <subperintah> [argumen...]")
	fmt.Println("\nSubperintah yang tersedia:")
	fmt.Println("  add <domain> <target>   Menambahkan situs proxy baru")
	fmt.Println("  remove <domain>         Menghapus situs proxy")
	fmt.Println("  list                    Menampilkan daftar situs proxy")
}

func printModuleHelp() {
	fmt.Println("Penggunaan: webpanel module <subperintah> [argumen...]")
	fmt.Println("\nSubperintah yang tersedia:")
	fmt.Println("  enable <module> <domain>    Mengaktifkan modul untuk domain")
	fmt.Println("  disable <module> <domain>   Menonaktifkan modul untuk domain")
	fmt.Println("  list <domain>               Menampilkan modul yang diaktifkan untuk domain")
	fmt.Println("  list-available              Menampilkan semua modul yang tersedia")
}

func printBackupHelp() {
	fmt.Println("Penggunaan: webpanel backup <subperintah> [argumen...]")
	fmt.Println("\nSubperintah yang tersedia:")
	fmt.Println("  enable <daily|weekly> <domain>    Mengaktifkan backup untuk domain")
	fmt.Println("  disable <daily|weekly> <domain>   Menonaktifkan backup untuk domain")
	fmt.Println("  dbbackup add <dbname>             Menambahkan backup database")
}

func printDatabaseHelp() {
	fmt.Println("Penggunaan: webpanel db <subperintah> [argumen...]")
	fmt.Println("\nSubperintah yang tersedia:")
	fmt.Println("  create <database> <user> <password>   Membuat database baru")
	fmt.Println("  delete <database>                     Menghapus database")
}

func printPHPHelp() {
	fmt.Println("Penggunaan: webpanel php <subperintah> [argumen...]")
	fmt.Println("\nSubperintah yang tersedia:")
	fmt.Println("  list                 Menampilkan semua versi PHP yang tersedia")
	fmt.Println("  installed            Menampilkan versi PHP yang terinstal")
	fmt.Println("  install <version>    Menginstal versi PHP tertentu")
	fmt.Println("  uninstall <version>  Menghapus instalasi versi PHP tertentu")
	fmt.Println("  module               Mengelola modul PHP")
}

func printPHPModuleHelp() {
	fmt.Println("Penggunaan: webpanel php module <subperintah> [argumen...]")
	fmt.Println("\nSubperintah yang tersedia:")
	fmt.Println("  list <version>       Menampilkan modul PHP yang tersedia")
	fmt.Println("  install <module>     Menginstal modul PHP")
}
