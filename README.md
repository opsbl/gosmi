# gosmi

This project is modified from https://github.com/sleepinggenius2/gosmi

Starting with v0.2.0, this library is native Go and no longer a wrapper around libsmi. The implementation is currently very close, but may change in the future.

For the native implementation, two additional components have been added:

* SMIv1/2 parser in [parser](parser)
* libsmi-compatible Go implementation in [smi](smi)

## Usage

### Examples

Examples can now be found in:

* [cmd/parse](cmd/parse)
* [cmd/smi](cmd/smi)
* [cmd/embed](cmd/embed)

### Singleton

```go
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

```

### Multi-instance

```go
package main

import (
	"fmt"
	"github.com/opsbl/gosmi"
	"github.com/opsbl/gosmi/smi"
	"os"
)

func main() {
	smiHandle := smi.NewSmiHandle()
	smiHandle.SetPath("/var/lib/mibs/ietf/")
	helper := gosmi.NewHandleHelper(smiHandle)

	module, err := helper.LoadModule("IF-MIB")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Module: %v\n", module)
	loadedModules := helper.GetLoadedModules()
	fmt.Println("Loaded modules:")
	for _, module := range loadedModules {
		fmt.Println("  ", module.Name)
	}
}

```