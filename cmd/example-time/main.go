// This is an example service for fedwiki
//
// this fedwiki server only has
//   "/status" page
package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/egonelbre/fedwiki"
	"github.com/egonelbre/fedwiki/item"
)

var (
	addr = flag.String("listen", "", "HTTP listening address")
)

func main() {
	flag.Parse()

	port := os.Getenv("PORT")
	if port != "" {
		*addr = "localhost:" + port
	}
	if *addr == "" {
		*addr = ":8080"
	}

	http.ListenAndServe(*addr, &fedwiki.Server{fedwiki.HandlerFunc(serve), nil})
}

func serve(r *http.Request) (code int, template string, data interface{}) {
	if r.Method != "GET" {
		return fedwiki.ErrorResponse(http.StatusForbidden, "Method %s is not allowed", r.Method)
	}

	switch r.URL.Path {
	case "/status":
		return http.StatusOK, "", fedwiki.Page{
			PageHeader: fedwiki.PageHeader{
				Slug:  "/status",
				Title: "Status",
				Date:  fedwiki.Now(),
			},
			Story: fedwiki.Story{
				item.Paragraph("Current time is: " + time.Now().String()),
			},
		}
	default:
		return fedwiki.ErrorResponse(http.StatusNotFound, "Page %s not found", r.URL.Path)
	}
}
