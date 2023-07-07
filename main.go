package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"

	"github.com/atotto/clipboard"
	"github.com/spf13/viper"
)

const (
	redColor    = "\033[31m"
	greenColor  = "\033[32m"
	yellowColor = "\033[33m"
	resetColor  = "\033[0m"
)

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Entry   Entry    `xml:"entry"`
}

type Entry struct {
	XMLName    xml.Name   `xml:"entry"`
	Properties Properties `xml:"properties"`
}

type Properties struct {
	XMLName xml.Name `xml:"properties"`
	Version string   `xml:"Version"`
}

type Nuspec struct {
	XMLName  xml.Name `xml:"package"`
	Metadata struct {
		XMLName xml.Name `xml:"metadata"`
		Version string   `xml:"version"`
	} `xml:"metadata"`
}

func main() {

	if isAdmin() {
		fmt.Println("Admin mode is turned on.")
	}

	fmt.Println("Loading config...")
	viper.SetConfigName("config")
	// Set the path to look for the config file
	viper.AddConfigPath(getPath())
	// Read the config file
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
		return
	}

	// Get the URL from the config file
	urls := viper.GetStringSlice("urls")

	fmt.Println("Loading packages and checking")
	fmt.Println("Package Name | current version | available version ")
	nuspecFiles, err := filepath.Glob(viper.GetString("baseDirectory"))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	hasUpdates := false
	var installScript strings.Builder
	installScript.WriteString("choco upgrade -y ")

	// Channel to store the output messages
	outputCh := make(chan string)
	var results sync.Map
	var wg sync.WaitGroup

	for _, file := range nuspecFiles {

		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			packageVersion, err := getVersionFromNuspec(file)
			if err != nil {
				outputCh <- fmt.Sprintf("Error: %v", err)
				return
			}

			packageName := filepath.Base(filepath.Dir(file))
			latestVersion, err := getLatestVersionFromAPI(packageName, urls)
			if err != nil {
				outputCh <- fmt.Sprintf("Error: %v", err)
				return
			}
			if latestVersion == "" {
				outputCh <- fmt.Sprintf(redColor+"%s | %s | unknown?", packageName, packageVersion)
			} else if packageVersion != latestVersion {
				outputCh <- fmt.Sprintf(yellowColor+"%s |%s | %s", packageName, packageVersion, latestVersion)
				results.Store(packageName, latestVersion)
			} else {
				outputCh <- fmt.Sprintf(greenColor+"%s | %s | %s", packageName, packageVersion, latestVersion)
			}
		}(file)
	}
	go func() {
		wg.Wait()
		close(outputCh)
	}()

	for msg := range outputCh {
		fmt.Println(msg)
	}

	results.Range(func(key, value interface{}) bool {
		hasUpdates = true
		installScript.WriteString(" ")
		installScript.WriteString(fmt.Sprintf("%s", key))
		installScript.WriteString(" ")
		return true // continue iterating
	})

	if hasUpdates {
		args := os.Args[1:]
		if len(args) > 0 && args[0] == "auto" {
			handleAutoMode(installScript.String())
		} else {
			fmt.Println(resetColor + "To install updates, run the following")
			fmt.Println(installScript.String())
		}
	}
}

func getPath() string {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("Failed to get executable path: %v\n", err)
		return "."
	}

	// Get the directory containing the executable
	appDir := filepath.Dir(exePath)
	return appDir
}

func isAdmin() bool {
	if runtime.GOOS == "windows" {
		_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
		return err == nil
	}
	return os.Getuid() == 0
}

func runCommand(command string) error {
	cmd := exec.Command("cmd", "/C", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

func handleAutoMode(path string) {
	if isAdmin() {
		fmt.Println(resetColor + "calling choco to install updates - " + path)
		err := runCommand(path)
		if err != nil {
			fmt.Println("Error running 'choco upgrade':", err)
		}
	} else {
		fmt.Println("An admin terminal needs to be opened for the command to run.")
		promptClipboard(path)
	}
}

func promptClipboard(path string) {
	fmt.Println("Would you like me to enter it into your clipboard? (y/n)")
	var response string
	fmt.Scanln(&response)
	response = strings.TrimSpace(strings.ToLower(response))
	if response == "y" || response == "yes" {
		err := clipboard.WriteAll(path)
		if err != nil {
			fmt.Println("Error writing to clipboard:", err)
		} else {
			fmt.Println("The command 'choco upgrade' has been copied to your clipboard.")
		}
	}
}
