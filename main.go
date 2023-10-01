package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Mihalic2040/mdk/src/manager"
	"github.com/Mihalic2040/mdk/src/packages"
)

func main() {
	//Define custom flag set
	fs := flag.NewFlagSet("mdk", flag.ExitOnError)

	// Parse the command-line arguments
	fs.Parse(os.Args[1:])

	// Check the number of arguments
	if fs.NArg() < 1 {
		fmt.Println("Error: Module name is required.")
		help()
		return
	}

	// Get the module name and its arguments
	module := fs.Arg(0)
	moduleArguments := fs.Args()[1:]

	// Perform actions based on the module name
	switch module {
	case "new_project":
		manager.InitNew()
	case "init":
		manager.Init()
	case "run":
		manager.Run(moduleArguments)
	case "get":
		packages.Get(moduleArguments)
	default:
		fmt.Printf("Error: Unknown module '%s'\n", module)
		help()
	}
}

func help() {
	fmt.Println("BASE")
	fmt.Println("new_project - setup project json")
	fmt.Println("init - download and install dependencies on project")
	fmt.Println("run [[client || server || cli || srv]] [[username]] - launch instances, if client you can specify usename.")

	fmt.Println("PACKAGE MANAGER")
	fmt.Println("get [[name]] [[name2]] - install pkg into project")
}
