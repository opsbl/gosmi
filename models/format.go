package models

//go:generate enumer -type=FormatKind -trimprefix=FormatKind -json -yaml -output=format_string.go

import (
	"fmt"
	"strconv"
	"time"

	"github.com/opsbl/gosmi/types"
)

type FormatKind byte

const (
	FormatKindNone     FormatKind = 0
	FormatKindEnumName FormatKind = 1 << iota
	FormatKindEnumValue
	FormatKindBits
	FormatKindString
	FormatKindUnits
	FormatKindDurationShort
	FormatKindAll FormatKind = 0xff & ^FormatKindUnits
)

func ResolveFormat(formats []FormatKind, defaultFormat ...FormatKind) (format FormatKind) {
	if len(formats) == 0 {
		if len(defaultFormat) == 0 {
			return FormatKindAll
		}
		return defaultFormat[0]
	}
	for _, f := range formats {
		format |= f
	}
	return
}

type Value struct {
	Format    FormatKind
	Formatted string
	Raw       interface{}
}

func (v Value) Bytes() []byte {
	if b, ok := v.Raw.([]byte); ok {
		return b
	}
	if s, ok := v.Raw.(string); ok {
		return []byte(s)
	}
	return []byte{}
}

func (v Value) Duration() time.Duration {
	if d, ok := v.Raw.(time.Duration); ok {
		return d
	}
	return time.Duration(0)
}

func (v Value) Int64() int64 {
	if i, ok := v.Raw.(int64); ok {
		return i
	}
	return 0
}

func (v Value) Uint64() uint64 {
	if i, ok := v.Raw.(int64); ok {
		return uint64(i)
	}
	return 0
}

func (v Value) String() string {
	if v.Format != FormatKindNone {
		return v.Formatted
	}
	if v.Raw == nil {
		return ""
	}
	switch r := v.Raw.(type) {
	case string:
		return r
	case []byte:
		return string(r)
	}
	return fmt.Sprintf("%v", v.Raw)
}

func ToInt64(value interface{}) (val int64, err error) {
	switch value := value.(type) {
	case int64:
		val = value
	case uint64:
		val = int64(value)
	case int:
		val = int64(value)
	case int8:
		val = int64(value)
	case int16:
		val = int64(value)
	case int32:
		val = int64(value)
	case uint:
		val = int64(value)
	case uint8:
		val = int64(value)
	case uint16:
		val = int64(value)
	case uint32:
		val = int64(value)
	case types.SmiSubId:
		val = int64(value)
	case string:
		return strconv.ParseInt(value, 10, 64)
	default:
		err = fmt.Errorf("Value has invalid type: %T", value)
	}
	return
}

type ValueFormatter func(interface{}) Value

func (n Node) FormatValue(value interface{}, flags ...FormatKind) Value {
	return n.Type.FormatValue(value, flags...)
}

func (n Node) GetValueFormatter(flags ...FormatKind) ValueFormatter {
	return n.Type.GetValueFormatter(flags...)
}

func (n ScalarNode) FormatValue(value interface{}, flags ...FormatKind) Value {
	return n.Type.FormatValue(value, flags...)
}

func (n ScalarNode) GetValueFormatter(flags ...FormatKind) ValueFormatter {
	return n.Type.GetValueFormatter(flags...)
}

func (n ColumnNode) FormatValue(value interface{}, flags ...FormatKind) Value {
	return n.Type.FormatValue(value, flags...)
}

func (n ColumnNode) GetValueFormatter(flags ...FormatKind) ValueFormatter {
	return n.Type.GetValueFormatter(flags...)
}

func (t Type) FormatValue(value interface{}, flags ...FormatKind) Value {
	formatFlags := ResolveFormat(flags)
	switch t.BaseType {
	case types.BaseTypeOctetString:
		switch t.Name {
		case "IpAddress", "InetAddress", "IpV4orV6Addr":
			return GetInetAddressFormatted(value, formatFlags)
		}
		return GetOctetStringFormatted(value, formatFlags, t.Format)
	case types.BaseTypeBits:
		if t.Enum == nil {
			return GetBitsFormatted(value, formatFlags)
		}
		return GetEnumBitsFormatted(value, formatFlags, t.Enum)
	case types.BaseTypeEnum:
		return GetEnumFormatted(value, formatFlags, t.Enum)
	}
	switch t.Name {
	case "TimeTicks", "TimeInterval", "TimeStamp":
		return GetDurationFormatted(value, formatFlags)
	}
	return GetIntFormatted(value, formatFlags, t.Format)
}

func (t Type) GetValueFormatter(flags ...FormatKind) ValueFormatter {
	formatFlags := ResolveFormat(flags)
	switch t.BaseType {
	case types.BaseTypeOctetString:
		switch t.Name {
		case "IpAddress", "InetAddress", "IpV4orV6Addr":
			return GetInetAddressFormatter(formatFlags)
		}
		return GetOctetStringFormatter(formatFlags, t.Format)
	case types.BaseTypeBits:
		if t.Enum == nil {
			return GetBitsFormatter(formatFlags)
		}
		return GetEnumBitsFormatter(formatFlags, t.Enum)
	case types.BaseTypeEnum:
		return GetEnumFormatter(formatFlags, t.Enum)
	}
	switch t.Name {
	case "TimeTicks", "TimeInterval", "TimeStamp":
		return GetDurationFormatter(formatFlags)
	}
	return GetIntFormatter(formatFlags, t.Format)
}
