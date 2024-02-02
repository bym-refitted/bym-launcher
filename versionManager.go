package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// const buildFolder = "bymr"
// const runtimeFolder = "flashRuntimes"
const versionInfoPathBase = "api.bymrefitted.com/launcher.json"


type VersionManifest struct {
    CurrentGameVersion    string `json:"currentGameVersion"`
    CurrentLauncherVersion string `json:"currentLauncherVersion"`
    Builds                Builds `json:"builds"`
    FlashRuntimes         FlashRuntimes `json:"flashRuntimes"`
	httpsWorked bool
}

func (v VersionManifest) String() string {
    return fmt.Sprintf("CurrentGameVersion: %s, CurrentLauncherVersion: %s, Builds: %+v, FlashRuntimes: %+v, httpsWorked: %v",
        v.CurrentGameVersion, v.CurrentLauncherVersion, v.Builds, v.FlashRuntimes, v.httpsWorked)
}

type Builds struct {
    Stable string `json:"stable"`
    Http   string `json:"http"`
    Local  string `json:"local"`
}

type FlashRuntimes struct {
    Windows string `json:"windows"`
    Darwin  string `json:"darwin"`
    Linux   string `json:"linux"`
}

func getVersionInfo(ctx context.Context) (VersionManifest,error) {
	httpsWorked := false
	// First we try https 
		resp, err := http.Get("https://"+versionInfoPathBase)
		if err != nil {
			runtime.EventsEmit(ctx, "infoLog", fmt.Sprintf("Could not access over https, attempting http %s", err))
			// try via http if that fails
			resp, err = http.Get("http://"+versionInfoPathBase)
			if err != nil {
				runtime.EventsEmit(ctx, "infoLog", fmt.Sprintf("Could not access over http, please check the server status on our discord %s", err))
				return VersionManifest{}, err
			}
		}else{
			runtime.EventsEmit(ctx, "infoLog","ONLINE: Launcher successfully connected over https")
			httpsWorked = true
		}
		defer resp.Body.Close()

		// Create a variable to hold the decoded data
		var data VersionManifest

		// Decode the JSON response
		err = json.NewDecoder(resp.Body).Decode(&data)
		data.httpsWorked = httpsWorked
		if err != nil {
			runtime.EventsEmit(ctx, "infoLog","Failed to decode version info")
			return VersionManifest{}, err
		}
	return data, nil
}


// const latestBuildUrl = "https://api.github.com/repos/bym-refitted/backyard-monsters-refitted/releases/latest"

// const flashPlayerArchiveUrl = "https://archive.org/download/standalone_flash_player/Flash%20Player%2032.0%20r101%20%5BWin%5D%20%5BStand%20Alone%5D.exe"

// type Assets struct {
// 	Url  string `json:"browser_download_url"`
// 	Name string `json:"name"`
// }

// type latestBuild struct {
// 	ID     int      `json:"id"`
// 	Assets []Assets `json:"assets"`
// }

// func getLatestBuild() (latestBuild, error) {
// 	var latest latestBuild

// 	// Create a new request
//     req, err := http.NewRequest("GET", latestBuildUrl, nil)
//     if err != nil {
//         log.Fatalln(err)
//     }

//     // Add headers
//     req.Header.Add("User-Agent", "BYM-Launcher")

//     // Create a new HTTP client and send the request
//     client := &http.Client{}
//     resp, err := client.Do(req)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		return latest, fmt.Errorf("failed to fetch latest release: %s", resp.Status)
// 	}

// 	err = json.NewDecoder(resp.Body).Decode(&latest)
// 	// if err != nil {
// 	// 	return latest, err
// 	// }
// 	return latest, nil
// }

// func createBuildFolderAndVersionFile() (int, error) {
// 	err := ensureFolderExists(buildFolder)
// 	err = ensureFolderExists(runtimeFolder)
// 	_, err = downloadRuntime("flashplayer_32.exe")
// 	// Check if "version.txt" file exists
// 	versionFilePath := filepath.Join(buildFolder, "version.txt")
// 	if !fileExists(versionFilePath) {
// 		// "version.txt" file does not exist, create it
// 		file, err := os.Create(versionFilePath)
// 		if err != nil {
// 			return 0, fmt.Errorf("failed to create version.txt file: %v", err)
// 		}
// 		defer file.Close()

// 		// Write default content to "version.txt" file
// 		_, err = file.WriteString("0")
// 		if err != nil {
// 			return 0, fmt.Errorf("failed to write to version.txt file: %v", err)
// 		}
// 	}

// 	content, err := os.ReadFile(versionFilePath)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to read version.txt file: %v", err)
// 	}

// 	version, err := strconv.Atoi(string(content))
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to parse version as integer: %v", err)
// 	}

// 	return version, err
// }


// func ensureFolderExists(folder string) error {
// 	_, err := os.Stat(folder)
// 	if os.IsNotExist(err) {
// 		err := os.Mkdir(folder, 0755)
// 		if err != nil {
// 			return fmt.Errorf("failed to create %v folder: %v", folder, err)
// 		}
// 	}
// 	return nil
// }

// func downloadRuntime(fileName string) (string, error) {
// 	// Ensure the runtimes folder exists
// 	ensureFolderExists(runtimeFolder)
// 	// Construct the path for the downloaded file within the "bymr" folder
// 	filePath := filepath.Join(runtimeFolder, fileName)
// 	if !fileExists(filePath) {
// 		fmt.Printf("Flash player not installed - Downloading")
// 		// Send GET request to download the latest build
// 		resp, err := http.Get(flashPlayerArchiveUrl)
// 		if err != nil {
// 			return "", fmt.Errorf("failed to download latest build: %v", err)
// 		}
// 		defer resp.Body.Close()

// 		// Create the file to save the downloaded build
// 		out, err := os.Create(filePath)
// 		if err != nil {
// 			return "", fmt.Errorf("failed to create file: %v", err)
// 		}
// 		defer out.Close()

// 		// Copy the downloaded content to the file
// 		_, err = io.Copy(out, resp.Body)
// 		if err != nil {
// 			return "", fmt.Errorf("failed to write to file: %v", err)
// 		}
// 		fmt.Printf("Flash player downloaded successfully ")
// 	}
// 	// Return the absolute path to the downloaded file
// 	return filepath.Abs(filePath)
// }

// func downloadLatestBuild(url string, fileName string) (string, error) {
// 	// Ensure the "bymr" folder exists
// 	err := ensureFolderExists(buildFolder)


// 	// Construct the path for the downloaded file within the "bymr" folder
// 	filePath := filepath.Join(buildFolder, fileName)

// 	// Send GET request to download the latest build
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to download latest build: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// Create the file to save the downloaded build
// 	out, err := os.Create(filePath)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to create file: %v", err)
// 	}
// 	defer out.Close()

// 	// Copy the downloaded content to the file
// 	_, err = io.Copy(out, resp.Body)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to write to file: %v", err)
// 	}

// 	// Return the absolute path to the downloaded file
// 	return filepath.Abs(filePath)
// }


// func fileExists(filePath string) bool {
// 	_, err := os.Stat(filePath)
// 	return !os.IsNotExist(err)
// }

// func patcher() (latestBuild, error) {
// 	lVersion, err := getLatestBuild()
// 	if err != nil {
// 		fmt.Println("Cannot get latest build d", err)
// 		return lVersion, err
// 	}
// 	cVersion, _ := createBuildFolderAndVersionFile()
// 	fmt.Printf("Current version: %d | Latest version: %d \n", cVersion, lVersion.ID)
	

// 	if cVersion != lVersion.ID {
// 		fmt.Println("Downloading latest build")
// 		for _, asset := range lVersion.Assets {
// 			name := asset.Name

// 			if strings.Contains(name, "local") {
// 				name = "bymr-local.swf"
// 			} else if strings.Contains(name, "http") {
// 				name = "bymr-http.swf"
// 			} else if strings.Contains(name, "stable") {
// 				name = "bymr-stable.swf"
// 			}

// 			filePath, err := downloadLatestBuild(asset.Url, name)
// 			if err != nil {
// 				fmt.Printf("Error downloading build %s: %v\n", asset.Name, err)
// 				continue
// 			}
// 			fmt.Printf("Build %s downloaded successfully to: %s\n", asset.Name, filePath)
// 		}

// 		versionFilePath := filepath.Join(buildFolder, "version.txt")
// 		err = os.WriteFile(versionFilePath, []byte(strconv.Itoa(lVersion.ID)), 0644)
// 		if err != nil {
// 			fmt.Printf("failed to update version.txt: %v \n", err)
// 			return lVersion, nil
// 		}
// 	}

// 	fmt.Println("HERE", lVersion.ID, cVersion)
// 	return lVersion, nil
// }
