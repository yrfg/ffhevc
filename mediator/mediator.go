package mediator

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

// M is a universal public variable of Mediator
var M *Mediator = &Mediator{}

// Mediator is a toolbox of video processing
type Mediator struct {
}

// VideoTranscode return a VideoTranscoder
func (m *Mediator) VideoTranscode() *VideoTranscoder {
	return NewVideoTransCoder()
}

// VideoTranscoder is a toolset for transcoding
// ! use NewVideoTransCoder() instead of initialize inplace
// ! see https://trac.ffmpeg.org/wiki/Encode/H.265
type VideoTranscoder struct {
	inputFile         string
	outputFile        string
	outVideoCode      string
	outVideoContainer string
	crf               int
	preset            string
}

// PresetMedium medium
const PresetMedium = "medium"

// PresetFast fast
const PresetFast = "fast"

// PresetSlow slow
const PresetSlow = "slow"

// DefaultCrf 28
const DefaultCrf = 28

// DefaultPreset medium
const DefaultPreset = PresetMedium

// OutVideoCodeLibx265 libx265
const OutVideoCodeLibx265 = "libx265"

// DefaultOutVideoCode libx265
const DefaultOutVideoCode = OutVideoCodeLibx265

// OutVideoContainerMp4 mp4
const OutVideoContainerMp4 = "mp4"

// DefaultOutVideoContainer mp4
const DefaultOutVideoContainer = OutVideoContainerMp4

// NewVideoTransCoder create VideoTranscoder
func NewVideoTransCoder() *VideoTranscoder {
	return &VideoTranscoder{
		crf:               DefaultCrf,
		preset:            DefaultPreset,
		outVideoCode:      DefaultOutVideoCode,
		outVideoContainer: DefaultOutVideoContainer,
	}
}

// SetCrf update crf
func (v *VideoTranscoder) SetCrf(crf int) {
	v.crf = crf
}

// SetPreset update preset
func (v *VideoTranscoder) SetPreset(preset string) {
	v.preset = preset
}

// SetInputFile update input file
func (v *VideoTranscoder) SetInputFile(input string) {
	v.inputFile = input
}

// SetOutputFile update output file
// 建议不手动指定输出文件目录。
// ! 设置该项将使 outVideoContainer 无效化。
func (v *VideoTranscoder) SetOutputFile(output string) {
	v.outputFile = output
}

// SetOutVideoContainer update out video container
func (v *VideoTranscoder) SetOutVideoContainer(container string) {
	v.outVideoContainer = container
}

// IsVideoContainerSuffix 判断是否是视频容器格式
func IsVideoContainerSuffix(suffix string) bool {
	switch suffix {
	case "mp4", "mov", "mkv", "avi", "rmvb", "wmv", "flv":
		return true
	default:
		return false
	}
}

func HasVideoContainerSuffix(filename string) bool {
	filenamePart := strings.Split(filename, ".")
	partCount := len(filenamePart)
	if partCount == 1 {
		return false
	}
	return IsVideoContainerSuffix(filenamePart[partCount-1])
}

// DefaultOutputFile 根据输入文件名，生成默认的输出文件名。
func (v *VideoTranscoder) DefaultOutputFile() string {
	inputSplit := strings.Split(v.inputFile, ".")
	l := len(inputSplit)
	if l == 1 {
		return v.inputFile + "_" + v.outVideoCode + "." + v.outVideoContainer
	}
	var b strings.Builder
	for i, s := range inputSplit {
		if i < l-2 {
			b.WriteString(s + ".")
		} else if i == l-2 {
			b.WriteString(s + "_" + v.outVideoCode + ".")
		} else {
			// 判断后缀
			if IsVideoContainerSuffix(s) {
				b.WriteString(v.outVideoContainer)
			} else {
				b.WriteString(s + "." + v.outVideoContainer)
			}
		}
	}
	return b.String()
}

// SetOutputVideoCode update output video code
func (v *VideoTranscoder) SetOutputVideoCode(videoCode string) {
	v.outVideoCode = videoCode
}

// Run trancode begin
func (v *VideoTranscoder) Run() {
	checkoutFfmpeg()
	// 若输出未指定，使用默认输出文件位置。
	if len(v.outputFile) == 0 {
		v.outputFile = v.DefaultOutputFile()
	}
	ffcmd := fmt.Sprintf(
		"-i %s -c:v %s -crf %d -preset %s %s",
		v.inputFile,
		v.outVideoCode,
		v.crf,
		v.preset,
		v.outputFile,
	)

	cmd := exec.Command("ffmpeg", strings.Split(ffcmd, " ")...)

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("VideoTranscoder.Start() failed with '%s'\n", err)
	}
	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()
	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("VideoTranscoder.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}
}

func checkoutFfmpeg() {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Fatal("ffmpeg is not installed")
		os.Exit(1)
	}
}
