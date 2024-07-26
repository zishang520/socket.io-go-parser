package parser

import (
	"io"
	"reflect"
	"strings"

	"github.com/zishang520/engine.io-go-parser/types"
)

// IsBinary returns true if the data is a binary type (Buffer or File).
func IsBinary(data any) bool {
	switch data.(type) {
	case *types.StringBuffer, *strings.Reader:
		return false
	case []byte, io.Reader:
		return true
	default:
		return false
	}
}

// HasBinary checks recursively if the data contains any binary data.
func HasBinary(data any) bool {
	switch v := data.(type) {
	case nil:
		return false
	case []any:
		for _, item := range v {
			if HasBinary(item) {
				return true
			}
		}
		return false
	case map[string]any:
		for _, value := range v {
			if HasBinary(value) {
				return true
			}
		}
		return false
	}

	if IsBinary(data) {
		return true
	}

	dv := reflect.ValueOf(data)
	switch dv.Kind() {
	case reflect.Pointer:
		return IsBinary(dv.Elem().Interface())
	case reflect.Struct:
		for fi := range dv.NumField() {
			dfv := dv.Field(fi)
			if dfv.CanInterface() && HasBinary(dfv.Interface()) {
				return true
			}
		}
		return false
	case reflect.Array, reflect.Slice:
		for i := range dv.Len() {
			av := dv.Index(i)
			if av.CanInterface() && HasBinary(av.Interface()) {
				return true
			}
		}
		return false
	case reflect.Map:
		mr := dv.MapRange()
		for mr.Next() {
			// Keys can't be binary blobs in json, so only check values
			mv := mr.Value()
			if mv.CanInterface() && HasBinary(mv.Interface()) {
				return true
			}
		}
		return false
	default:
		return false
	}
}
