package main

import (
	"encoding/xml"
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
	"sync"
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
	fmt.Println("Loading config...")
	viper.SetConfigName("config")
	// Set the path to look for the config file
	viper.AddConfigPath(".")
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
				outputCh <- fmt.Sprintf(redColor + "%s | %s | unknown?", packageName, packageVersion)
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
		fmt.Println(resetColor + "To install updates, run the following")
		fmt.Println(installScript.String())
	}
}