// internal/chatwoot/paginate_test.go
package chatwoot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chatwoot/chatwoot-cli/internal/contract"
)

func TestListAll(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		var items []map[string]any
		switch page {
		case "", "1":
			items = []map[string]any{{"id": 1}, {"id": 2}}
		case "2":
			items = []map[string]any{{"id": 3}, {"id": 4}}
		case "3":
			items = []map[string]any{{"id": 5}, {"id": 6}}
		default:
			items = []map[string]any{}
		}

		resp := map[string]any{
			"data": items,
			"meta": map[string]any{
				"page":        page,
				"total_pages": 3,
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	type Item struct {
		ID int `json:"id"`
	}

	fetcher := func(ctx context.Context, page int) ([]Item, *contract.Pagination, error) {
		c := NewClient(srv.URL, "sk-test", "api_access_token")
		resp, err := c.Do(ctx, http.MethodGet, fmt.Sprintf("/items?page=%d", page), nil)
		if err != nil {
			return nil, nil, err
		}
		defer resp.Body.Close()

		var body struct {
			Data []Item         `json:"data"`
			Meta map[string]any `json:"meta"`
		}
		json.NewDecoder(resp.Body).Decode(&body)
		pag := &contract.Pagination{Page: page, TotalPages: 3, PerPage: 2, TotalCount: 6}
		return body.Data, pag, nil
	}

	items, pag, err := ListAll(context.Background(), fetcher)
	if err != nil {
		t.Fatalf("ListAll error: %v", err)
	}
	if len(items) != 6 {
		t.Errorf("len(items) = %d, want 6", len(items))
	}
	if pag.TotalCount != 6 {
		t.Errorf("TotalCount = %d, want 6", pag.TotalCount)
	}
}

func TestListAllSinglePage(t *testing.T) {
	type Item struct{ ID int }

	fetcher := func(ctx context.Context, page int) ([]Item, *contract.Pagination, error) {
		return []Item{{ID: 1}}, &contract.Pagination{Page: 1, TotalPages: 1, PerPage: 25, TotalCount: 1}, nil
	}

	items, _, err := ListAll(context.Background(), fetcher)
	if err != nil {
		t.Fatalf("ListAll error: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("len(items) = %d, want 1", len(items))
	}
}
