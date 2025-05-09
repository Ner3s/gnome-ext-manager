package main

import (
	"fmt"
	"gnome-ext-manager/services"
)

func main() {
	option := 0

	for option != 3 {
		fmt.Println("\n###############################")
		fmt.Println("##     GNOME-EXT-MANAGER     ##")
		fmt.Println("###############################")
		fmt.Println("1. Backup GNOME extensions")
		fmt.Println("2. Restore GNOME extensions")
		fmt.Println("3. Exit")
		fmt.Print("Choose an option: ")
		fmt.Scan(&option)

		switch option {
		case 1:
			fmt.Println("\nBackup GNOME Extensions")
			var outputPath string
			fmt.Print("Enter backup file name (without extension) or press Enter for default: ")
			fmt.Scanln(&outputPath)

			err := services.NewBackupExtensions(outputPath)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("\nBackup completed successfully!")
			}
		case 2:
			fmt.Println("\nRestore GNOME Extensions")
			var zipFilePath string
			fmt.Print("Enter the backup file path (.zip): ")
			fmt.Scanln(&zipFilePath)

			err := services.NewRestoreExtensions(zipFilePath)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("\nRestore completed successfully!")
			}
		case 3:
			fmt.Println("Exiting...")
		default:
			fmt.Println("Invalid option, please try again.")
		}
	}
}
