package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gnome-ext-manager/models"
	"gnome-ext-manager/utils"
	"os"
	"os/exec"
	"strings"
	"time"
)

func NewBackupExtensions(outputPath string) error {
	progressBar := utils.NewProgressBar(100, 50)
	progressBar.Start()
	defer progressBar.Stop()

	progressBar.SetMessage("Starting backup process")
	progressBar.SetProgress(5)

	progressBar.SetMessage("Detecting GNOME Shell version")
	versionCmd := exec.Command("bash", "-c", "gnome-shell --version | awk '{print $3}'")
	versionOut, err := versionCmd.Output()
	if err != nil {
		return fmt.Errorf("erro ao obter versão do GNOME Shell: %w", err)
	}
	versionStr := strings.TrimSpace(string(versionOut))
	majorVersion := strings.Split(versionStr, ".")[0]

	progressBar.SetProgress(20)
	fmt.Println("GNOME Shell version detected:", majorVersion)

	progressBar.SetMessage("Listing enabled extensions")
	listCmd := exec.Command("gnome-extensions", "list", "--enabled")
	listOut, err := listCmd.Output()
	if err != nil {
		return fmt.Errorf("erro ao listar extensões: %w", err)
	}
	uuids := strings.Fields(string(listOut))

	progressBar.SetProgress(35)

	var extensions []models.Extension

	progressBar.SetMessage("Processing extensions")
	totalExtensions := len(uuids)
	for i, uuid := range uuids {
		ext := models.Extension{
			UUID:       uuid,
			Enabled:    true,
			URL:        fmt.Sprintf("https://extensions.gnome.org/ajax/detail/?uuid=%s", uuid),
			GnomeShell: majorVersion,
		}
		extensions = append(extensions, ext)

		progress := 35 + int(float64(i+1)/float64(totalExtensions)*25)
		progressBar.SetProgress(progress)

		time.Sleep(100 * time.Millisecond)
	}

	progressBar.SetMessage("Saving extension data")
	progressBar.SetProgress(60)

	jsonData, err := json.MarshalIndent(extensions, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao gerar JSON: %w", err)
	}

	jsonFile := "backup_extensions.json"
	if outputPath != "" {
		jsonFile = outputPath + "_extensions.json"
	}

	err = os.WriteFile(jsonFile, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar arquivo JSON: %w", err)
	}
	fmt.Println("Arquivo JSON salvo em:", jsonFile)

	progressBar.SetProgress(70)

	progressBar.SetMessage("Backing up extension settings")
	dconfFile := "extensions_dconf_backup.txt"
	if outputPath != "" {
		dconfFile = outputPath + "_dconf_backup.txt"
	}

	var dconfBuffer bytes.Buffer
	dconfCmd := exec.Command("dconf", "dump", "/org/gnome/shell/extensions/")
	dconfCmd.Stdout = &dconfBuffer
	err = dconfCmd.Run()
	if err != nil {
		fmt.Println("Erro ao fazer backup do dconf:", err)
	} else {
		err = os.WriteFile(dconfFile, dconfBuffer.Bytes(), 0644)
		if err != nil {
			fmt.Println("Erro ao salvar backup do dconf:", err)
		} else {
			fmt.Println("Backup das configurações salvo em:", dconfFile)
		}
	}

	progressBar.SetProgress(90)

	progressBar.SetMessage("Creating zip archive")
	zipFile := "backup_extensions.zip"
	if outputPath != "" {
		zipFile = outputPath + ".zip"
	}

	files := []string{jsonFile, dconfFile}
	err = utils.CreateZip(zipFile, files)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo zip: %w", err)
	}

	progressBar.SetMessage("Backup complete")
	progressBar.SetProgress(100)

	return nil
}
