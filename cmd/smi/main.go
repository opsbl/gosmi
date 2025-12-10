package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/opsbl/gosmi"
	"github.com/opsbl/gosmi/types"
)

type arrayStrings []string

var modules arrayStrings
var paths arrayStrings
var debug bool
var api *gosmi.Instance

func (a arrayStrings) String() string {
	return strings.Join(a, ",")
}

func (a *arrayStrings) Set(value string) error {
	*a = append(*a, value)
	return nil
}

func main() {
	flag.BoolVar(&debug, "d", false, "Debug")
	flag.Var(&modules, "m", "Module to load")
	flag.Var(&paths, "p", "Path to add")
	flag.Parse()

	Init()

	oid := flag.Arg(0)
	if oid == "" {
		ModuleTrees()
	} else {
		Subtree(oid)
	}

	Exit()
}

func Init() {
	api = gosmi.Must()

	for _, path := range paths {
		api.AppendPath(path)
	}

	for i, module := range modules {
		moduleName, err := api.LoadModule(module)
		if err != nil {
			fmt.Printf("Init Error: %s\n", err)
			return
		}
		if debug {
			fmt.Printf("Loaded module %s\n", moduleName)
		}
		modules[i] = moduleName
	}

	if debug {
		path := api.GetPath()
		fmt.Printf("Search path: %s\n", path)
		loadedModules := api.GetLoadedModules()
		fmt.Println("Loaded modules:")
		for _, loadedModule := range loadedModules {
			fmt.Printf("  %s (%s)\n", loadedModule.Name, loadedModule.Path)
		}
	}
}

func Exit() {
	if api != nil {
		api.Close()
	}
}

func Subtree(oid string) {
	var node gosmi.SmiNode
	var err error
	if (oid[0] >= '0' && oid[0] <= '9') || oid[0] == '.' {
		node, err = api.GetNodeByOID(types.OidMustFromString(oid))
	} else {
		node, err = api.GetNode(oid)
	}
	if err != nil {
		fmt.Printf("Subtree Error: %s\n", err)
		return
	}

	subtree := node.GetSubtree()

	jsonBytes, _ := json.Marshal(subtree)
	os.Stdout.Write(jsonBytes)
}

func ModuleTrees() {
	for _, module := range modules {
		m, err := api.GetModule(module)
		if err != nil {
			fmt.Printf("ModuleTrees Error: %s\n", err)
			continue
		}

		nodes := m.GetNodes()
		types := m.GetTypes()

		jsonBytes, _ := json.Marshal(struct {
			Module gosmi.SmiModule
			Nodes  []gosmi.SmiNode
			Types  []gosmi.SmiType
		}{
			Module: m,
			Nodes:  nodes,
			Types:  types,
		})
		os.Stdout.Write(jsonBytes)
	}
}
