package server

import (
	"sync"

	"github.com/egonelbre/wiki-go-server/page"
)

// provides the /system pages
type Sitemap struct {
	Store page.Store

	mu      sync.RWMutex
	headers []*page.Header
}

//TODO: autoupdate when page.Store.IsDynamic()
func NewSitemap(store page.Store) *Sitemap {
	sitemap := &Sitemap{}
	sitemap.Store = store
	return sitemap
}

func (sitemap *Sitemap) Update() {
	//TODO: throttle updating
	sitemap.mu.Lock()
	defer sitemap.mu.Unlock()
	sitemap.headers, _ = sitemap.Store.List()
}
