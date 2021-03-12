package main

import (
	"fmt"

	"github.com/yrfg/ffhevc/mediator"
)

func main() {
	v := mediator.M.VideoTranscode()
	dir := "xxxx"
	tasklist := mediator.GetDirVideos(dir, true)
	fmt.Println("任务数量: ", len(tasklist))
	v.SetOutputDir("xxxx")
	for _, task := range tasklist {
		v.SetInputFile(task)
		v.Run()
	}
}
