package main

import(
	"encoding/xml"
	"io/ioutil"
	"os"
)

func getVersionFromNuspec(file string) (string, error) {
	xmlFile, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer xmlFile.Close()

	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return "", err
	}
	var nuspec Nuspec
	err = xml.Unmarshal(byteValue, &nuspec)
	if err != nil {
		return "", err
	}
	return nuspec.Metadata.Version, nil
}