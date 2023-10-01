package manager

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Mihalic2040/mdk/src/parsers/fabric"
	"github.com/Mihalic2040/mdk/src/parsers/forge/version"
	"github.com/axllent/semver"
)

func InitNew() {
	// CREATE PROJECT DIR
	var projectName string
	fmt.Print("Project name: ")
	fmt.Scanln(&projectName)

	// If the user presses Enter without typing anything, use the example name
	if projectName == "" {
		fmt.Println("Project Name is NULL!!!")
		return
	}

	err := projectDir(projectName)
	if err != nil {
		fmt.Println(err)
	}

	// //INIT DIR STURUCTURE
	// err = initDirStructure(projectName)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// CREATE PROJECT JSON
	var version string
	fmt.Println("\n\n\nDefault: 1.20.1")
	fmt.Print("Minecraft version: ")
	fmt.Scanln(&version)

	// If the user presses Enter without typing anything, use the example version
	if version == "" {
		version = "1.20.1" // Example version
	}

	fmt.Println("\n\n\nExample Mod Loaders: forge, fabric")
	fmt.Println("Default: forge")
	var modLoader string
	fmt.Print("Mod Loader: ")
	reader := bufio.NewReader(os.Stdin)
	modLoader, _ = reader.ReadString('\n')
	modLoader = strings.TrimSpace(modLoader)

	// If the user presses Enter without typing anything, use the example mod loader
	if modLoader == "" {
		modLoader = "forge" // Example mod loader
	}

	modLoaderVer, err := getLoaderNewest(modLoader, version)
	if err != nil {
		fmt.Println(err)
		return
	}

	cfg := Config{
		Name:         projectName,
		Version:      version,
		ModLoader:    modLoader,
		ModLoaderVer: modLoaderVer,
	}

	jsonCfg, err := Serialize(cfg)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Save project configuration to project.json
	file, err := os.Create(filepath.Join(projectName, "project.json"))
	if err != nil {
		fmt.Println("Error creating project.json file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonCfg)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("PROJECT SETUP DONE\n\n\nTo start working type:\ncd " + projectName + "\nmdk prepare\nmdk run client")
}

func projectDir(Name string) error {
	// Create the main project directory
	err := os.Mkdir(Name, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func initDirStructure(projectName string) error {
	// Define subdirectories
	directories := []string{
		"client",
		"server",
	}

	// Create subdirectories inside the project
	for _, dir := range directories {
		err := os.Mkdir(filepath.Join(projectName, dir), os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func getLoaderNewest(loader string, ver string) (string, error) {
	switch loader {
	case "forge":
		forge, err := version.FromDefault()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		f := &version.Filter{
			Minecraft: ver,
		}

		var versions []string
		for _, v := range forge.Releases.Filter(f) {
			versions = append(versions, v.ID)
		}

		// Sort the versions in ascending
		newset := semver.SortMax(versions)

		//fmt.Println("Version", newset)
		// Get the latest version (which is the last element in the sorted slice)
		if len(versions) > 0 {
			//fmt.Println(semver.SortMax(versions))
			fmt.Println("\n\n\nDefault version:", newset[0], "LATEST")
		} else {
			fmt.Println("No matching versions found.")
		}

		var modLoaderV string
		fmt.Print("Mod Loader Version: ")
		reader := bufio.NewReader(os.Stdin)
		modLoaderV, _ = reader.ReadString('\n')
		modLoaderV = strings.TrimSpace(modLoaderV)

		var modLoaderExists bool
		for _, version := range versions {
			if modLoaderV == version {
				modLoaderExists = true
				break
			}
		}

		if modLoaderExists == false {
			//fmt.Println("Version:", modLoaderV, "not valid setting default.")
			modLoaderV = newset[0]
		}
		return modLoaderV, nil

	case "fabric":
		vers := fabric.GetVersions(ver)
		newest := semver.SortMax(vers)

		fmt.Println("\n\n\nDafult version:", newest[0], "LATEST")

		var modLoaderV string
		fmt.Print("Mod Loader Version: ")
		reader := bufio.NewReader(os.Stdin)
		modLoaderV, _ = reader.ReadString('\n')
		modLoaderV = strings.TrimSpace(modLoaderV)

		var modLoaderExists bool
		for _, version := range vers {
			if modLoaderV == version {
				modLoaderExists = true
				break
			}
		}

		if modLoaderExists == false {
			//fmt.Println("Version:", modLoaderV, "not valid setting default.")
			modLoaderV = newest[0]
		}
		return modLoaderV, nil

	}

	return "", nil
}
