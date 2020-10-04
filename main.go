package main

import (
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/yrfg/ffhevc/mediator"
)

func GetDirVideos(dir string, ignoreHidden bool) []string {
	fileInfoList, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	files := make([]string, 0)
	for _, f := range fileInfoList {
		if f.IsDir() {
			continue
		}
		filename := f.Name()
		if ignoreHidden && strings.HasPrefix(filename, ".") {
			continue
		}
		if !mediator.HasVideoContainerSuffix(filename) {
			continue
		}
		fullPath := path.Join(dir, filename)
		files = append(files, fullPath)
	}
	return files
}

func main() {
	v := mediator.M.VideoTranscode()
	v.SetCrf(33)
	tasklist := GetDirVideos("/Users/yira/Documents/yira/collection/剪辑", true)
	for _, task := range tasklist {
		v.SetInputFile(task)
		v.SetOutputFile("")
		v.Run()
	}
}
