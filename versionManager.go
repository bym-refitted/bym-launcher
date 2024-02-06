package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const versionInfoPathBase = "api.bymrefitted.com/launcher.json"
const downloadBasePath = "api.bymrefitted.com/launcher/downloads/"
const downloadsFolder = "bymDownloads"
const buildFolder = "bymDownloads/swfBuilds"
const runtimeFolder = "bymDownloads/flashRuntimes"

type LocalVersionManifest struct {
	CurrentGameVersion     string        `json:"currentGameVersion"`
	CurrentLauncherVersion string        `json:"currentLauncherVersion"`
	Builds                 Builds        `json:"builds"`
	FlashRuntimes          FlashRuntimes `json:"flashRuntimes"`
}

type VersionManifest struct {
	CurrentGameVersion     string        `json:"currentGameVersion"`
	CurrentLauncherVersion string        `json:"currentLauncherVersion"`
	Builds                 Builds        `json:"builds"`
	FlashRuntimes          FlashRuntimes `json:"flashRuntimes"`
	httpsWorked            bool
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

func getVersionInfo(ctx context.Context) (VersionManifest, error) {
	httpsWorked := false
	// First we try https
	resp, err := http.Get("https://" + versionInfoPathBase)
	if err != nil {
		runtime.EventsEmit(ctx, "infoLog", fmt.Sprintf("Could not access over https, attempting http %s", err))
		// try via http if that fails
		resp, err = http.Get("http://" + versionInfoPathBase)
		if err != nil {
			runtime.EventsEmit(ctx, "infoLog", fmt.Sprintf("Could not access over http, please check the server status on our discord %s", err))
			return VersionManifest{}, err
		}
	} else {
		runtime.EventsEmit(ctx, "infoLog", "ONLINE: Launcher successfully connected over https")
		httpsWorked = true
	}
	defer resp.Body.Close()

	// Create a variable to hold the decoded data
	var data VersionManifest

	// Decode the JSON response
	err = json.NewDecoder(resp.Body).Decode(&data)
	data.httpsWorked = httpsWorked
	if err != nil {
		runtime.EventsEmit(ctx, "infoLog", "Failed to decode version info")
		return VersionManifest{}, err
	}
	return data, nil
}

func localFilesStatus() (bool, LocalVersionManifest, error) {
	ensureFolderExists(downloadsFolder)
	ensureFolderExists(buildFolder)
	ensureFolderExists(runtimeFolder)
	// get versions from local json file if it exists
	return getLocalVersions()
}

func downloadSwfs(builds Builds, version string, useHttps bool) error {
	buildsToDownload := map[string]string{
		builds.Stable: "stable",
		builds.Http:   "http",
		builds.Local:  "local",
	}
	for build, buildName := range buildsToDownload {
		buildPath := filepath.Join(buildFolder, fmt.Sprintf("bymr-%s-%s.swf", buildName, version))
		err := downloadFile(buildPath, build, useHttps)
		if err != nil {
			return err
		}
	}

	return nil
}

func doAllSwfsExist(builds Builds, version string) bool {
	buildsToDownload := map[string]string{
		builds.Stable: "stable",
		builds.Http:   "http",
		builds.Local:  "local",
	}
	for _, buildName := range buildsToDownload {
		buildPath := filepath.Join(buildFolder, fmt.Sprintf("bymr-%s-%s.swf", buildName, version))
		if !fileExists(buildPath) {
			return false
		}
	}
	return true
}

func downloadRuntimes(flashRuntimeUrl string, flashRuntimeFileName string, useHttps bool) error {
	flashFilePath := filepath.Join(runtimeFolder, flashRuntimeFileName)
	return downloadFile(flashFilePath, flashRuntimeUrl, useHttps)
}

func getPlatformFlashRuntime(envInfo runtime.EnvironmentInfo, serverManifest VersionManifest) (string, string, error) {
	switch envInfo.Platform {
	case "windows":
		return serverManifest.FlashRuntimes.Windows, "flashplayer.exe", nil
	case "darwin":
		return serverManifest.FlashRuntimes.Darwin, "flashplayer.dmg", nil
	case "linux":
		return serverManifest.FlashRuntimes.Linux, "flashplayer.tar.gz", nil
	}
	return "", "", fmt.Errorf("Unsupported platform: %s", envInfo.Platform)
}
