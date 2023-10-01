package minecraft

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadClient(ver string) {
	file, err := os.ReadFile("./client/" + ver + ".json")
	if err != nil {
		fmt.Println(err)
		return
	}
	var data MinecraftJson
	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Downloading:", data.Downloads.Client.URL)
	resp, err := http.Get(data.Downloads.Client.URL)
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	filePath := "./client/" + ver + ".jar"

	// Ensure the directory exists, create it if necessary
	dir := filepath.Dir(filePath)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	// Write the bytes to the file
	err = os.WriteFile(filePath, jsonData, os.ModePerm)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Client DONE")
	return
}
