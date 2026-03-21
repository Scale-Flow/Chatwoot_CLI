package chatwoot

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Timestamp handles Chatwoot API timestamps that may arrive as either
// an ISO 8601 string or a Unix epoch integer.
type Timestamp struct {
	Value string
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// Try string first (ISO 8601)
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		t.Value = s
		return nil
	}
	// Try number (Unix epoch)
	var n int64
	if err := json.Unmarshal(data, &n); err == nil {
		t.Value = strconv.FormatInt(n, 10)
		return nil
	}
	return fmt.Errorf("timestamp: cannot unmarshal %s", string(data))
}

func (t Timestamp) String() string {
	return t.Value
}
