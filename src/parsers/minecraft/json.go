package minecraft

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type VersionData struct {
	Latest   LatestVersion `json:"latest"`
	Versions []VersionInfo `json:"versions"`
}

type LatestVersion struct {
	Release  string `json:"release"`
	Snapshot string `json:"snapshot"`
}

type VersionInfo struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	Time        string `json:"time"`
	ReleaseTime string `json:"releaseTime"`
}

type MinecraftJson struct {
	Arguments struct {
		Game []any `json:"game"`
		Jvm  []any `json:"jvm"`
	} `json:"arguments"`
	AssetIndex struct {
		ID        string `json:"id"`
		Sha1      string `json:"sha1"`
		Size      int    `json:"size"`
		TotalSize int    `json:"totalSize"`
		URL       string `json:"url"`
	} `json:"assetIndex"`
	Assets          string `json:"assets"`
	ComplianceLevel int    `json:"complianceLevel"`
	Downloads       struct {
		Client struct {
			Sha1 string `json:"sha1"`
			Size int    `json:"size"`
			URL  string `json:"url"`
		} `json:"client"`
		ClientMappings struct {
			Sha1 string `json:"sha1"`
			Size int    `json:"size"`
			URL  string `json:"url"`
		} `json:"client_mappings"`
		Server struct {
			Sha1 string `json:"sha1"`
			Size int    `json:"size"`
			URL  string `json:"url"`
		} `json:"server"`
		ServerMappings struct {
			Sha1 string `json:"sha1"`
			Size int    `json:"size"`
			URL  string `json:"url"`
		} `json:"server_mappings"`
	} `json:"downloads"`
	ID          string `json:"id"`
	JavaVersion struct {
		Component    string `json:"component"`
		MajorVersion int    `json:"majorVersion"`
	} `json:"javaVersion"`
	Libraries []Library `json:"libraries"`
	Logging   struct {
		Client struct {
			Argument string `json:"argument"`
			File     struct {
				ID   string `json:"id"`
				Sha1 string `json:"sha1"`
				Size int    `json:"size"`
				URL  string `json:"url"`
			} `json:"file"`
			Type string `json:"type"`
		} `json:"client"`
	} `json:"logging"`
	MainClass              string    `json:"mainClass"`
	MinimumLauncherVersion int       `json:"minimumLauncherVersion"`
	ReleaseTime            time.Time `json:"releaseTime"`
	Time                   time.Time `json:"time"`
	Type                   string    `json:"type"`
}

type Library struct {
	Downloads struct {
		Artifact struct {
			Path string `json:"path"`
			Sha1 string `json:"sha1"`
			Size int    `json:"size"`
			URL  string `json:"url"`
		} `json:"artifact"`
		Classifiers Classifiers `json:"classifiers"`
	} `json:"downloads"`
	Name  string `json:"name"`
	Rules []struct {
		Action string `json:"action"`
		Os     struct {
			Name string `json:"name"`
		} `json:"os"`
	} `json:"rules,omitempty"`
}

type Classifiers struct {
	NativesOsx struct {
		Path string `json:"path"`
		Sha0 string `json:"sha1"`
		Size int    `json:"size"`
		URL  string `json:"url"`
	} `json:"natives-osx"`
	NativesLinux struct {
		Path string `json:"path"`
		Sha0 string `json:"sha1"`
		Size int    `json:"size"`
		URL  string `json:"url"`
	} `json:"natives-linux"`
	NativesWindows struct {
		Path string `json:"path"`
		Sha0 string `json:"sha1"`
		Size int    `json:"size"`
		URL  string `json:"url"`
	} `json:"natives-windows"`
}

func ReadAll(ver string) *MinecraftJson {
	file, err := os.ReadFile("./client/" + ver + ".json")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var data MinecraftJson
	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &data
}

func DownloadJSON(MinecraftVersion string) error {
	var manifest string = "https://launchermeta.mojang.com/mc/game/version_manifest.json"

	resp, err := http.Get(manifest)
	if err != nil {
		fmt.Println(err)
		return err
	}

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var versionData VersionData
	err = json.Unmarshal([]byte(jsonData), &versionData)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	var json_uri string
	for _, version := range versionData.Versions {
		if version.ID == MinecraftVersion {
			json_uri = version.URL
		}
	}

	resp, err = http.Get(json_uri)
	if err != nil {
		fmt.Println(err)
		return err
	}

	jsonData, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	filePath := "./client/" + MinecraftVersion + ".json"

	// Ensure the directory exists, create it if necessary
	dir := filepath.Dir(filePath)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return err
	}

	// Write the bytes to the file
	err = os.WriteFile(filePath, jsonData, os.ModePerm)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	return nil
}
