package types

//go:generate enumer -type=Language -trimprefix=Language -json -yaml -output=language_string.go
type Language int

const (
	LanguageUnknown Language = iota
	LanguageSMIv1
	LanguageSMIv2
	LanguageSMIng
	LanguageSPPI
)
