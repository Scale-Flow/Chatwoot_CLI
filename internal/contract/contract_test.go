package contract

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestSuccessSingleResource(t *testing.T) {
	resp := Success(map[string]any{"id": 42, "name": "Support Inbox"})

	if !resp.OK {
		t.Fatal("expected ok to be true")
	}
	if resp.Error != nil {
		t.Fatal("expected error to be nil")
	}

	var buf bytes.Buffer
	if err := Write(&buf, resp, false); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(buf.Bytes(), &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if raw["ok"] != true {
		t.Errorf("ok = %v, want true", raw["ok"])
	}
	data, ok := raw["data"].(map[string]any)
	if !ok {
		t.Fatal("data is not an object")
	}
	if data["id"] != float64(42) {
		t.Errorf("data.id = %v, want 42", data["id"])
	}
	if _, exists := raw["meta"]; exists {
		t.Error("meta should be absent for single resource")
	}
	if _, exists := raw["warnings"]; exists {
		t.Error("warnings should be absent when empty")
	}
	if _, exists := raw["error"]; exists {
		t.Error("error should be absent on success")
	}
}

func TestSuccessCollection(t *testing.T) {
	items := []map[string]any{
		{"id": 1, "name": "Alice"},
		{"id": 2, "name": "Bob"},
	}
	meta := Meta{
		Pagination: &Pagination{
			Page: 1, PerPage: 25, TotalCount: 142, TotalPages: 6,
		},
	}
	resp := SuccessList(items, meta)

	var buf bytes.Buffer
	if err := Write(&buf, resp, false); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(buf.Bytes(), &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	data, ok := raw["data"].([]any)
	if !ok {
		t.Fatal("data is not an array")
	}
	if len(data) != 2 {
		t.Errorf("data length = %d, want 2", len(data))
	}

	metaRaw, ok := raw["meta"].(map[string]any)
	if !ok {
		t.Fatal("meta is not an object")
	}
	pag, ok := metaRaw["pagination"].(map[string]any)
	if !ok {
		t.Fatal("pagination is not an object")
	}
	if pag["total_count"] != float64(142) {
		t.Errorf("total_count = %v, want 142", pag["total_count"])
	}
	if pag["per_page"] != float64(25) {
		t.Errorf("per_page = %v, want 25", pag["per_page"])
	}
}

func TestSuccessEmptyCollection(t *testing.T) {
	var items []map[string]any // nil slice
	meta := Meta{
		Pagination: &Pagination{
			Page: 1, PerPage: 25, TotalCount: 0, TotalPages: 0,
		},
	}
	resp := SuccessList(items, meta)

	var buf bytes.Buffer
	if err := Write(&buf, resp, false); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(buf.Bytes(), &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	data, ok := raw["data"].([]any)
	if !ok {
		t.Fatalf("data is not an array, got %T", raw["data"])
	}
	if len(data) != 0 {
		t.Errorf("data length = %d, want 0", len(data))
	}
}

func TestSuccessWithWarnings(t *testing.T) {
	resp := Success(map[string]any{"id": 42})
	resp.Warnings = []Warning{
		{Code: "deprecated_endpoint", Message: "This endpoint will be removed"},
	}

	var buf bytes.Buffer
	if err := Write(&buf, resp, false); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(buf.Bytes(), &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	warnings, ok := raw["warnings"].([]any)
	if !ok {
		t.Fatal("warnings is not an array")
	}
	if len(warnings) != 1 {
		t.Errorf("warnings length = %d, want 1", len(warnings))
	}
}

func TestSuccessWithoutWarningsOmitsField(t *testing.T) {
	resp := Success(map[string]any{"id": 42})

	var buf bytes.Buffer
	if err := Write(&buf, resp, false); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	if bytes.Contains(buf.Bytes(), []byte(`"warnings"`)) {
		t.Error("warnings field should be absent when empty")
	}
}

func TestErrorEnvelope(t *testing.T) {
	resp := Err(ErrCodeNotFound, "Conversation 12345 not found")

	var buf bytes.Buffer
	if err := Write(&buf, resp, false); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(buf.Bytes(), &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if raw["ok"] != false {
		t.Errorf("ok = %v, want false", raw["ok"])
	}
	if _, exists := raw["data"]; exists {
		t.Error("data should be absent on error")
	}

	errObj, ok := raw["error"].(map[string]any)
	if !ok {
		t.Fatal("error is not an object")
	}
	if errObj["code"] != "not_found" {
		t.Errorf("error.code = %v, want not_found", errObj["code"])
	}
	if errObj["message"] != "Conversation 12345 not found" {
		t.Errorf("error.message = %v, want expected message", errObj["message"])
	}
}

func TestErrorWithDetail(t *testing.T) {
	detail := map[string]any{"fields": map[string]any{"email": "is not valid"}}
	resp := ErrWithDetail(ErrCodeValidation, "Invalid request", detail)

	var buf bytes.Buffer
	if err := Write(&buf, resp, false); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(buf.Bytes(), &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	errObj := raw["error"].(map[string]any)
	if errObj["detail"] == nil {
		t.Error("error.detail should be present")
	}
}

func TestPrettyOutput(t *testing.T) {
	resp := Success(map[string]any{"id": 1})

	var compact bytes.Buffer
	if err := Write(&compact, resp, false); err != nil {
		t.Fatalf("Write compact error: %v", err)
	}

	var pretty bytes.Buffer
	if err := Write(&pretty, resp, true); err != nil {
		t.Fatalf("Write pretty error: %v", err)
	}

	if bytes.Contains(compact.Bytes(), []byte("  ")) {
		t.Error("compact output should not contain indentation")
	}
	if !bytes.Contains(pretty.Bytes(), []byte("  ")) {
		t.Error("pretty output should contain indentation")
	}
	if pretty.Len() <= compact.Len() {
		t.Error("pretty output should be longer than compact")
	}
}

func TestTrailingNewline(t *testing.T) {
	resp := Success(map[string]any{"id": 1})

	var buf bytes.Buffer
	if err := Write(&buf, resp, false); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	if buf.Bytes()[buf.Len()-1] != '\n' {
		t.Error("output should end with a newline")
	}
}

func TestValidationRejectsOkTrueWithError(t *testing.T) {
	resp := Response{OK: true, Data: "x", Error: &ErrorDetail{Code: "bad", Message: "bad"}}
	var buf bytes.Buffer
	err := Write(&buf, resp, false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidationRejectsOkFalseWithData(t *testing.T) {
	resp := Response{OK: false, Data: "x", Error: &ErrorDetail{Code: "bad", Message: "bad"}}
	var buf bytes.Buffer
	err := Write(&buf, resp, false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidationRejectsOkFalseWithoutError(t *testing.T) {
	resp := Response{OK: false}
	var buf bytes.Buffer
	err := Write(&buf, resp, false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestSnakeCaseFieldNames(t *testing.T) {
	meta := Meta{
		Pagination: &Pagination{
			Page: 1, PerPage: 10, TotalCount: 50, TotalPages: 5,
		},
	}
	resp := SuccessList([]any{}, meta)

	var buf bytes.Buffer
	if err := Write(&buf, resp, false); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	output := buf.String()
	for _, field := range []string{"per_page", "total_count", "total_pages"} {
		if !bytes.Contains(buf.Bytes(), []byte(field)) {
			t.Errorf("expected snake_case field %q in output: %s", field, output)
		}
	}
	for _, bad := range []string{"perPage", "totalCount", "totalPages", "PerPage", "TotalCount"} {
		if bytes.Contains(buf.Bytes(), []byte(bad)) {
			t.Errorf("unexpected camelCase field %q in output: %s", bad, output)
		}
	}
}
