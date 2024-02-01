package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const buildFolder = "bymr"
const latestBuildUrl = "https://api.github.com/repos/bym-refitted/backyard-monsters-refitted/releases/latest"

type Assets struct {
	Url  string `json:"browser_download_url"`
	Name string `json:"name"`
}

type latestBuild struct {
	ID     int      `json:"id"`
	Assets []Assets `json:"assets"`
}

func getLatestBuild() (latestBuild, error) {
	var latest latestBuild

	resp, err := http.Get(latestBuildUrl)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != http.StatusOK {
		return latest, fmt.Errorf("failed to fetch latest release: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&latest)
	if err != nil {
		return latest, err
	}
	return latest, nil
}

func createBuildFolderAndVersionFile() (int, error) {
	_, err := os.Stat(buildFolder)
	if os.IsNotExist(err) {
		err := os.Mkdir(buildFolder, 0755)
		if err != nil {
			return 0, fmt.Errorf("failed to create bymr folder: %v", err)
		}
	}

	// Check if "version.txt" file exists
	versionFilePath := filepath.Join(buildFolder, "version.txt")
	if !fileExists(versionFilePath) {
		// "version.txt" file does not exist, create it
		file, err := os.Create(versionFilePath)
		if err != nil {
			return 0, fmt.Errorf("failed to create version.txt file: %v", err)
		}
		defer file.Close()

		// Write default content to "version.txt" file
		_, err = file.WriteString("0")
		if err != nil {
			return 0, fmt.Errorf("failed to write to version.txt file: %v", err)
		}
	}

	content, err := os.ReadFile(versionFilePath)
	if err != nil {
		return 0, fmt.Errorf("failed to read version.txt file: %v", err)
	}

	version, err := strconv.Atoi(string(content))
	if err != nil {
		return 0, fmt.Errorf("failed to parse version as integer: %v", err)
	}

	return version, err
}

func downloadLatestBuild(url string, fileName string) (string, error) {
	// Ensure the "bymr" folder exists
	err := os.MkdirAll(buildFolder, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create bymr folder: %v", err)
	}

	// Construct the path for the downloaded file within the "bymr" folder
	filePath := filepath.Join(buildFolder, fileName)

	// Send GET request to download the latest build
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download latest build: %v", err)
	}
	defer resp.Body.Close()

	// Create the file to save the downloaded build
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	// Copy the downloaded content to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write to file: %v", err)
	}

	// Return the absolute path to the downloaded file
	return filepath.Abs(filePath)
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func patcher() (latestBuild, error) {
	lVersion, err := getLatestBuild()
	if err != nil {
		fmt.Println("Cannot get latest build d", err)
		return lVersion, err
	}
	cVersion, _ := createBuildFolderAndVersionFile()
	fmt.Printf("Current version: %d | Latest version: %d \n", cVersion, lVersion.ID)

	if cVersion != lVersion.ID {
		fmt.Println("Downloading latest build")
		for _, asset := range lVersion.Assets {
			name := asset.Name

			if strings.Contains(name, "local") {
				name = "bymr-local.swf"
			} else if strings.Contains(name, "http") {
				name = "bymr-http.swf"
			} else if strings.Contains(name, "stable") {
				name = "bymr-stable.swf"
			}

			filePath, err := downloadLatestBuild(asset.Url, name)
			if err != nil {
				fmt.Printf("Error downloading build %s: %v\n", asset.Name, err)
				continue
			}
			fmt.Printf("Build %s downloaded successfully to: %s\n", asset.Name, filePath)
		}

		versionFilePath := filepath.Join(buildFolder, "version.txt")
		err = os.WriteFile(versionFilePath, []byte(strconv.Itoa(lVersion.ID)), 0644)
		if err != nil {
			fmt.Printf("failed to update version.txt: %v \n", err)
			return lVersion, nil
		}
	}

	fmt.Println("HERE", lVersion.ID, cVersion)
	return lVersion, nil
}
