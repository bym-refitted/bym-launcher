package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// sartup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type InitialInfo struct {
	Platform     string          `json:"platform"`
	Architecture string          `json:"architecture"`
	Manifest     VersionManifest `json:"manifest"`
}

func (a *App) InitializeApp() error {
	// First - get current info
	if a.ctx == nil {
		return errors.New("context is nil, something went severely wrong. Please restart the app and contact GHark on Discord.")
	}
	// Get OS info
	os := runtime.Environment(a.ctx)

	runtime.EventsEmit(a.ctx, "infoLog", fmt.Sprintf("Platform: %s %v", os.Platform, os.Arch))

	serverManifest, err := getVersionInfo(a.ctx)
	if err != nil {
		runtime.EventsEmit(a.ctx, "infoLog", "Server manifest could not be retieved. Please check your internet connection.")
		return err
	}

	// Get server info - is online, latest version, runtime info + links
	runtime.EventsEmit(a.ctx, "initialLoad", InitialInfo{
		Platform:     os.Platform,
		Architecture: os.Arch,
		Manifest:     serverManifest,
	})

	localManifestExists, localManifest, err := localFilesStatus()

	noLocalManifest := !localManifestExists || err != nil

	shouldRefreshBuilds := noLocalManifest || serverManifest.CurrentGameVersion != localManifest.CurrentGameVersion || !doAllSwfsExist(serverManifest.Builds, serverManifest.CurrentGameVersion)

	if shouldRefreshBuilds {
		// download swfs
		runtime.EventsEmit(a.ctx, "infoLog", fmt.Sprintf("Downloading latest SWFs"))
		err := downloadSwfs(serverManifest.Builds, serverManifest.CurrentGameVersion, serverManifest.httpsWorked)

		if err != nil {
			runtime.EventsEmit(a.ctx, "infoLog", fmt.Sprintf("Could not download latest swfs %s", err))
		}
	}

	flashRuntimeUrl, flashRuntimeFileName, err := getPlatformFlashRuntime(os, serverManifest)

	if noLocalManifest || !fileExists(flashRuntimeFileName) {
		// download players
		runtime.EventsEmit(a.ctx, "infoLog", fmt.Sprintf("Downloading flash player"))
		downloadRuntimes(flashRuntimeUrl, flashRuntimeFileName, serverManifest.httpsWorked)
		if err != nil {
			runtime.EventsEmit(a.ctx, "infoLog", fmt.Sprintf("Could not download latest flash runtime %s", err))
		}
	}

	// Store the locally downloaded versions
	setLocalVersions(LocalVersionManifest{
		CurrentGameVersion:     serverManifest.CurrentGameVersion,
		CurrentLauncherVersion: serverManifest.CurrentLauncherVersion,
		Builds:                 serverManifest.Builds,
		FlashRuntimes:          serverManifest.FlashRuntimes,
	})

	// Pass to frontend
	return nil
}

func (a *App) LaunchGame(buildName string, version string, flashRuntimeName string) error {

	swfPath := filepath.Join(".", buildFolder, fmt.Sprintf("bymr-%s-%s.swf", buildName, version))

	if !fileExists(swfPath) {
		fmt.Print("Cannot find file: ", swfPath)
		return errors.New("cannot find swf build")
	}
	// const runtimeFolder = "flashRuntimes"

	flashRuntimePath := filepath.Join(".", runtimeFolder, flashRuntimeName)
	if !fileExists(flashRuntimePath) {
		fmt.Print("cannot find file: ", flashRuntimePath)
		return errors.New("Cannot find flashplayer")
	}
	fmt.Print("Opening: ", flashRuntimePath, swfPath)
	cmd := exec.Command(flashRuntimePath, swfPath)
	// cmd.SysProcAttr = &syscall.SysProcAttr{
	// 	HideWindow:    true,
	// 	CreationFlags: 0x08000000,
	// }
	if err := cmd.Start(); err != nil {
		log.Printf("[BYMR LAUNCHER] Failed to start BYMR build %s: %v", buildName, err)
		return err
	}
	return nil
}
