// internal/chatwoot/paginate.go
package chatwoot

import (
	"context"

	"github.com/chatwoot/chatwoot-cli/internal/contract"
)

// PageFetcher fetches a single page of results.
type PageFetcher[T any] func(ctx context.Context, page int) ([]T, *contract.Pagination, error)

// ListAll fetches all pages by calling the fetcher repeatedly.
func ListAll[T any](ctx context.Context, fetch PageFetcher[T]) ([]T, *contract.Pagination, error) {
	var all []T
	var lastPag *contract.Pagination

	for page := 1; ; page++ {
		items, pag, err := fetch(ctx, page)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, items...)
		lastPag = pag

		if pag == nil || page >= pag.TotalPages {
			break
		}
	}

	if lastPag != nil {
		lastPag.TotalCount = len(all)
	}

	return all, lastPag, nil
}
