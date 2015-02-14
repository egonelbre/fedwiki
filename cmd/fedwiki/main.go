package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
	addr = flag.String("listen", "localhost:8080", "HTTP listening address")

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
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	datadir := *dirdata
	if !filepath.IsAbs(datadir) {
		var err error
		datadir, err = filepath.Abs(*datadir)
		check(err)
	}

	dir := Dirs{
		Content:  datadir,
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
}

func copyfiles(src, dst string) error {
	return filepath.Walk(src,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return os.Mkdir(filepath.Join(dst, path), info.Mode())
			}

			data, err := ioutil.ReadFile(filepath.Join(src, path))
			if err != nil {
				return err
			}
			return ioutil.WriteFile(filepath.Join(dst, path), data, info.Mode())
		})
}
