package controller

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/waas-app/WaaS/config"
	"github.com/waas-app/WaaS/util"
	"go.uber.org/zap"
)

func WebsiteRouter(ctx context.Context) *mux.Router {
	router := mux.NewRouter()

	staticFiles, err := filepath.Abs("website/build")
	if err != nil {
		util.Logger(ctx).Fatal("failed to get absolute path to static files", zap.Error(err))
	}

	if _, err := os.Stat(staticFiles); os.IsNotExist(err) {
		// if the static files directory doesn't exist
		// then proxy to a local webpack development server
		// i.e. we're developing wg-access-server locally
		util.Logger(ctx).Info("serving website from webpack dev server")
		u, _ := url.Parse(config.Spec.RootURL)
		router.NotFoundHandler = httputil.NewSingleHostReverseProxy(u)
	} else {
		// if the static files directory exists then
		// handle static file requests.
		// the react app handles routing so we also
		// add a catch-all route to serve the react index page.
		util.Logger(ctx).Info("serving website from static files", zap.String("path", staticFiles))
		router.PathPrefix("/").Handler(
			FileServerWith404(
				http.Dir(staticFiles),
				func(w http.ResponseWriter, r *http.Request) bool {
					http.ServeFile(w, r, filepath.Join(staticFiles, "index.html"))
					return false
				},
			),
		)
	}
	return router
}

// credit: https://gist.github.com/lummie/91cd1c18b2e32fa9f316862221a6fd5c
type FSHandler404 = func(w http.ResponseWriter, r *http.Request) (doDefaultFileServe bool)

// credit: https://gist.github.com/lummie/91cd1c18b2e32fa9f316862221a6fd5c
func FileServerWith404(root http.FileSystem, handler404 FSHandler404) http.Handler {
	fs := http.FileServer(root)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//make sure the url path starts with /
		upath := r.URL.Path
		if !strings.HasPrefix(upath, "/") {
			upath = "/" + upath
			r.URL.Path = upath
		}
		upath = path.Clean(upath)

		// attempt to open the file via the http.FileSystem
		f, err := root.Open(upath)
		if err != nil {
			if os.IsNotExist(err) {
				// call handler
				if handler404 != nil {
					doDefault := handler404(w, r)
					if !doDefault {
						return
					}
				}
			}
		}

		// close if successfully opened
		if err == nil {
			f.Close()
		}

		// default serve
		fs.ServeHTTP(w, r)
	})
}
