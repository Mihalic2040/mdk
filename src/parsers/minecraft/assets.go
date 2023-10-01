package minecraft

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type Object struct {
	Hash string `json:"hash"`
	Size int    `json:"size"`
}

func DownloadAssets(ver string) {
	downloadAssetsIndex(ver)

	min := ReadAll(ver)

	file, err := os.ReadFile("./client/assets/indexes/" + min.AssetIndex.ID + ".json")
	if err != nil {
		fmt.Println(err)
		return
	}
	var data map[string]map[string]Object
	if err := json.Unmarshal([]byte(file), &data); err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Extract and print the names
	// Create a buffered channel with a capacity of 100 to limit concurrent goroutines
	sem := make(chan struct{}, 100)
	var wg sync.WaitGroup
	for key, obj := range data["objects"] {
		// Block until there is space in the semaphore
		sem <- struct{}{}

		wg.Add(1)
		go func(key string, obj Object) {
			defer func() {
				// Release the semaphore after the goroutine completes
				<-sem
				wg.Done()
			}()

			// Implement your download logic here
			downloadAsset("./client/assets/objects/"+obj.Hash[:2]+"/"+obj.Hash, getAssetURL(obj.Hash))
			// ...
		}(key, obj)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("Assets DONE")
}

func downloadAsset(path string, url string) {
	fmt.Println("Downloading asset:", path, url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Downloading", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Server:", err)
		return
	}

	// Create directories if they don't exist
	err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		fmt.Println("Make dir:", err)
		return
	}

	// Create the local file
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Create:", err)
		return
	}
	defer file.Close()

	// Copy the content from the HTTP response to the local file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Write:", err)
		return
	}

	fmt.Println("Download successful:", path)
	return
}

func getAssetURL(hash string) string {
	return "https://resources.download.minecraft.net/" + hash[:2] + "/" + hash
}

func downloadAssetsIndex(ver string) {

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

	resp, err := http.Get(data.AssetIndex.URL)
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	filePath := "./client/assets/indexes/" + data.AssetIndex.ID + ".json"

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
}
