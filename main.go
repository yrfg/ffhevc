package main

import (
	"github.com/yrfg/ffhevc/mediator"
)

func main() {
	v := mediator.M.VideoTranscode()
	v.SetCrf(32)
	v.SetInputFile("/Volumes/MyBook_Data/Phone_Video1.mov")
	v.Run()
}
