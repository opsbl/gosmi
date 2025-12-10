package types

//go:generate enumer -type=Access -trimprefix=Access -json -yaml -output=access_string.go

type Access int

const (
	AccessUnknown Access = iota
	AccessNotImplemented
	AccessNotAccessible
	AccessNotify
	AccessReadOnly
	AccessReadWrite
	AccessInstall
	AccessInstallNotify
	AccessReportOnly
	AccessEventOnly
)
