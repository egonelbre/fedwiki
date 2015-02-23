package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/egonelbre/fedwiki"
	"github.com/egonelbre/fedwiki/pagestore"
	"github.com/egonelbre/fedwiki/pagestore/folderstore"
	"github.com/egonelbre/fedwiki/renderer"
	"github.com/egonelbre/fedwiki/sitemap"
)

var (
	addr = flag.String("listen", "", "HTTP listening address")

	clientpage  = flag.String("client", "client.html", "client html page")
	dirdefpages = flag.String("default", "default-pages", "directory for default pages")
	dirviews    = flag.String("views", "views", "directory for storing views")
	dirstatic   = flag.String("static", "static", "directory for storing static content")

	dirpages  = flag.String("pages", filepath.Join("data", "pages"), "directory for storing pages")
	dirstatus = flag.String("status", filepath.Join("data", "status"), "directory for storing status")
)

func absolute(s string) string {
	if filepath.IsAbs(s) {
		return s
	}
	r, err := filepath.Abs(s)
	if err != nil {
		log.Fatal(err)
	}
	return r
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	port := os.Getenv("PORT")
	if port != "" {
		*addr = "localhost:" + port
	}
	if *addr == "" {
		*addr = ":8080"
	}

	*clientpage = absolute(*clientpage)
	*dirdefpages = absolute(*dirdefpages)
	*dirpages = absolute(*dirpages)
	*dirstatus = absolute(*dirstatus)
	*dirviews = absolute(*dirviews)
	*dirstatic = absolute(*dirstatic)

	// if we don't have a pages directory assume that we haven't
	// setup the content yet and copy everything from default data
	if _, err := os.Stat(*dirpages); os.IsNotExist(err) {
		check(copyfiles(*dirdefpages, *dirpages))
	}

	store := folderstore.New(*dirpages)

	sitemap := sitemap.New(store)
	sitemap.Update()

	render := renderer.New(filepath.Join(*dirviews, "*"))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(*dirstatic))))

	http.Handle("/system/sitemap", &fedwiki.Server{sitemap, render})
	http.Handle("/system/slugs", &fedwiki.Server{sitemap, render})

	pageserver := &fedwiki.Server{pagestore.Handler{store}, render}
	http.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/favicon.png" {
			http.ServeFile(rw, r, filepath.Join(*dirstatus, "favicon.png"))
			return
		}

		if r.URL.Path == "" || r.URL.Path == "/" {
			http.ServeFile(rw, r, *clientpage)
			return
		}

		pageserver.ServeHTTP(rw, r)
	}))

	log.Printf("Listening on %v...\n", *addr)
	check(http.ListenAndServe(*addr, nil))
}

func copyfiles(src, dst string) error {
	return filepath.Walk(src,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			path, err = filepath.Rel(src, path)
			if err != nil {
				return err
			}

			if info.IsDir() {
				return os.MkdirAll(filepath.Join(dst, path), info.Mode())
			}

			data, err := ioutil.ReadFile(filepath.Join(src, path))
			if err != nil {
				return err
			}
			return ioutil.WriteFile(filepath.Join(dst, path), data, info.Mode())
		})
}
