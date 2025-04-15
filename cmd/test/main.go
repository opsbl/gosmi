package main

import (
	"fmt"
	"github.com/opsbl/gosmi"
	"os"
)

func main() {
	gosmi.Init()

	gosmi.SetPath("/var/lib/mibs/ietf/")
	module, err := gosmi.LoadModule("IF-MIB")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Module: %v\n", module)
	loadedModules := gosmi.GetLoadedModules()
	fmt.Println("Loaded modules:")
	for _, module := range loadedModules {
		fmt.Println("  ", module.Name)
	}
}
