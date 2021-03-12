package mediator

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

// IsVideoContainerSuffix 判断是否是视频容器格式
func IsVideoContainerSuffix(suffix string) bool {
	switch suffix {
	case "mp4", "mov", "mkv", "avi", "rmvb", "wmv", "flv":
		return true
	default:
		return false
	}
}

// (param)file: filename or file path
func HasVideoContainerSuffix(file string) bool {
	_, filename := path.Split(file)
	frags := strings.Split(filename, ".")
	if len(frags) == 1 {
		return false
	}
	return IsVideoContainerSuffix(frags[len(frags)-1])
}

func GetReplacedVideoContainerSuffix(oldFilename, newSuffix string) string {
	frags := strings.Split(oldFilename, ".")
	l := len(frags)
	if l == 1 {
		return oldFilename + "." + newSuffix
	}
	var b strings.Builder
	for i, s := range frags {
		if i < l-1 {
			b.WriteString(s + ".")
		} else {
			// 判断后缀
			// 非视频后缀保留，视频后缀做替换。
			if IsVideoContainerSuffix(s) {
				b.WriteString(newSuffix)
			} else {
				b.WriteString(s + "." + newSuffix)
			}
		}
	}
	return b.String()
}

func CheckFFmpeg() {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Fatal("ffmpeg is not installed")
		os.Exit(1)
	}
}

// GetDirVideos
// 不递归
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
		if !HasVideoContainerSuffix(filename) {
			continue
		}
		fullPath := path.Join(dir, filename)
		files = append(files, fullPath)
	}
	return files
}
