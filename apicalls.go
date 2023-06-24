package main

import(
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func getLatestVersionFromAPI(packageName string, urls []string) (string, error) {
	var latestVersion string

	for _, url := range urls {
		apiURL := url + "/Packages()?$filter=(tolower(Id)%20eq%20'" + strings.ToLower(packageName) + "')%20and%20IsLatestVersion&semVerLevel=2.0.0"
		resp, err := http.Get(apiURL)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			continue
		}
		var feed Feed
		err = xml.NewDecoder(resp.Body).Decode(&feed)
		if err != nil {
			return "", err
		}
		version := feed.Entry.Properties.Version
		if version != "" && (latestVersion == "" || compareVersions(version, latestVersion) > 0) {
			latestVersion = version
		}
	}
	if latestVersion == "" {
		return "", fmt.Errorf("Unable to find the latest version from the URLs")
	}
	return latestVersion, nil
}

func compareVersions(version1, version2 string) int {
	// Split the versions into individual components
	components1 := strings.Split(version1, ".")
	components2 := strings.Split(version2, ".")

	// Compare each component from left to right
	for i := 0; i < len(components1) && i < len(components2); i++ {
		// Convert the component strings to integers
		num1, _ := strconv.Atoi(components1[i])
		num2, _ := strconv.Atoi(components2[i])

		// Compare the components
		if num1 < num2 {
			return -1
		} else if num1 > num2 {
			return 1
		}
	}
	return len(components1) - len(components2)
}