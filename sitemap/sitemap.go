package sitemap

import (
	"net/http"
	"sync"

	"github.com/egonelbre/fedwiki"
)

// provides the /system pages
type Sitemap struct {
	Store fedwiki.PageStore

	mu      sync.RWMutex
	headers []*fedwiki.PageHeader
}

func New(store fedwiki.PageStore) *Sitemap {
	sitemap := &Sitemap{}
	sitemap.Store = store
	return sitemap
}

func (sitemap *Sitemap) Update() {
	sitemap.mu.Lock()
	defer sitemap.mu.Unlock()
	sitemap.headers, _ = sitemap.Store.List()
}

func (sitemap *Sitemap) PageChanged(page *fedwiki.Page, err error) {
	//TODO: throttle updating
	sitemap.Update()
}

func (sitemap *Sitemap) Handle(r *http.Request) (code int, template string, data interface{}) {
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
