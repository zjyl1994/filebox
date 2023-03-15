package service

import (
	"io/ioutil"
	"sort"
	"strings"
	"time"
)

const (
	timeFormat     = "2006-01-02 15:04:05"
	ReadmeFilename = "README.md"
)

type fileInfo struct {
	Name    string
	IsDir   bool
	LastMod time.Time
	Size    int64
}

type FileItem struct {
	Name    string `json:"name"`
	IsDir   bool   `json:"is_dir"`
	LastMod int64  `json:"last_mod"`
	Size    int64  `json:"size"`
}

func ListDir(path string) (result []FileItem, hasREADME bool, err error) {
	infos, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, false, err
	}

	fileInfos := make([]fileInfo, 0, len(infos))
	for _, v := range infos {
		name := v.Name()
		if name == ReadmeFilename {
			hasREADME = true
		}
		if !strings.HasPrefix(name, ".") {
			fileInfos = append(fileInfos, fileInfo{
				Name:    name,
				IsDir:   v.IsDir(),
				LastMod: v.ModTime(),
				Size:    v.Size(),
			})
		}
	}

	sort.Slice(fileInfos, func(i, j int) bool {
		if fileInfos[i].IsDir != fileInfos[j].IsDir {
			return fileInfos[i].IsDir
		} else {
			return fileInfos[i].Name < fileInfos[j].Name
		}
	})

	result = make([]FileItem, len(fileInfos))
	for i, v := range fileInfos {
		result[i] = FileItem{
			Name:    v.Name,
			IsDir:   v.IsDir,
			LastMod: v.LastMod.Unix(),
			Size:    v.Size,
		}
	}

	return result, hasREADME, nil
}

func InHiddenPath(path string) bool {
	parts := strings.Split(path, "/")
	for _, name := range parts {
		if strings.HasPrefix(name, ".") {
			return true
		}
	}
	return false
}
