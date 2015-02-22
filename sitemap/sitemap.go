package sitemap

import (
	"net/http"
	"sync"

	"github.com/egonelbre/fedwiki/page"
)

// provides the /system pages
type Sitemap struct {
	Store page.Store

	mu      sync.RWMutex
	headers []*page.Header
}

func New(store page.Store) *Sitemap {
	sitemap := &Sitemap{}
	sitemap.Store = store
	return sitemap
}

func (sitemap *Sitemap) Update() {
	sitemap.mu.Lock()
	defer sitemap.mu.Unlock()
	sitemap.headers, _ = sitemap.Store.List()
}

func (sitemap *Sitemap) PageChanged(p *page.Page, err error) {
	//TODO: throttle updating
	sitemap.Update()
}

func (sitemap *Sitemap) Handle(rw http.ResponseWriter, r *http.Request) (response interface{}, code int, template string) {
	switch r.URL.Path {
	case "/system/sitemap":
		sitemap.mu.RLock()
		defer sitemap.mu.RUnlock()
		return sitemap.headers, http.StatusOK, "sitemap"
	case "/system/slugs":
		sitemap.mu.RLock()
		defer sitemap.mu.RUnlock()
		slugs := make([]page.Slug, 0, len(sitemap.headers))
		for _, h := range sitemap.headers {
			slugs = append(slugs, h.Slug)
		}

		return slugs, http.StatusOK, "slugs"
	}
	return nil, http.StatusNotFound, ""
}
