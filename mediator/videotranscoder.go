package mediator

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
)

// VideoTranscoder is a toolset for transcoding
// ! use NewVideoTransCoder() instead of initialize inplace
// ! see https://trac.ffmpeg.org/wiki/Encode/H.265
// 结构体上的字段的意义是配置，不应被改动。
type VideoTranscoder struct {
	inputFile         string
	outputDir         string
	outputName        string
	outVideoCoding    string
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

// DefaultPreset medium
const DefaultPreset = PresetMedium

// DefaultCrf 30
const DefaultCrf = 30

// OutVideoCodingLibx265 libx265
const OutVideoCodingLibx265 = "libx265"

// DefaultOutVideoCoding libx265
const DefaultOutVideoCoding = OutVideoCodingLibx265

// OutVideoContainerMp4 mp4
const OutVideoContainerMp4 = "mp4"

// DefaultOutVideoContainer mp4
const DefaultOutVideoContainer = OutVideoContainerMp4

// NewVideoTransCoder create VideoTranscoder
func NewVideoTransCoder() *VideoTranscoder {
	return &VideoTranscoder{
		crf:               DefaultCrf,
		preset:            DefaultPreset,
		outVideoCoding:    DefaultOutVideoCoding,
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

func (v *VideoTranscoder) SetOutputDir(out string) {
	v.outputDir = out
}

// SetOutputName update output file
func (v *VideoTranscoder) SetOutputName(name string) {
	v.outputName = name
}

// SetOutVideoContainer update out video container
func (v *VideoTranscoder) SetOutVideoContainer(container string) {
	v.outVideoContainer = container
}

// SetOutputVideoCoding update output video code
func (v *VideoTranscoder) SetOutputVideoCoding(videoCode string) {
	v.outVideoCoding = videoCode
}

// Run trancode begin
func (v *VideoTranscoder) Run() {
	CheckFFmpeg()
	if v.inputFile == "" {
		log.Fatal("no inputfile")
		os.Exit(1)
	}

	inDir, inFilename := path.Split(v.inputFile)
	outDir, outFilename := v.outputDir, v.outputName
	// 若输出目录未指定，使用相同目录
	if outDir == "" {
		outDir = inDir
	}
	// 若输出文件名未指定，使用默认文件名规则
	if outFilename == "" {
		outFilename = GetReplacedVideoContainerSuffix(inFilename, v.outVideoContainer)
	}

	outPath := path.Join(outDir, outFilename)

	ffcmd := fmt.Sprintf(
		"-i %s -c:v %s -f %s -crf %d -preset %s %s",
		v.inputFile,
		v.outVideoCoding,
		v.outVideoContainer,
		v.crf,
		v.preset,
		outPath,
	)

	fmt.Println(ffcmd)

	cmd := exec.Command(
		"ffmpeg",
		"-i", v.inputFile,
		"-c:v", v.outVideoCoding,
		"-f", v.outVideoContainer,
		"-crf", strconv.Itoa(v.crf),
		"-preset", v.preset,
		outPath,
	)

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
		log.Fatalf("VideoTranscoder Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}
}
