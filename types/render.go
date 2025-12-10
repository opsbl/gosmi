package types

//go:generate enumer -type=Render -trimprefix=Render -json -yaml -output=render_string.go

type Render int

const (
	RenderNumeric Render = 1 << iota
	RenderName
	RenderQualified
	RenderFormat
	RenderPrintable
	RenderUnknown
	RenderAll Render = 0xff
)
