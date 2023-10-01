package minecraft

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func DownloadLibs(ver string) {
	//READ MINECRAFT JSON
	file, err := os.ReadFile("./client/" + ver + ".json")
	if err != nil {
		fmt.Println(err)
	}

	var data MinecraftJson
	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println(err)
	}

	var wg sync.WaitGroup
	for _, lib := range data.Libraries {
		if checkRules(lib.Rules, runtime.GOOS) {
			//fmt.Println(lib.Name, lib.Rules)
			// Implement download logic here

			wg.Add(1)
			go func(lib Library) {
				// Decrement the counter when the goroutine completes
				defer wg.Done()
				// Implement download logic here for the specific library
				downloadLib(lib)
				//DOWNLADO NATIVE IF AVALIBLE
				if !isEmptyClassifiers(lib.Downloads.Classifiers) {
					if runtime.GOOS == "darwin" && lib.Downloads.Classifiers.NativesOsx.Path != "" {
						downloadNativeLib("./client/natives/"+filepath.Base(lib.Downloads.Classifiers.NativesOsx.Path), lib.Downloads.Classifiers.NativesOsx.URL)
					} else if runtime.GOOS == "linux" && lib.Downloads.Classifiers.NativesLinux.Path != "" {
						downloadNativeLib("./client/natives/"+filepath.Base(lib.Downloads.Classifiers.NativesLinux.Path), lib.Downloads.Classifiers.NativesLinux.URL)
					} else if runtime.GOOS == "windows" && lib.Downloads.Classifiers.NativesWindows.Path != "" {
						downloadNativeLib("./client/natives/"+filepath.Base(lib.Downloads.Classifiers.NativesWindows.Path), lib.Downloads.Classifiers.NativesWindows.URL)
					}
				}

				// ...
			}(lib)
			// ...
		} else {
			//fmt.Println("Library", lib.Name, "cannot be downloaded on", runtime.GOOS)
		}
	}
	wg.Wait()
	fmt.Println("Libs DONE")
}

func isEmptyClassifiers(classifiers Classifiers) bool {
	return classifiers.NativesOsx.Path == "" &&
		classifiers.NativesLinux.Path == "" &&
		classifiers.NativesWindows.Path == ""
}

func checkRules(rules []struct {
	Action string `json:"action"`
	Os     struct {
		Name string `json:"name"`
	} `json:"os"`
}, currentOS string) bool {
	for _, rule := range rules {
		if rule.Action == "allow" && strings.ToLower(rule.Os.Name) == "osx" {
			// If "allow" rule specifies osx and current OS is not macOS, skip download
			if currentOS != "darwin" && currentOS != "osx" {
				return false
			}
		}
		if rule.Action == "disallow" && strings.ToLower(rule.Os.Name) == currentOS {
			return false // Disallow download when "disallow" rule matches current OS
		}
	}
	// If no rules match, allow by default
	return true
}

func downloadLib(lib Library) {
	fmt.Println("Downloading artifact:", lib.Name)
	resp, err := http.Get(lib.Downloads.Artifact.URL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(err)
		return
	}

	// Create directories if they don't exist
	err = os.MkdirAll(filepath.Dir("./client/libraries/"+lib.Downloads.Artifact.Path), os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create the local file
	file, err := os.Create("./client/libraries/" + lib.Downloads.Artifact.Path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Copy the content from the HTTP response to the local file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(lib.Name + "Download successful!")
	return
}

func downloadNativeLib(filePath, url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to download natives library:", resp.Status)
		return
	}

	err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Unzip the library file into the natives folder
	err = unzip(filePath, "./client/natives")
	if err != nil {
		fmt.Println("Error unzipping library:", err)
		return
	}

	// Delete the original downloaded file
	err = os.Remove(filePath)
	if err != nil {
		fmt.Println("Error deleting downloaded file:", err)
		return
	}

	fmt.Println("Native library downloaded successfully:", filePath)
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		filePath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
		} else {
			os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
			w, err := os.Create(filePath)
			if err != nil {
				return err
			}
			_, err = io.Copy(w, rc)
			if err != nil {
				w.Close()
				return err
			}
			w.Close()
		}
	}

	return nil
}
