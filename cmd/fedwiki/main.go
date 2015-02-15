package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/egonelbre/wiki-go-server/page/folderstore"
	"github.com/egonelbre/wiki-go-server/server"
)

type Dirs struct {
	Default  string
	Data     string
	Packages string
	Client   string
	Pages    string
	Status   string
	Identity string
}

var (
	addr = flag.String("listen", "", "HTTP listening address")

	dirdefault  = flag.String("default", "default-data", "directory for default-data")
	dirdata     = flag.String("data", "data", "directory for storing all data")
	dirpackages = flag.String("packages", "packages", "directory storing all the packages (relative to data)")
	dirclient   = flag.String("client", "wiki-client", "directory for client (relative to data)")
	dirpages    = flag.String("pages", "pages", "directory for storing pages (relative to data)")
	dirstatus   = flag.String("status", "status", "directory for storing status (relative to data)")
	diridentity = flag.String("identity", "status/persona.identity", "directory for storing identity (relative to data)")
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

	datadir := *dirdata
	if !filepath.IsAbs(datadir) {
		var err error
		datadir, err = filepath.Abs(datadir)
		check(err)
	}

	defdir := *dirdefault
	if !filepath.IsAbs(defdir) {
		var err error
		defdir, err = filepath.Abs(defdir)
		check(err)
	}

	dir := Dirs{
		Default:  defdir,
		Data:     datadir,
		Packages: filepath.Join(datadir, *dirpackages),
		Client:   filepath.Join(datadir, *dirclient),
		Pages:    filepath.Join(datadir, *dirpages),
		Status:   filepath.Join(datadir, *dirstatus),
		Identity: filepath.Join(datadir, *diridentity),
	}

	// if we don't have a pages directory assume that we haven't
	// setup the content yet and copy everything from default data
	if _, err := os.Stat(dir.Pages); os.IsNotExist(err) {
		check(copyfiles(dir.Default, dir.Data))
	}

	store := folderstore.New(dir.Pages)
	server := server.New(store)

	log.Printf("Listening on %v...\n", *addr)
	check(http.ListenAndServe(*addr,
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			log.Printf("REQ '%s' > '%s'\n", r.URL, r.URL.Path)
			if strings.HasPrefix(r.URL.Path, "/client/") {
				upath := filepath.Join("client", "client", r.URL.Path[len("/client"):])
				fmt.Println(upath)
				http.ServeFile(rw, r, path.Clean(upath))
				return
			}
			if strings.HasPrefix(r.URL.Path, "/static/") {
				upath := filepath.Join("static", r.URL.Path[len("/static"):])
				fmt.Println(upath)
				http.ServeFile(rw, r, path.Clean(upath))
				return
			}
			if r.URL.Path == "" || r.URL.Path == "/" {
				http.ServeFile(rw, r, "client.html")
				return
			}

			server.ServeHTTP(rw, r)
		})))
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
