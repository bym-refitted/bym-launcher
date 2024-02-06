package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func downloadFile(filePath string, url string, useHttps bool) error {
	fullUrl := fmt.Sprintf("http%s://%s%s", map[bool]string{true: "s", false: ""}[useHttps], downloadBasePath, url)
	resp, err := http.Get(fullUrl)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Could not connect & download file - Status: %v %s", resp.StatusCode, url)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}
	return nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func ensureFolderExists(folder string) error {
	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		fmt.Printf("Creating %v folder\n", folder)
		err := os.Mkdir(folder, 0755)
		if err != nil {
			return fmt.Errorf("failed to create %v folder: %v", folder, err)
		}
	}
	return nil
}

func getLocalVersions() (bool, LocalVersionManifest, error) {
	versionFilePath := filepath.Join(downloadsFolder, "version.json")
	if !fileExists(versionFilePath) {
		return false, LocalVersionManifest{}, nil
	}

	file, err := os.Open(versionFilePath)
	if err != nil {
		return false, LocalVersionManifest{}, fmt.Errorf("failed to read version.json file: %v", err)
	}
	defer file.Close()

	var localManifest LocalVersionManifest
	if err := json.NewDecoder(file).Decode(&localManifest); err != nil {
		return false, LocalVersionManifest{}, fmt.Errorf("Failed to decode version.json file: %v", err)
	}

	return true, localManifest, nil
}

func setLocalVersions(localManifest LocalVersionManifest) error {
	versionFilePath := filepath.Join(downloadsFolder, "version.json")
	file, err := os.Create(versionFilePath)
	if err != nil {
		return fmt.Errorf("Failed to create local version manifest: %v", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(localManifest); err != nil {
		return fmt.Errorf("Failed to encode local version manifest: %v", err)
	}
	return nil
}
