package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/egonelbre/fedwiki"
	"github.com/egonelbre/fedwiki/pagestore"
	"github.com/egonelbre/fedwiki/pagestore/folderstore"
	"github.com/egonelbre/fedwiki/plugin"
	"github.com/egonelbre/fedwiki/sitemap"
	"github.com/egonelbre/fedwiki/template"
)

var (
	addr = flag.String("listen", "", "HTTP listening address")

	clientpage  = flag.String("client", "client.html", "client html page")
	dirdefpages = flag.String("default", "default-pages", "directory for default pages")
	dirviews    = flag.String("views", "views", "directory for storing views")
	dirstatic   = flag.String("static", "static", "directory for storing static content")
	dirplugins  = flag.String("plugins", "plugins", "directory for storing plugins content")

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
	*dirplugins = absolute(*dirplugins)

	// if we don't have a pages directory assume that we haven't
	// setup the content yet and copy everything from default data
	if _, err := os.Stat(*dirpages); os.IsNotExist(err) {
		log.Println("Initializing default pages.")
		check(copyfiles(*dirdefpages, *dirpages))
		check(copyglob(filepath.Join(*dirplugins, "*", "pages", "*"), *dirpages))
	}

	render := template.New(filepath.Join(*dirviews, "*"))

	mainstore := folderstore.New(*dirpages)
	sitemap := sitemap.New(mainstore)
	sitemap.Update()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(*dirstatic))))

	plugins := plugin.NewServer(*dirplugins)
	plugins.Update()
	pluginsserver := &fedwiki.Server{plugins, render}

	http.Handle("/plugins/", plugins)

	sitemapserver := &fedwiki.Server{sitemap, render}

	pageserver := &fedwiki.Server{pagestore.Handler{mainstore}, render}
	http.Handle("/",
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/system/sitemap") ||
				strings.HasPrefix(r.URL.Path, "/system/slugs") {
				sitemapserver.ServeHTTP(rw, r)
				return
			}

			if strings.HasPrefix(r.URL.Path, "/system/plugins") ||
				strings.HasPrefix(r.URL.Path, "/system/factories") {
				pluginsserver.ServeHTTP(rw, r)
				return
			}

			if r.URL.Path == "/favicon.png" {
				http.ServeFile(rw, r, filepath.Join(*dirstatus, "favicon.png"))
				return
			}

			if r.URL.Path == "" || r.URL.Path == "/" {
				http.ServeFile(rw, r, *clientpage)
				return
			}

			if r.Method != "GET" {
				defer sitemap.Update()
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

func copyglob(glob, dst string) error {
	matches, err := filepath.Glob(glob)
	if err != nil {
		return err
	}
	for _, match := range matches {
		if filepath.Ext(match) != "" {
			continue
		}

		data, err := ioutil.ReadFile(match)
		if err != nil {
			continue
		}

		info, err := os.Stat(match)
		if err != nil || info.IsDir() {
			continue
		}

		_ = ioutil.WriteFile(filepath.Join(dst, filepath.Base(match)), data, info.Mode())
	}
	return nil
}
