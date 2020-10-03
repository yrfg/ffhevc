package main

import (
	"github.com/yrfg/ffhevc/mediator"
)

func main() {
	v := mediator.M.VideoTranscode()
	v.SetCrf(32)
	tasklist := []string{"/Volumes/MyBook_Data/Phone_Video1.mov"}
	for _, task := range tasklist {
		v.SetInputFile(task)
		v.Run()
	}
}
