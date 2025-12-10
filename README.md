# gosmi (forked)

Forked from [github.com/sleepinggenius2/gosmi](https://github.com/sleepinggenius2/gosmi). The upstream project is a native Go reimplementation of libsmi and provides:

* SMIv1/2 parser in [parser](parser)
* libsmi-compatible Go implementation in [smi](smi)

## changes

* Added `smi.Instance` so multiple SMI contexts can run in parallel without sharing module state. The existing package-level functions now operate on a default instance.
* Added circular-import protection during `LoadModule`, returning a clear error instead of recursing forever.
* Added unit tests with bundled MIB fixtures (copied from `/var/lib/mibs/ietf`) covering instance isolation and circular-import detection.

## Usage

```go
// create an instance (required for all operations)
inst, err := smi.NewInstance("my-app")
if err != nil {
    panic(err)
}
inst.SetPath("/path/to/mibs")

// load a module with circular-import protection
if _, err := inst.LoadModuleWithError("RFC1155-SMI"); err != nil {
    panic(err)
}

// gosmi API requires an instance too
g, _ := gosmi.New("app")
g.SetPath("/path/to/mibs")
name, _ := g.LoadModule("RFC1155-SMI")
mod, _ := g.GetModule(name)
node, _ := g.GetNode("sysDescr", mod)
fmt.Println(node.RenderQualified())
```

### Examples

Examples from upstream remain available:

* [cmd/parse](cmd/parse)
* [cmd/smi](cmd/smi)
* [cmd/embed](cmd/embed)
