package parser

import (
	"io"
	"reflect"
	"strings"

	"github.com/zishang520/engine.io-go-parser/types"
)

// Returns true if obj is a Buffer or a File.
func IsBinary(data any) bool {
	switch data.(type) {
	case *types.StringBuffer: // false
	case *strings.Reader: // false
	case []byte:
		return true
	case io.Reader:
		return true
	}
	return false
}

func HasBinary(data any) bool {
	switch o := data.(type) {
	case nil:
		return false
	case []any:
		for _, v := range o {
			if HasBinary(v) {
				return true
			}
		}
		return false
	case map[string]any:
		for _, v := range o {
			if HasBinary(v) {
				return true
			}
		}
		return false
	}
	dv := reflect.ValueOf(data)
	if dv.Kind() == reflect.Struct {
		for fi := range dv.NumField() {
			dfv := dv.Field(fi)
			if dfv.CanInterface() && HasBinary(dfv.Interface()) {
				return true
			}
		}
		return false
	}

	return IsBinary(data)
}
