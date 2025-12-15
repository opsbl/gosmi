package models

import (
	"fmt"
	"strings"
)

func GetBitsFormatted(value interface{}, flags FormatKind) (v Value) {
	v.Format = flags
	v.Raw = value
	if flags&FormatKindBits != 0 {
		if bytes, ok := value.([]byte); ok {
			v.Formatted = fmt.Sprintf("% X", bytes)
		}
	}
	return
}

func GetBitsFormatter(flags FormatKind) (f ValueFormatter) {
	return func(value interface{}) Value {
		return GetBitsFormatted(value, flags)
	}
}

func GetEnumBitsFormatted(value interface{}, flags FormatKind, enum *Enum) (v Value) {
	v.Format = flags
	v.Raw = value
	if flags == FormatKindNone {
		return
	}
	octets := value.([]byte)
	if flags&FormatKindBits != 0 {
		v.Formatted = fmt.Sprintf("% X", octets)
	}
	if (flags&FormatKindEnumName)+(flags&FormatKindEnumValue) == 0 {
		return
	}
	bitsFormatted := make([]string, 0, 8*len(octets))
	for i, octet := range octets {
		for j := 7; j >= 0; j-- {
			if octet&(1<<uint(j)) != 0 {
				bit := uint64(8*i + (7 - j))
				var bitFormatted string
				if flags&FormatKindEnumName != 0 {
					bitFormatted = enum.Name(int64(bit))
					if flags&FormatKindEnumValue != 0 || bitFormatted == "unknown" {
						bitFormatted += "(" + fmt.Sprintf("%d", bit) + ")"
					}
				} else if flags&FormatKindEnumValue != 0 {
					bitFormatted = fmt.Sprintf("%d", bit)
				}
				bitsFormatted = append(bitsFormatted, bitFormatted)
			}
		}
	}
	if v.Formatted == "" {
		v.Formatted = strings.Join(bitsFormatted, " ")
	} else {
		v.Formatted += "[" + strings.Join(bitsFormatted, " ") + "]"
	}
	return
}

func GetEnumBitsFormatter(flags FormatKind, enum *Enum) (f ValueFormatter) {
	return func(value interface{}) Value {
		return GetEnumBitsFormatted(value, flags, enum)
	}
}
