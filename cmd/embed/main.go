package main

import (
	"embed"
	"encoding/json"
	"flag"
	"os"

	"github.com/opsbl/gosmi"
)

//go:embed FIZBIN-MIB.mib
var fs embed.FS

func main() {
	module := flag.String("m", "FIZBIN-MIB", "Module to load")
	flag.Parse()

	api := gosmi.Must()
	defer api.Close()

	api.SetFS(gosmi.NewNamedFS("Embed Example", fs))

	m, err := api.GetModule(*module)
	if err != nil {
		panic(err)
	}

	nodes := m.GetNodes()
	types := m.GetTypes()

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(struct {
		Module gosmi.SmiModule
		Nodes  []gosmi.SmiNode
		Types  []gosmi.SmiType
	}{
		Module: m,
		Nodes:  nodes,
		Types:  types,
	})
}
