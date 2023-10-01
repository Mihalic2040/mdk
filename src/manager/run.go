package manager

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"log"

	"github.com/Mihalic2040/mdk/src/parsers/minecraft"
)

type Launcher struct {
	project string
	assets  string
	libs    string
	game    string
	cfg     Config
}

func (L *Launcher) PrepareLibs() string {
	// Get a list of all files in the folder
	var foder_path string = L.project + L.game + L.libs
	var filePaths []string

	// Walk through the folder and its subfolders
	err := filepath.WalkDir(foder_path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Skip directories
		if d.IsDir() {
			return nil
		}
		// Add file path to the slice
		filePaths = append(filePaths, path)
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	//Add client to path
	jar_path := L.project + L.game + "/" + L.cfg.Version + ".jar"
	filePaths = append(filePaths, jar_path)

	// Format the file paths
	formattedFiles := strings.Join(filePaths, ":")

	// Print the formatted file paths
	//fmt.Println(formattedFiles)
	return formattedFiles
}

func Run(args []string) {
	//fmt.Println(args)
	if args[0] == "cli" || args[0] == "client" {
		wd, _ := os.Getwd()
		lnc := Launcher{
			project: wd,
			assets:  "/assets/",
			libs:    "/libraries/",
			game:    "/client",
			cfg:     *ReadDesirialize(),
		}
		if len(args) > 1 {
			//fmt.Println("Username:", args[1])
			//run_client(args[1])
			lnc.run_client(args[1])
		} else {
			lnc.run_client("dev")
		}
	}
}

func (L *Launcher) run_client(username string) {
	cfg := L.cfg
	minecraft := minecraft.ReadAll(cfg.Version)

	//fmt.Println(minecraft.Arguments)

	//SETUP ARGS
	libs := L.PrepareLibs()
	natives := L.project + "/client/natives"
	game := L.game
	assets := L.assets
	currentDir := L.project

	// Convert []interface{} to []string
	var filteredArgs []interface{}
	for _, v := range minecraft.Arguments.Game {
		switch v.(type) {
		case string, int, float64:
			filteredArgs = append(filteredArgs, v)
		}
	}

	// Convert filteredArgs to string
	var result string
	for _, v := range filteredArgs {
		switch value := v.(type) {
		case int:
			result += strconv.Itoa(value) + " "
		case float64:
			result += strconv.FormatFloat(value, 'f', -1, 64) + " "
		case string:
			result += value + " "
		}
	}

	// Remove the trailing space, if any
	result = result[:len(result)-1]

	cmdArgs := []string{
		"-Dminecraft.client.jar=" + L.project + L.game + "/" + L.cfg.Version + ".jar",
		"-Xmx6048m",
		"-Xms2048m",
		"-Djava.library.path=" + natives,
		"-Dminecraft.launcher.brand=mdk",
		"-Dminecraft.launcher.version=0.0.0",
		"-cp",
		libs,
		minecraft.MainClass,
		"--version", cfg.Version,
		"--accessToken", "c6c43b32bef3c48a644fe1d4c106c17",
		"--username", username,
		"--gameDir", currentDir + game,
		"--assetsDir", currentDir + game + assets,
		"--assetIndex", minecraft.AssetIndex.ID,
		"--userType", "Local",
		"--versionType", "release",
		"--uuid", "fpsdfosdkfsfmksdfs0-vsd9vsd9d",
	}

	fmt.Println(cmdArgs)

	gameProcess := exec.Command("java", cmdArgs...)
	gameProcess.Dir = L.project + game
	gameProcess.Stdout = os.Stdout
	gameProcess.Stderr = os.Stderr

	log.Println("Starting Minecraft ", cfg.Version, "...")
	err := gameProcess.Run()
	if err != nil {
		log.Println(err)
	}

	log.Println("Minecraft ", cfg.Version, " closed gracefully")
}
