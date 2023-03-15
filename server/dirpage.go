package server

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/zjyl1994/filebox/assets"
	"github.com/zjyl1994/filebox/service"
	"github.com/zjyl1994/filebox/vars"
	"github.com/zjyl1994/utilz"
)

type fileRenderItem struct {
	service.FileItem
	Link string `json:"link"`
}

func handleDirList(w http.ResponseWriter, r *http.Request) {
	if service.InHiddenPath(r.URL.Path) {
		if !basicAuthPass(w, r) {
			return
		}
	}
	dataMap := make(map[string]any)
	dataMap["base"] = filepath.Base(r.URL.Path)
	dataMap["uppath"] = filepath.Dir(r.URL.Path)
	dataMap["path"] = r.URL.Path
	dataMap["toplevel"] = dataMap["uppath"] == dataMap["path"]
	dataMap["title"] = vars.Title

	diskPath := filepath.Join(vars.DataDir, r.URL.Path)
	fileItems, hasREADME, err := service.ListDir(diskPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if hasREADME { // add readme
		if md, err := service.ProcMarkdown(filepath.Join(diskPath, service.ReadmeFilename), r.URL.Path); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			dataMap["readme"] = md
		}
	}

	renderItems := make([]fileRenderItem, len(fileItems))
	for i, v := range fileItems {
		link := filepath.Join(r.URL.Path, v.Name)
		if !v.IsDir {
			link = link + "?" + service.GenSecureLinkStr(link)
		}
		renderItems[i] = fileRenderItem{
			FileItem: v,
			Link:     link,
		}
	}
	dataMap["items"] = renderItems

	w.Header().Add("Content-Type", "text/html")
	io.WriteString(w, strings.Replace(assets.IndexTemplateString, "__DATA_JSON__", utilz.ToJSONStringNoError(dataMap), 1))
}
