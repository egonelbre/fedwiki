package server

import (
	"sync"
	"time"

	"github.com/egonelbre/wiki-go-server/page"
)

// provides the /system pages
type Sitemap struct {
	Store page.Store

	mu      sync.RWMutex
	headers []*page.Header
}

func NewSitemap(store page.Store) *Sitemap {
	sitemap := &Sitemap{}
	sitemap.Store = store
	go sitemap.run()
	return sitemap
}

func (sitemap *Sitemap) run() {
	for {
		sitemap.Update()
		time.Sleep(1 * time.Minute)
	}
}

func (sitemap *Sitemap) Update() {
	sitemap.mu.Lock()
	defer sitemap.mu.Unlock()
	sitemap.headers, _ = sitemap.Store.List()
}
