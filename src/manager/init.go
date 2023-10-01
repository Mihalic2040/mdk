package manager

import (
	"fmt"

	"github.com/Mihalic2040/mdk/src/parsers/minecraft"
)

func Init() {
	//DOWNLOAD MINECRAFT JSON
	//INSTALL CLIENT
	//INSTALL MODLOADER ON CLINET
	setup_client()

	//INSTALL SERVER

	//CREATE PKG DB
}

func setup_client() {
	cfg := ReadDesirialize()

	//DOWNLAOD VERSION JSON
	err := minecraft.DownloadJSON(cfg.Version)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	//DOWNLOAD LIBS
	minecraft.DownloadLibs(cfg.Version)
	//DOWNLOAD ASSETS
	minecraft.DownloadAssets(cfg.Version)
	//DOWNLOAD JAR
	minecraft.DownloadClient(cfg.Version)

}

func setup_server() {

}
