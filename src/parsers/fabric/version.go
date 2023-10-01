package fabric

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ModLoaderResponse struct {
	Loader struct {
		Version string `json:"version"`
	} `json:"loader"`
}

var MetaUrl = "https://meta.fabricmc.net"

func GetVersions(ver string) []string {
	url, err := url.JoinPath(MetaUrl, "v2", "versions", "loader", ver)
	if err != nil {
		fmt.Println(err)
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var modLoaderResponses []ModLoaderResponse

	// Unmarshal the JSON data into the slice of ModLoaderResponse
	err = json.Unmarshal(data, &modLoaderResponses)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Extract loader versions into a separate array
	var loaderVersions []string
	for _, response := range modLoaderResponses {
		loaderVersions = append(loaderVersions, response.Loader.Version)
	}

	// Print the extracted loader versions
	return loaderVersions
}
