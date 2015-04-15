// This package implements sitemap and slugs Handler for a PageStore
package sitemap

import (
	"net/http"
	"sync"

	"github.com/egonelbre/fedwiki"
)

// Handler implements a page handler for
// /system/sitemap
// /system/slugs
type Handler struct {
	Store fedwiki.PageStore

	mu      sync.RWMutex
	headers []*fedwiki.PageHeader
}

func New(store fedwiki.PageStore) *Handler {
	sitemap := &Handler{}
	sitemap.Store = store
	return sitemap
}

func (sitemap *Handler) Update() {
	//TODO: throttle updates
	sitemap.mu.Lock()
	defer sitemap.mu.Unlock()
	sitemap.headers, _ = sitemap.Store.List()
}

func (sitemap *Handler) Handle(r *http.Request) (code int, template string, data interface{}) {
	switch r.URL.Path {
	case "/system/sitemap":
		sitemap.mu.RLock()
		defer sitemap.mu.RUnlock()
		return http.StatusOK, "sitemap", sitemap.headers
	case "/system/slugs":
		sitemap.mu.RLock()
		defer sitemap.mu.RUnlock()
		slugs := make([]fedwiki.Slug, 0, len(sitemap.headers))
		for _, h := range sitemap.headers {
			slugs = append(slugs, h.Slug)
		}

		return http.StatusOK, "slugs", slugs
	}

	return fedwiki.ErrorResponse(http.StatusNotFound, "Page not found.")
}
