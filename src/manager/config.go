package manager

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	ModLoader    string `json:"modLoader"`
	ModLoaderVer string `json:"modLoaderVer"`
}

func Serialize(cfg Config) ([]byte, error) {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ReadDesirialize() *Config {
	// Read JSON file
	file, err := os.Open("./project.json")
	if err != nil {
		fmt.Println("ERROR:", err)
		fmt.Println("Maybe you are not in project folder.")
		os.Exit(1)
	}
	defer file.Close()

	// Read file content
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("error reading file: %v", err)
		os.Exit(1)
	}

	// Unmarshal JSON data into struct
	var result Config
	err = json.Unmarshal(data, &result)
	if err != nil {
		fmt.Println("error unmarshaling JSON: %v", err)
		os.Exit(1)
	}

	return &result
}
