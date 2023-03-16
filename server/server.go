package server

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/zjyl1994/filebox/assets"
	"github.com/zjyl1994/filebox/service"
	"github.com/zjyl1994/filebox/vars"
	"github.com/zjyl1994/utilz"
	"golang.org/x/net/webdav"
)

func Run() error {
	dav := &webdav.Handler{
		FileSystem: webdav.Dir(vars.DataDir),
		LockSystem: webdav.NewMemLS(),
	}
	handler := &handler{
		davHandler: dav,
	}
	return http.ListenAndServe(vars.Listen, handler)
}

type handler struct {
	davHandler *webdav.Handler
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Filebox")
	diskPath := filepath.Join(vars.DataDir, r.URL.Path)
	queryString := r.URL.Query()
	if r.Method == http.MethodGet {
		if isDir(diskPath) {
			handleDirList(w, r)
			return
		} else if queryString.Has("s") && queryString.Has("e") {
			if service.CheckSecureLink(r.URL.Path, queryString.Get("s"), queryString.Get("e")) {
				http.ServeFile(w, r, diskPath)
			} else {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			}
			return
		} else if r.URL.Path == "/favicon.ico" {
			if utilz.FileExist(diskPath) {
				http.ServeFile(w, r, diskPath)
			} else {
				w.Header().Add("Content-Type", "image/x-icon")
				w.Write(assets.FaviconBytes)
			}
			return
		}
	}

	// Add CORS header for WebDAV when set --cors
	if len(vars.CorsOrigin) > 0 {
		w.Header().Add("Access-Control-Allow-Origin", vars.CorsOrigin)
		w.Header().Add("Access-Control-Allow-Headers", "Overwrite, Destination, Content-Type, Depth, User-Agent, Translate, Range, Content-Range, Timeout, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control, Location, Lock-Token, If")
		w.Header().Add("Access-Control-Allow-Methods", "ACL, CANCELUPLOAD, CHECKIN, CHECKOUT, COPY, DELETE, GET, HEAD, LOCK, MKCALENDAR, MKCOL, MOVE, OPTIONS, POST, PROPFIND, PROPPATCH, PUT, REPORT, SEARCH, UNCHECKOUT, UNLOCK, UPDATE, VERSION-CONTROL")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Expose-Headers", "DAV, Content-Length, Allow, ETag")
		w.Header().Add("Access-Control-Max-Age", "3600")
	}
	if !basicAuthPass(w, r) {
		return
	}

	h.davHandler.ServeHTTP(w, r)
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func basicAuthPass(w http.ResponseWriter, r *http.Request) bool {
	if len(vars.Username) == 0 && len(vars.Password) == 0 {
		return true
	}
	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return false
	}
	if username != vars.Username || password != vars.Password {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return false
	}
	return true
}
