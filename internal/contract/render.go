package contract

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
)

// Success builds a success envelope for a single resource.
func Success(data any) Response {
	return Response{OK: true, Data: data}
}

// SuccessList builds a success envelope for a collection response.
// It normalizes nil slices to empty slices to ensure JSON output is []
// rather than null.
func SuccessList(data any, meta Meta) Response {
	return Response{OK: true, Data: normalizeSlice(data), Meta: &meta}
}

// Err builds an error envelope.
func Err(code string, message string) Response {
	return Response{OK: false, Error: &ErrorDetail{Code: code, Message: message}}
}

// ErrWithDetail builds an error envelope with additional detail.
func ErrWithDetail(code string, message string, detail any) Response {
	return Response{OK: false, Error: &ErrorDetail{Code: code, Message: message, Detail: detail}}
}

// Write serializes a Response as JSON to w.
// It validates envelope invariants before serializing.
// When pretty is true, output is indented with 2 spaces.
// A trailing newline is always appended.
func Write(w io.Writer, resp Response, pretty bool) error {
	if err := validate(resp); err != nil {
		return err
	}

	var out []byte
	var err error
	if pretty {
		out, err = json.MarshalIndent(resp, "", "  ")
	} else {
		out, err = json.Marshal(resp)
	}
	if err != nil {
		return fmt.Errorf("marshal response: %w", err)
	}
	out = append(out, '\n')
	_, err = w.Write(out)
	return err
}

func validate(resp Response) error {
	if resp.OK && resp.Error != nil {
		return errors.New("contract violation: ok is true but error is set")
	}
	if !resp.OK && resp.Data != nil {
		return errors.New("contract violation: ok is false but data is set")
	}
	if !resp.OK && resp.Error == nil {
		return errors.New("contract violation: ok is false but error is nil")
	}
	return nil
}

// normalizeSlice ensures a nil slice becomes an empty slice so that
// json.Marshal produces [] instead of null.
func normalizeSlice(data any) any {
	if data == nil {
		return []any{}
	}
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice && v.IsNil() {
		return reflect.MakeSlice(v.Type(), 0, 0).Interface()
	}
	return data
}
