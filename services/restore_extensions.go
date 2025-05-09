package services

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"gnome-ext-manager/models"
	"gnome-ext-manager/utils"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func unzip(src, dest string, progressBar *utils.ProgressBar) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// Set total for progress bar
	// fileCount := len(r.File)
	progressBar.SetMessage("Extracting files")

	for i, f := range r.File {
		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		os.MkdirAll(filepath.Dir(path), os.ModePerm)

		outFile, err := os.Create(path)
		if err != nil {
			return err
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}

		// Update progress
		progressBar.SetProgress(i + 1)
	}
	return nil
}

func getGnomeShellVersion() (string, error) {
	cmd := exec.Command("bash", "-c", "gnome-shell --version | awk '{print $3}'")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	version := strings.TrimSpace(string(out))
	major := strings.Split(version, ".")[0]
	return major, nil
}

func restoreExtensions(jsonPath string, tempDir string, progressBar *utils.ProgressBar) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return err
	}

	var extensions []models.Extension
	if err := json.Unmarshal(data, &extensions); err != nil {
		return err
	}

	gnomeVersion, err := getGnomeShellVersion()
	if err != nil {
		return err
	}
	fmt.Println("Versão do GNOME Shell detectada:", gnomeVersion)

	// Set total for progress bar
	progressBar.SetMessage("Installing extensions")
	progressBar.SetProgress(0)
	// progressBar.total = len(extensions)

	for i, ext := range extensions {
		resp, err := http.Get(ext.URL)
		if err != nil {
			fmt.Println("Erro ao consultar API para", ext.UUID, ":", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Println("Erro HTTP para", ext.UUID, "status:", resp.StatusCode)
			continue
		}

		var apiData map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&apiData); err != nil {
			fmt.Println("Erro ao decodificar JSON da API para", ext.UUID, ":", err)
			continue
		}

		// encontrar pk na shell_version_map
		shellMap := apiData["shell_version_map"].(map[string]interface{})
		// var pk float64
		// if val, ok := shellMap[gnomeVersion]; ok {
		// 	pk = val.(map[string]interface{})["pk"].(float64)
		// } else {
		// 	fmt.Println("Versão", gnomeVersion, "não encontrada para extensão", ext.UUID)
		// 	continue
		// }

		versionNum := int(shellMap[gnomeVersion].(map[string]interface{})["version"].(float64))

		cleanUUID := strings.ReplaceAll(ext.UUID, "@", "")
		downloadURL := fmt.Sprintf("https://extensions.gnome.org/extension-data/%s.v%d.shell-extension.zip", cleanUUID, versionNum)

		respZip, err := http.Get(downloadURL)
		if err != nil {
			fmt.Println("Erro ao baixar extensão:", downloadURL, err)
			continue
		}
		defer respZip.Body.Close()

		if respZip.StatusCode != 200 {
			fmt.Println("Erro HTTP ao baixar extensão:", respZip.StatusCode)
			continue
		}

		zipPath := filepath.Join(tempDir, fmt.Sprintf("%s.v%d.shell-extension.zip", cleanUUID, versionNum))
		outFile, err := os.Create(zipPath)
		if err != nil {
			fmt.Println("Erro ao criar arquivo:", zipPath, err)
			continue
		}

		_, err = io.Copy(outFile, respZip.Body)
		outFile.Close()
		if err != nil {
			fmt.Println("Erro ao salvar arquivo da extensão:", err)
			continue
		}

		cmd := exec.Command("gnome-extensions", "install", zipPath, "--force")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Erro ao instalar extensão:", err)
		}

		cmdEnable := exec.Command("gnome-extensions", "enable", ext.UUID)
		cmdEnable.Stdout = os.Stdout
		cmdEnable.Stderr = os.Stderr

		progressBar.SetProgress(i + 1)
	}

	return nil
}

func restoreDconf(dconfPath string) error {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("dconf load /org/gnome/shell/extensions/ < %s", dconfPath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func NewRestoreExtensions(zipFilePath string) error {
	progressBar := utils.NewProgressBar(100, 50)
	progressBar.Start()
	defer progressBar.Stop()

	if zipFilePath == "" {
		fmt.Print("Informe o caminho do arquivo de backup (.zip): ")
		fmt.Scanln(&zipFilePath)
	}

	tempDir, err := os.MkdirTemp("", "gnome-restore")
	if err != nil {
		return fmt.Errorf("erro ao criar diretório temporário: %w", err)
	}
	defer os.RemoveAll(tempDir)

	progressBar.SetMessage("Extracting backup archive")
	progressBar.SetProgress(10)
	err = unzip(zipFilePath, tempDir, progressBar)
	if err != nil {
		return fmt.Errorf("erro ao extrair arquivo zip: %w", err)
	}

	jsonPath := filepath.Join(tempDir, "backup_extensions.json")
	dconfPath := filepath.Join(tempDir, "extensions_dconf_backup.txt")

	progressBar.SetProgress(30)
	err = restoreExtensions(jsonPath, tempDir, progressBar)
	if err != nil {
		fmt.Println("Erro ao restaurar extensões:", err)
	}

	progressBar.SetMessage("Restoring extension settings")
	progressBar.SetProgress(80)
	err = restoreDconf(dconfPath)
	if err != nil {
		fmt.Println("Erro ao restaurar configurações dconf:", err)
	}

	progressBar.SetMessage("Finalizing")
	progressBar.SetProgress(100)
	return nil
}
