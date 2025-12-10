package types

//go:generate enumer -type=Status -trimprefix=Status -json -yaml -output=status_string.go

type Status int

const (
	StatusUnknown Status = iota
	StatusCurrent
	StatusDeprecated
	StatusMandatory
	StatusOptional
	StatusObsolete
)
